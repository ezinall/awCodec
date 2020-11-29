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

	iblen = 576 // Frequency lines of each granule

	maxTableEntry = 15 // Maximum Huffman table entry index
)

type Scalefac struct {
	L [2][2][21]byte
	S [2][2][12][3]byte
}

var scalefacCompress = [16][2]int{
	{0, 0}, {0, 1}, {0, 2}, {0, 3}, {3, 0}, {1, 1}, {1, 2}, {1, 3},
	{2, 1}, {2, 2}, {2, 3}, {3, 1}, {3, 2}, {3, 3}, {4, 2}, {4, 3},
}

var bandIndex = [3][2][]int{
	{ // Layer 3
		{0, 4, 8, 12, 16, 20, 24, 30, 36, 44, 52, 62, 74, 90, 110, 134, 162, 196, 238, 288, 342, 418, 576},
		{0, 4, 8, 12, 16, 22, 30, 40, 52, 66, 84, 106, 136, 192},
	},
	{ // Layer 2
		{0, 4, 8, 12, 16, 20, 24, 30, 36, 42, 50, 60, 72, 88, 106, 128, 156, 190, 230, 276, 330, 384, 576},
		{0, 4, 8, 12, 16, 22, 28, 38, 50, 64, 80, 100, 126, 192},
	},
	{ // Layer 1
		{0, 4, 8, 12, 16, 20, 24, 30, 36, 44, 54, 66, 82, 102, 126, 156, 194, 240, 296, 364, 448, 550, 576},
		{0, 4, 8, 12, 16, 22, 30, 42, 58, 78, 104, 138, 180, 192},
	},
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

		bitReader := NewBitReader(buf)

		sideInfo := SideInformation{}
		sideInfo.MainDataBegin = uint16(bitReader.ReadBits(9)) // main_data_begin

		if header.Mode == modeSingleChannel {
			sideInfo.PrivateBits = byte(bitReader.ReadBits(5)) // private_bits
		} else {
			sideInfo.PrivateBits = byte(bitReader.ReadBits(3)) // private_bits
		}

		for ch := 0; ch < nch; ch++ {
			for band := 0; band < 4; band++ {
				sideInfo.Scfsi[ch][band] = byte(bitReader.ReadBits(1)) // scfsi[ch][scfsi_band]
			}
		}

		for gr := 0; gr < 2; gr++ { // 2 granules for MPEG1, 1 granules for MPEG2
			for ch := 0; ch < nch; ch++ {
				sideInfo.Part23Length[gr][ch] = uint16(bitReader.ReadBits(12))      // part2_3_length[gr][ch]
				sideInfo.BigValues[gr][ch] = uint16(bitReader.ReadBits(9))          // big_values[gr][ch]
				sideInfo.GlobalGain[gr][ch] = uint8(bitReader.ReadBits(8))          // global_gain[gr][ch]
				sideInfo.ScalefacCompress[gr][ch] = byte(bitReader.ReadBits(4))     // scalefac_compress[gr][ch]
				sideInfo.WindowsSwitchingFlag[gr][ch] = byte(bitReader.ReadBits(1)) // window_switching_flag[gr][ch]

				if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 {
					sideInfo.BlockType[gr][ch] = byte(bitReader.ReadBits(2))       // block_type[gr][ch]
					sideInfo.MixedBlockFlag[gr][ch] = uint8(bitReader.ReadBits(1)) // mixed_block_flag[gr][ch]

					for region := 0; region < 2; region++ {
						sideInfo.TableSelect[gr][ch][region] = byte(bitReader.ReadBits(5)) // table_select[gr][ch][region]
					}

					for window := 0; window < 3; window++ {
						sideInfo.SubblockGain[gr][ch][window] = uint8(bitReader.ReadBits(3)) // subblock_gain[gr][ch][window]
					}

					// Set default if window switching set
					blockType := sideInfo.BlockType[gr][ch]
					mixedBlock := sideInfo.MixedBlockFlag[gr][ch]
					if blockType == 1 || blockType == 3 || blockType == 2 && mixedBlock == 1 {
						sideInfo.Region0Count[gr][ch] = 7
					} else if blockType == 2 && mixedBlock != 1 {
						sideInfo.Region0Count[gr][ch] = 8
					}
					sideInfo.Region1Count[gr][ch] = 20 - sideInfo.Region0Count[gr][ch]

				} else {
					// Set default if window not switching
					//sideInfo.BlockType[gr][ch] = 0
					//sideInfo.MixedBlockFlag[gr][ch] = 0

					for region := 0; region < 3; region++ {
						sideInfo.TableSelect[gr][ch][region] = byte(bitReader.ReadBits(5)) // table_select[gr][ch][region]
					}

					sideInfo.Region0Count[gr][ch] = byte(bitReader.ReadBits(4)) // region0_count[gr][ch]
					sideInfo.Region1Count[gr][ch] = byte(bitReader.ReadBits(3)) // region1_count[gr][ch]
				}

				sideInfo.Preflag[gr][ch] = byte(bitReader.ReadBits(1))           // preflag[gr][ch]
				sideInfo.ScalfacScale[gr][ch] = byte(bitReader.ReadBits(1))      // scalefac_scale[gr][ch]
				sideInfo.Count1tableSelect[gr][ch] = byte(bitReader.ReadBits(1)) // count1table_select[gr][ch]
			}
		}

		fmt.Printf("%+v\n", sideInfo)

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

		bitReader = NewBitReader(mainData)

		scalefac := Scalefac{}
		var samples [2][2][iblen]float32

		for gr := 0; gr < 2; gr++ {
			for ch := 0; ch < nch; ch++ {
				bitReader.counter = 0

				slen1 := scalefacCompress[sideInfo.ScalefacCompress[gr][ch]][0]
				slen2 := scalefacCompress[sideInfo.ScalefacCompress[gr][ch]][1]

				//var part2Length int // Number of bits used for scalefactors
				if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == blockShort {
					if sideInfo.MixedBlockFlag[gr][ch] == 1 { // Mixed blocks
						//part2Length = 17*slen1 + 18*slen2 // part2_length all bit length

						for sfb := 0; sfb < 8; sfb++ { // scalefactors bands
							scalefac.L[gr][ch][sfb] = byte(bitReader.ReadBits(slen1))
						}
						for sfb := 3; sfb < 6; sfb++ {
							for window := 0; window < 3; window++ {
								scalefac.S[gr][ch][sfb][window] = byte(bitReader.ReadBits(slen1))
							}
						}

					} else { // Short blocks
						//part2Length = 18*slen1 + 18*slen2 // part2_length all bit length

						for sfb := 0; sfb < 6; sfb++ {
							for window := 0; window < 3; window++ {
								scalefac.S[gr][ch][sfb][window] = byte(bitReader.ReadBits(slen1))
							}
						}
					}

					for sfb := 6; sfb < 12; sfb++ {
						for window := 0; window < 3; window++ {
							scalefac.S[gr][ch][sfb][window] = byte(bitReader.ReadBits(slen2))
						}
					}

				} else { // Long blocks
					//part2Length = 11*slen1 + 10*slen2 // part2_length all bit length

					if gr == 0 {
						for sfb := 0; sfb < 11; sfb++ {
							scalefac.L[gr][ch][sfb] = byte(bitReader.ReadBits(slen1))
						}
						for sfb := 11; sfb < 21; sfb++ {
							scalefac.L[gr][ch][sfb] = byte(bitReader.ReadBits(slen2))
						}

					} else {
						for sfb := 0; sfb < 6; sfb++ {
							if sideInfo.Scfsi[ch][0] == 0 {
								scalefac.L[gr][ch][sfb] = byte(bitReader.ReadBits(slen1))
							} else {
								scalefac.L[gr][ch][sfb] = scalefac.L[0][ch][sfb]
							}
						}
						for sfb := 6; sfb < 11; sfb++ {
							if sideInfo.Scfsi[ch][1] == 0 {
								scalefac.L[gr][ch][sfb] = byte(bitReader.ReadBits(slen1))
							} else {
								scalefac.L[gr][ch][sfb] = scalefac.L[0][ch][sfb]
							}
						}
						for sfb := 11; sfb < 16; sfb++ {
							if sideInfo.Scfsi[ch][2] == 0 {
								scalefac.L[gr][ch][sfb] = byte(bitReader.ReadBits(slen2))
							} else {
								scalefac.L[gr][ch][sfb] = scalefac.L[0][ch][sfb]
							}
						}
						for sfb := 16; sfb < 21; sfb++ {
							if sideInfo.Scfsi[ch][3] == 0 {
								scalefac.L[gr][ch][sfb] = byte(bitReader.ReadBits(slen2))
							} else {
								scalefac.L[gr][ch][sfb] = scalefac.L[0][ch][sfb]
							}
						}
					}
				}

				// Huffman code ===============================================================================================
				part23Length := int(sideInfo.Part23Length[gr][ch]) // int bits

				var region0 int
				var region1 int
				if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == 2 {
					region0 = 36
					region1 = 576
				} else {
					region0 = bandIndex[header.Layer-1][0][sideInfo.Region0Count[gr][ch]+1]
					region1 = bandIndex[header.Layer-1][0][sideInfo.Region0Count[gr][ch]+1+sideInfo.Region1Count[gr][ch]+1]
				}
				fmt.Printf("region0 %+v region1 %+v\n", region0, region1)

				sample := 0
				for ; sample < int(sideInfo.BigValues[gr][ch])*2; sample += 2 {
					tableNum := 0
					if sample < region0 {
						tableNum = int(sideInfo.TableSelect[gr][ch][0])
					} else if sample < region1 {
						tableNum = int(sideInfo.TableSelect[gr][ch][1])
					} else {
						tableNum = int(sideInfo.TableSelect[gr][ch][2])
					}

					if tableNum == 0 {
						continue
					}

					x, y, _, _ := decodeHuffman(bitReader, tableNum)

					samples[gr][ch][sample] = float32(x)
					samples[gr][ch][sample+1] = float32(y)
				}

				//count1 := 0
				for ; sample+4 <= iblen && bitReader.counter < part23Length; sample += 4 {
					//count1++

					var v, w, x, y int
					if sideInfo.Count1tableSelect[gr][ch] == 1 {
						v = bitReader.ReadBits(1) ^ 1
						w = bitReader.ReadBits(1) ^ 1
						x = bitReader.ReadBits(1) ^ 1
						y = bitReader.ReadBits(1) ^ 1
					} else {
						v, w, x, y = decodeHuffmanB(bitReader)
					}

					if v != 0 && bitReader.ReadBits(1) == 1 {
						v = -v
					}
					if w != 0 && bitReader.ReadBits(1) == 1 {
						w = -w
					}
					if x != 0 && bitReader.ReadBits(1) == 1 {
						x = -x
					}
					if y != 0 && bitReader.ReadBits(1) == 1 {
						y = -y
					}

					samples[gr][ch][sample] = float32(v)
					samples[gr][ch][sample+1] = float32(w)
					samples[gr][ch][sample+2] = float32(x)
					samples[gr][ch][sample+3] = float32(y)
				}

				//fmt.Println(sideInfo.BigValues[gr][ch]*2, count1, part23Length, iblen)
			}
		}
		fmt.Printf("%+v\n", scalefac)
		//fmt.Printf("%+v\n", samples)

		fmt.Println("==================================================")
	}
}

func decodeHuffman(r *BitReader, tableNumber int) (x, y, v, w int) {
	table := tables[tableNumber]

	bitSample := r.ReadBits(24)
	for x, v := range table.Table {
		for y, k := range v {
			hcod := k[0]
			hlen := k[1]

			if hcod == bitSample>>(24-hlen) {
				r.Seek(-(24 - hlen))

				if table.Linbits != 0 && x == maxTableEntry {
					x += r.ReadBits(table.Linbits)
				}
				if x != 0 && r.ReadBits(1) == 1 {
					x = -x
				}

				if table.Linbits != 0 && y == maxTableEntry {
					y += r.ReadBits(table.Linbits) //
				}
				if y != 0 && r.ReadBits(1) == 1 {
					y = -y
				}

				return x, y, 0, 0
			}
		}
	}
	r.Seek(-24)
	return x, y, v, w
}

func decodeHuffmanB(r *BitReader) (v, w, x, y int) {
	bitSample := r.ReadBits(24)
	for _, k := range huffmanTableA {
		hcod := k[0]
		hlen := k[1]

		if hcod == bitSample>>(24-hlen) {
			r.Seek(-(24 - hlen))

			v = k[2] & 0x8
			w = k[2] & 0x4
			x = k[2] & 0x2
			y = k[2] & 0x1

			return v, w, x, y
		}
	}
	r.Seek(-24)
	return v, w, x, y
}
