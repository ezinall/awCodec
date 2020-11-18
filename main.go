package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

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

const (
	syncWord = 0xFFF

	// ID
	MPEG1 = 0x1
	MPEG2 = 0x0

	// Layer
	layer1        = 0x3
	layer2        = 0x2
	layer3        = 0x1
	layerReserved = 0x0

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

type SideInformation struct {
	MainDataBegin        uint16         // 9 bits
	PrivateBits          byte           // 3 bits in mono, 5 in stereo
	Scfsi                [2][4]byte     // channels/bands 1 bit
	Part23Length         [2][2]uint16   // granule/channel 12 bits
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

const (
	blockReserved = 0x0 // Reserved block
	blockStart    = 0x1 // Start block
	blockShort    = 0x2 // 3 short block
	blockEnd      = 0x3 // End block
)

type Scalefac struct {
	L [2][2][21]byte
	S [2][2][12][3]byte
}

var scalefacCompress = [16][2]int{
	{0, 0}, {0, 1}, {0, 2}, {0, 3}, {3, 0}, {1, 1}, {1, 2}, {1, 3},
	{2, 1}, {2, 2}, {2, 3}, {3, 1}, {3, 2}, {3, 3}, {4, 2}, {4, 3},
}

type BitReader struct {
	offset uint
	bytes  []byte
}

func (si *BitReader) Bits(n uint) int {
	buf := make([]byte, 3)

	offset := si.offset / 8
	copy(buf, si.bytes[offset:])
	//fmt.Printf("%08b", buf)

	r := int(buf[0])<<16 | int(buf[1])<<8 | int(buf[2])
	r = r >> (uint(24) - n - si.offset%8) & (0xFFFFFF >> (24 - n))

	si.offset += n

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

	var prevData []byte
	for file.Len() != 0 {
		// Header ======================================================================================================
		buf := make([]byte, 4)
		file.Read(buf)
		//fmt.Printf("%08b\n", buf)

		h := int32(buf[0])<<24 | int32(buf[1])<<16 | int32(buf[2])<<8 | int32(buf[3])
		header := Header{
			uint16(h >> 20 & 0xFFF),
			uint8(h >> 19 & 0x1),
			uint8(h >> 17 & 0x3),
			uint8(h >> 16 & 0x1),
			uint8(h >> 12 & 0xF),
			uint8(h >> 10 & 0x3),
			uint8(h >> 9 & 0x1),
			uint8(h >> 8 & 0x1),
			uint8(h >> 6 & 0x3),
			uint8(h >> 4 & 0x3),
			uint8(h >> 3 & 0x1),
			uint8(h >> 2 & 0x1),
			uint8(h >> 0 & 0x3),
		}

		if header.SyncWord != syncWord {
			break
		}
		fmt.Printf("%+v\n", header)

		// Frame length ===============================================================================================
		frameSize := 144 // byte for layer2 and layer3 1152/(1b*8bit) = 144; for layer1 384/(4b*8bit) = 12
		bitrate := bitrateSpecified[header.Layer-1][header.BitrateIndex]
		samplingFrequency := frequencySpecified[header.SamplingFrequency]
		frameLength := frameSize * bitrate / samplingFrequency
		if header.PaddingBit == 1 {
			frameLength += 1
		}
		fmt.Println(bitrate/1000, samplingFrequency, frameLength)

		// CRC Check ===================================================================================================
		if header.ProtectionBit == protected {
			crc := make([]byte, 2)
			file.Read(crc)
		}

		// AudioData ===================================================================================================
		nch := 2 // Number of channels; equals 1 for single_channel mode, equals 2 for other modes.
		if header.Mode == modeSingleChannel {
			nch = 1
		}

		// Side Information ==========================================================================================
		sideInformationLength := 32
		if header.Mode == modeSingleChannel {
			sideInformationLength = 17
		}

		buf = make([]uint8, sideInformationLength)
		file.Read(buf)
		//fmt.Printf("%08b\n", bitReader.bytes)

		bitReader := BitReader{
			0,
			buf,
		}

		sideInfo := SideInformation{}
		sideInfo.MainDataBegin = uint16(bitReader.Bits(9)) // main_data_begin

		if header.Mode == modeSingleChannel {
			sideInfo.PrivateBits = byte(bitReader.Bits(5)) // private_bits
		} else {
			sideInfo.PrivateBits = byte(bitReader.Bits(3)) // private_bits
		}

		for ch := 0; ch < nch; ch++ {
			for band := 0; band < 4; band++ {
				sideInfo.Scfsi[ch][band] = byte(bitReader.Bits(1)) // scfsi[ch][scfsi_band]
			}
		}

		for gr := 0; gr < 2; gr++ { // 2 granules for MPEG1, 1 granules for MPEG2
			for ch := 0; ch < nch; ch++ {
				sideInfo.Part23Length[gr][ch] = uint16(bitReader.Bits(12))      // part2_3_length[gr][ch]
				sideInfo.BigValues[gr][ch] = uint16(bitReader.Bits(9))          // big_values[gr][ch]
				sideInfo.GlobalGain[gr][ch] = uint8(bitReader.Bits(8))          // global_gain[gr][ch]
				sideInfo.ScalefacCompress[gr][ch] = byte(bitReader.Bits(4))     // scalefac_compress[gr][ch]
				sideInfo.WindowsSwitchingFlag[gr][ch] = byte(bitReader.Bits(1)) // window_switching_flag[gr][ch]

				if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 {
					sideInfo.BlockType[gr][ch] = byte(bitReader.Bits(2))       // block_type[gr][ch]
					sideInfo.MixedBlockFlag[gr][ch] = uint8(bitReader.Bits(1)) // mixed_block_flag[gr][ch]

					for region := 0; region < 2; region++ {
						sideInfo.TableSelect[gr][ch][region] = byte(bitReader.Bits(5)) // table_select[gr][ch][region]
					}

					for window := 0; window < 3; window++ {
						sideInfo.SubblockGain[gr][ch][window] = uint8(bitReader.Bits(3)) // subblock_gain[gr][ch][window]
					}

					// TODO Clarify default values
					blockType := sideInfo.BlockType[gr][ch]
					mixedBlock := sideInfo.MixedBlockFlag[gr][ch]
					if blockType == 1 || blockType == 3 || blockType == 2 && mixedBlock == 1 {
						sideInfo.Region0Count[gr][ch] = 7 // 8
					} else if blockType == 2 && mixedBlock != 1 {
						sideInfo.Region0Count[gr][ch] = 8 // 9
					}
					sideInfo.Region1Count[gr][ch] = 36 // 63 or 0

				} else {
					for region := 0; region < 3; region++ {
						sideInfo.TableSelect[gr][ch][region] = byte(bitReader.Bits(5)) // table_select[gr][ch][region]
					}

					sideInfo.Region0Count[gr][ch] = byte(bitReader.Bits(4)) // region0_count[gr][ch]
					sideInfo.Region1Count[gr][ch] = byte(bitReader.Bits(3)) // region1_count[gr][ch]
				}

				sideInfo.Preflag[gr][ch] = byte(bitReader.Bits(1))           // preflag[gr][ch]
				sideInfo.ScalfacScale[gr][ch] = byte(bitReader.Bits(1))      // scalefac_scale[gr][ch]
				sideInfo.Count1tableSelect[gr][ch] = byte(bitReader.Bits(1)) // count1table_select[gr][ch]
			}
		}

		fmt.Printf("%+v, %d\n", sideInfo, bitReader.offset)

		// Main Data ==================================================================================================
		mainDataLength := frameLength - sideInformationLength - 4 // 4 bytes header
		if header.ProtectionBit == protected {
			mainDataLength -= 2
		}

		mainData := make([]byte, mainDataLength)
		file.Read(mainData)

		if sideInfo.MainDataBegin != 0 {
			mainData = append(prevData[len(prevData)-int(sideInfo.MainDataBegin):], mainData...)
		}
		//fmt.Printf("%d %x\n", len(mainData), mainData)

		prevData = mainData

		// TODO refactor reader?
		bitReader = BitReader{
			0,
			mainData,
		}

		scalefac := Scalefac{}

		for gr := 0; gr < 2; gr++ {
			for ch := 0; ch < nch; ch++ {
				slen1 := scalefacCompress[sideInfo.ScalefacCompress[gr][ch]][0]
				slen2 := scalefacCompress[sideInfo.ScalefacCompress[gr][ch]][1]

				var part2Length int
				if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == blockShort {
					if sideInfo.MixedBlockFlag[gr][ch] == 1 {
						part2Length = 17*slen1 + 18*slen2 // part2_length all bit length

						for sfb := 0; sfb < 8; sfb++ { // scalefactors bands
							scalefac.L[gr][ch][sfb] = byte(bitReader.Bits(uint(slen1))) // TODO BitReader
						}
						for sfb := 3; sfb < 6; sfb++ {
							for window := 0; window < 3; window++ {
								scalefac.S[gr][ch][sfb][window] = byte(bitReader.Bits(uint(slen1))) // TODO BitReader
							}
						}
						for sfb := 6; sfb < 12; sfb++ {
							for window := 0; window < 3; window++ {
								scalefac.S[gr][ch][sfb][window] = byte(bitReader.Bits(uint(slen2))) // TODO BitReader
							}
						}
					} else {
						part2Length = 18*slen1 + 18*slen2 // part2_length all bit length

						for sfb := 0; sfb < 6; sfb++ {
							for window := 0; window < 3; window++ {
								scalefac.S[gr][ch][sfb][window] = byte(bitReader.Bits(uint(slen1))) // TODO BitReader
							}
						}
						for sfb := 6; sfb < 12; sfb++ {
							for window := 0; window < 3; window++ {
								scalefac.S[gr][ch][sfb][window] = byte(bitReader.Bits(uint(slen2))) // TODO BitReader
							}
						}
					}
				} else {
					part2Length = 11*slen1 + 10*slen2 // part2_length all bit length

					if sideInfo.Scfsi[ch][0] == 0 || gr == 0 {
						for sfb := 0; sfb < 6; sfb++ {
							scalefac.L[gr][ch][sfb] = byte(bitReader.Bits(uint(slen1))) // TODO BitReader
						}
					}
					if sideInfo.Scfsi[ch][1] == 0 || gr == 0 {
						for sfb := 6; sfb < 11; sfb++ {
							scalefac.L[gr][ch][sfb] = byte(bitReader.Bits(uint(slen1))) // TODO BitReader
						}
					}
					if sideInfo.Scfsi[ch][2] == 0 || gr == 0 {
						for sfb := 11; sfb < 16; sfb++ {
							scalefac.L[gr][ch][sfb] = byte(bitReader.Bits(uint(slen2))) // TODO BitReader
						}
					}
					if sideInfo.Scfsi[ch][3] == 0 || gr == 0 {
						for sfb := 16; sfb < 21; sfb++ {
							scalefac.L[gr][ch][sfb] = byte(bitReader.Bits(uint(slen2))) // TODO BitReader
						}
					}
				}

				// Huffman code ===============================================================================================
				//bit_pos_end := part2Length + int(sideInfo.Part23Length[gr][ch]) - 1
				huffmanCodeLength := int(sideInfo.Part23Length[gr][ch]) - part2Length // int bits
				fmt.Printf("bit %d %d", part2Length, huffmanCodeLength)

				for i := 0; i <= huffmanCodeLength; i++ {
					bitReader.Bits(1)
				}
				fmt.Println()
			}
		}
		fmt.Printf("%+v\n", scalefac)
		fmt.Println("==================================================")
	}
}
