package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	syncWord = 0xFFF

	MPEG = 0x1

	// Layer
	layerReserved = 0x0
	layer3        = 0x1
	layer2        = 0x2
	layer1        = 0x3

	protected = 0x0

	// Mode specified
	modeStereo        = 0x0
	modeJoinStereo    = 0x1
	modeDualChannel   = 0x2
	modeSingleChannel = 0x3
)

var bitrateSpecified = [3][16]int{
	// Layer 3
	{0, 32000, 40000, 48000, 56000, 64000, 80000, 96000,
		112000, 128000, 160000, 192000, 224000, 256000, 320000}, // kBit/s
	// Layer 2
	{0, 32000, 48000, 56000, 64000, 80000, 96000, 112000,
		128000, 160000, 192000, 224000, 256000, 320000, 384000}, // kBit/s
	// Layer 1
	{0, 32000, 64000, 96000, 128000, 160000, 192000, 224000,
		256000, 288000, 320000, 352000, 384000, 416000, 448000}, // kBit/s
}

var frequencySpecified = [4]int{
	44100, 48000, 32000, 0, // kHz
}

type Header struct {
	SyncWord          uint16 // 12 bits
	Id                byte   // 1 bit
	Layer             byte   // 2 bits
	ProtectionBit     byte   // 1 bit
	BitrateIndex      byte   // 4 bits
	SamplingFrequency byte   // 2 bits
	PaddingBit        byte   // 1 bits
	PrivateBit        byte   // 1 bit
	Mode              byte   // 2 bits
	ModeExtension     byte   // 2 bits
	Copyright         byte   // 1 bit
	Copy              byte   // 1 bit
	Emphasis          byte   // 2 bits
}

type SideInformation struct {
	MainDataBegin        uint16         // 9 bits
	PrivateBits          byte           // 3 bits in mono, 5 in stereo
	Scfsi                [2][4]byte     // channels/bands 1 bit
	Part23lenght         [2][2]uint16   // granule/channel 12 bits
	BigValues            [2][2]uint16   // granule/channel 9 bits
	GlobalGain           [2][2]uint8    // granule/channel 8 bits
	ScalefacCompress     [2][2]byte     // granule/channel 4 bits
	WindowsSwitchingFlag [2][2]byte     // granule/channel 1 bit
	BlockType            [2][2]byte     // granule/channel 2 bit
	MixedBlockFlag       [2][2]uint8    // granule/channel 1 bit
	TableSelect          [2][2][3]byte  // granule/channel/region 5 bits
	SubblockGain         [2][2][3]uint8 // granule/channel/window 3 bits
	Region0Count         [2][2]byte     // granule/channel 4 bits
	Region1Count         [2][2]byte     // granule/channel 3 bits
	Preflag              [2][2]byte     // granule/channel 1 bit
	ScalfacScale         [2][2]byte     // granule/channel 1 bit
	Count1tableSelect    [2][2]byte     // granule/channel 1 bit
}

type BitReader struct {
	offset uint
	bytes  []byte
}

func (si *BitReader) Read(n uint) int {
	buf := make([]byte, 3)

	offset := si.offset / 8
	copy(buf, si.bytes[offset:])
	//fmt.Printf("%08b", buf)

	r := int(buf[0])<<16 | int(buf[1])<<8 | int(buf[1])
	r = r >> (uint(24) - n - si.offset%8) & (0xFFFFFF >> (24 - n))

	si.offset = si.offset + n

	return r
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	content, err := ioutil.ReadFile("")
	if err != nil {
		log.Fatal(err)
	}
	file := bytes.NewReader(content)

	ReadID3(file)

	c := 3
	// Read frames ====================================================================================================
	for file.Len() != 0 {
		// Frame Header ==============================================================================================
		var buf uint32 // 4 bytes
		if err := binary.Read(file, binary.BigEndian, &buf); err != nil {
			log.Println(err)
		}
		//fmt.Printf("%032b\n", buf)

		header := Header{
			uint16(buf >> 20 & 0xFFF),
			uint8(buf >> 19 & 0x1),
			uint8(buf >> 17 & 0x3),
			uint8(buf >> 16 & 0x1),
			uint8(buf >> 12 & 0xF),
			uint8(buf >> 10 & 0x3),
			uint8(buf >> 9 & 0x1),
			uint8(buf >> 8 & 0x1),
			uint8(buf >> 6 & 0x3),
			uint8(buf >> 4 & 0x3),
			uint8(buf >> 3 & 0x1),
			uint8(buf >> 2 & 0x1),
			uint8(buf >> 0 & 0x3),
		}

		if header.SyncWord != syncWord {
			break
		}
		//fmt.Printf("%+v\n", header)

		// Number of channels; equals 1 for single_channel mode, equals 2 for other modes.
		nch := 2
		if header.Mode == modeSingleChannel {
			nch = 1
		}

		// Frame length =============================================================================================
		frameSize := 144 // byte for layer2 and layer3 1152/(1b*8bit) = 144; for layer1 384/(4b*8bit) = 12
		bitrate := bitrateSpecified[header.Layer-1][header.BitrateIndex]
		frequency := frequencySpecified[header.SamplingFrequency]
		frameLength := frameSize*bitrate/frequency + int(header.PaddingBit)
		fmt.Println(frameSize, bitrate/1000, frequency, frameLength)

		// CRC Check =================================================================================================
		if header.ProtectionBit == protected {
			crc := make([]byte, 2)
			file.Read(crc)
		}

		// AudioData =================================================================================================
		// Side Information ==========================================================================================
		sideInformationLength := 32
		if header.Mode == modeSingleChannel {
			sideInformationLength = 17
		}
		bitReader := BitReader{
			0,
			make([]uint8, sideInformationLength),
		}
		file.Read(bitReader.bytes)
		//fmt.Printf("%08b\n", bitReader.bytes)

		//fmt.Printf("main_data_begin %09b\n", bitReader.Read(9))   // 9 bits
		//fmt.Printf("private_bits %03b\n", bitReader.Read(3))      // 3 bits in mono, 5 in stereo
		//fmt.Printf("scfsi ch1 (share) %04b\n", bitReader.Read(4)) // 1 bit * 4
		//fmt.Printf("scfsi ch2 (share) %04b\n", bitReader.Read(4)) // 1 bit * 4
		//for gr := 0; gr < 2; gr++ { // 2 granules for MPEG1, 1 granules for MPEG2
		//	for ch := 0; ch < 2; ch++ {
		//		fmt.Println("--------------------------------------------------")
		//		fmt.Printf("par2_3_lenth %012b\n", bitReader.Read(12))    // 12 bits
		//		fmt.Printf("big_values %09b\n", bitReader.Read(9))        // 9 bits
		//		fmt.Printf("global_gain %08b\n", bitReader.Read(8))       // 8 bits
		//		fmt.Printf("scalefac_compress %04b\n", bitReader.Read(4)) // 4 bits
		//		windowsSwitchingFlag := bitReader.Read(1)
		//		fmt.Printf("windows_switching_flag %b\n", windowsSwitchingFlag) // 1 bit
		//
		//		if windowsSwitchingFlag == 1 {
		//			fmt.Printf("block_type %02b\n", bitReader.Read(2))     // 2 bits
		//			fmt.Printf("mixed_blockflag %b\n", bitReader.Read(1))  // 1 bit
		//			fmt.Printf("table_select %015b\n", bitReader.Read(10)) // 2 * 5 bits
		//			fmt.Printf("subblock_gain %09b\n", bitReader.Read(9))  // 3 * 3 bits
		//
		//		} else {
		//			fmt.Printf("table_select %015b\n", bitReader.Read(15)) // 3 * 5 bits
		//			fmt.Printf("region0_count %04b\n", bitReader.Read(4))  // 4 bits
		//			fmt.Printf("region1_count %03b\n", bitReader.Read(3))  // 3 bits
		//		}
		//
		//		fmt.Printf("preflag %b\n", bitReader.Read(1))            // 1 bit
		//		fmt.Printf("scalfac_scale %b\n", bitReader.Read(1))      // 1 bit
		//		fmt.Printf("count1table_select %b\n", bitReader.Read(1)) // 1 bit
		//	}
		//}

		sideInformation := SideInformation{}
		sideInformation.MainDataBegin = uint16(bitReader.Read(9)) // main_data_begin

		if header.Mode == modeSingleChannel {
			sideInformation.PrivateBits = byte(bitReader.Read(5)) // private_bits
		} else {
			sideInformation.PrivateBits = byte(bitReader.Read(3)) // private_bits
		}

		for ch := 0; ch < nch; ch++ {
			for band := 0; band < 4; band++ {
				sideInformation.Scfsi[ch][band] = byte(bitReader.Read(1)) // scfsi[ch][scfsi_band]
			}
		}

		for gr := 0; gr < 2; gr++ { // 2 granules for MPEG1, 1 granules for MPEG2
			for ch := 0; ch < nch; ch++ {
				sideInformation.Part23lenght[gr][ch] = uint16(bitReader.Read(12))      // part2_3_length[gr][ch]
				sideInformation.BigValues[gr][ch] = uint16(bitReader.Read(9))          // big_values[gr][ch]
				sideInformation.GlobalGain[gr][ch] = uint8(bitReader.Read(8))          // global_gain[gr][ch]
				sideInformation.ScalefacCompress[gr][ch] = byte(bitReader.Read(4))     // scalefac_compress[gr][ch]
				sideInformation.WindowsSwitchingFlag[gr][ch] = byte(bitReader.Read(1)) // window_switching_flag[gr][ch]

				if sideInformation.WindowsSwitchingFlag[gr][ch] == 1 {
					sideInformation.BlockType[gr][ch] = byte(bitReader.Read(2))       // block_type[gr][ch]
					sideInformation.MixedBlockFlag[gr][ch] = uint8(bitReader.Read(1)) // mixed_block_flag[gr][ch]

					for region := 0; region < 2; region++ {
						sideInformation.TableSelect[gr][ch][region] = byte(bitReader.Read(5)) // table_select[gr][ch][region]
					}

					for window := 0; window < 3; window++ {
						sideInformation.SubblockGain[gr][ch][window] = uint8(bitReader.Read(3)) // subblock_gain[gr][ch][window]
					}

				} else {
					for region := 0; region < 3; region++ {
						sideInformation.TableSelect[gr][ch][region] = byte(bitReader.Read(5)) // table_select[gr][ch][region]
					}

					sideInformation.Region0Count[gr][ch] = byte(bitReader.Read(4)) // region0_count[gr][ch]
					sideInformation.Region1Count[gr][ch] = byte(bitReader.Read(3)) // region1_count[gr][ch]
				}

				sideInformation.Preflag[gr][ch] = byte(bitReader.Read(1))           // preflag[gr][ch]
				sideInformation.ScalfacScale[gr][ch] = byte(bitReader.Read(1))      // scalefac_scale[gr][ch]
				sideInformation.Count1tableSelect[gr][ch] = byte(bitReader.Read(1)) // count1table_select[gr][ch]
			}
		}
		fmt.Printf("%+v\n", sideInformation)

		// Main Data =================================================================================================
		mainDataLength := frameLength - sideInformationLength - 4
		mainData := make([]byte, mainDataLength)
		file.Read(mainData)
		//fmt.Printf("%x\n", mainData)

		c--
		if c == 0 {
			break
		}
		fmt.Println("==================================================")
	}
}
