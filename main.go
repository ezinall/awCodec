package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
)

type Header struct {
	SyncWord          uint16 // 12 bits
	ID                byte   // 1 bit
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
	// Header.SyncWord
	syncWord = 0xFFF

	// Header.ID
	MPEG1 = 0b1
	MPEG2 = 0b0

	// Header.Layer
	layer1        = 0b11
	layer2        = 0b10
	layer3        = 0b01
	layerReserved = 0b00

	// Header.ProtectionBit
	protected = 0b0

	// Header.Mode specified
	modeStereo        = 0b00
	modeJoinStereo    = 0b01 // intensity_stereo and/or ms_stereo
	modeDualChannel   = 0b10
	modeSingleChannel = 0b11

	// Header.ModeExtension specifies
	intensityStereo = 0b01
	msStereo        = 0b10
)

var bitrateSpecified = [3][16]int{
	// Layer 3
	{0, 32000, 40000, 48000, 56000, 64000, 80000, 96000,
		112000, 128000, 160000, 192000, 224000, 256000, 320000, 0}, // Bit/s
	// Layer 2
	{0, 32000, 48000, 56000, 64000, 80000, 96000, 112000,
		128000, 160000, 192000, 224000, 256000, 320000, 384000, 0}, // Bit/s
	// Layer 1
	{0, 32000, 64000, 96000, 128000, 160000, 192000, 224000,
		256000, 288000, 320000, 352000, 384000, 416000, 448000, 0}, // Bit/s
}

var frequencySpecified = [4]int{
	44100, 48000, 32000, 0, // Hz
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
	// SideInformation.BlockType
	blockReserved = 0b00 // Reserved block
	blockStart    = 0b01 // Start block
	blockShort    = 0b10 // 3 short block
	blockEnd      = 0b11 // End block

	mixed = 0b1 // Mixed flag
)

const (
	iblen = 576 // Frequency lines of each granule
)

type Scalefac struct {
	L [2][2][21]byte
	S [2][2][12][3]byte // gr/ch/sfb/window
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

// If it is set, then the values of the Table 11 are added to the scale factors.
//Preflag table only for SideInformation.BlockType blockShort windows.
var pretab = [21]int{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 3, 3, 3, 2,
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	content, err := ioutil.ReadFile("file_example_MP3_700KB.mp3")
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
		//fmt.Printf("%+v\n", header)

		// Frame length ===============================================================================================
		frameSize := 144 // byte for layer2 and layer3 1152/(1b*8bit) = 144; for layer1 384/(4b*8bit) = 12
		bitrate := bitrateSpecified[header.Layer-1][header.BitrateIndex]
		samplingFrequency := frequencySpecified[header.SamplingFrequency]
		frameLength := frameSize * bitrate / samplingFrequency
		if header.PaddingBit == 1 {
			frameLength += 1
		}
		//fmt.Println(bitrate/1000, samplingFrequency, frameLength)

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

		//fmt.Printf("%+v\n", sideInfo)

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
		var is [2][2][iblen]float32
		var count1 [2][2]int

		for gr := 0; gr < 2; gr++ {
			for ch := 0; ch < nch; ch++ {
				bitReader.counter = 0

				slen1 := scalefacCompress[sideInfo.ScalefacCompress[gr][ch]][0]
				slen2 := scalefacCompress[sideInfo.ScalefacCompress[gr][ch]][1]

				// Scalefactor
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

				// Huffman code =======================================================================================
				var region0 int
				var region1 int
				if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == blockShort {
					region0 = 36
					region1 = 576
				} else {
					region0 = bandIndex[header.Layer-1][0][sideInfo.Region0Count[gr][ch]+1]
					region1 = bandIndex[header.Layer-1][0][sideInfo.Region0Count[gr][ch]+1+sideInfo.Region1Count[gr][ch]+1]
				}
				//fmt.Printf("region0 %+v region1 %+v\n", region0, region1)

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

					is[gr][ch][sample] = float32(x)
					is[gr][ch][sample+1] = float32(y)
				}

				for ; sample+4 < iblen && bitReader.counter < int(sideInfo.Part23Length[gr][ch]); sample += 4 {
					count1[gr][ch] += 4

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

					is[gr][ch][sample] = float32(v)
					is[gr][ch][sample+1] = float32(w)
					is[gr][ch][sample+2] = float32(x)
					is[gr][ch][sample+3] = float32(y)
				}

				//fmt.Println(sideInfo.BigValues[gr][ch]*2, count1, part23Length, iblen)
			}
		}

		// ??? ================================================================================================
		for gr := 0; gr < 2; gr++ {
			for ch := 0; ch < nch; ch++ {
				requantize(gr, ch, header, sideInfo, scalefac, &is)
				reorder(gr, ch, header, sideInfo, &is, count1)
			}
			stereo(gr, header, &is)
			for ch := 0; ch < nch; ch++ {
				aliasReduction(gr, ch, sideInfo, &is)
				imdct(gr, ch, sideInfo.BlockType[gr][ch], &is)
				frequencyInversion()
				synthFilterbank()
			}
		}

		//fmt.Printf("%+v\n", scalefac)
		//fmt.Printf("%+v\n", is)
		//fmt.Println("==================================================")
	}
}

func requantize(gr, ch int, header Header, sideInfo SideInformation, scalefac Scalefac, is *[2][2][iblen]float32) {
	var (
		A, B float64
	)
	sfb := 0
	window := 0

	scalefacMultiplier := 0.5
	if sideInfo.ScalfacScale[gr][ch] == 0 {
		scalefacMultiplier = 1
	}

	for i, sample := 0, 0; sample < iblen; i, sample = i+1, sample+1 {
		if sideInfo.BlockType[gr][ch] == blockShort || sideInfo.MixedBlockFlag[gr][ch] == 1 && sfb >= 8 { // Short blocks
			if i == bandIndex[header.Layer-1][1][sfb] {
				i = 0
				if window == 2 {
					window = 0
					sfb++
				} else {
					window++
				}
			}

			A = float64(sideInfo.GlobalGain[gr][ch] - 210 - 8*sideInfo.SubblockGain[gr][ch][window])
			B = scalefacMultiplier * float64(scalefac.S[gr][ch][sfb][window])

		} else { // Long blocks
			if sample == bandIndex[header.Layer-1][0][sfb+2] { // sfb+1 ВОЗМОЖНО ОШИБКА!!!
				sfb++
			}

			A = float64(sideInfo.GlobalGain[gr][ch]) - 210
			B = scalefacMultiplier * float64(scalefac.L[gr][ch][sfb]+sideInfo.Preflag[gr][ch]*byte(pretab[sfb]))
		}

		sign := math.Copysign(1, float64(is[gr][ch][sample]))
		is[gr][ch][sample] = float32(sign * math.Pow(math.Abs(float64(is[gr][ch][sample])), 4./3.) *
			math.Pow(2, A/4.) * math.Pow(2, -B))

	}
}

func reorder(gr, ch int, header Header, sideInfo SideInformation, is *[2][2][iblen]float32, count1 [2][2]int) {
	samplesBuf := make([]float32, iblen)
	band := bandIndex[header.Layer-1][1]

	// Only reorder short blocks
	if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == blockShort {
		sfb := 0
		if sideInfo.MixedBlockFlag[gr][ch] != 0 {
			sfb = 3
		}

		nextSfb := band[sfb+1] * 3
		windowLen := band[sfb+1] - band[sfb]

		i := 0
		if sfb == 0 {
			i = 0
		}

		for i < iblen {
			if i == nextSfb {
				j := band[sfb] * 3
				copy(is[gr][ch][j:j+3*windowLen], samplesBuf[:3*windowLen])

				if i >= count1[gr][ch] {
					return
				}

				sfb++
				nextSfb = band[sfb+1] * 3
				windowLen = band[sfb+1] - band[sfb]
			}
			for window := 0; window < 3; window++ {
				for j := 0; j < windowLen; j++ {
					samplesBuf[j*3+window] = is[gr][ch][i]
					i++
				}
			}
		}
		j := 3 * band[12]
		copy(is[gr][ch][j:j+3*windowLen], samplesBuf[:3*windowLen])
	}
}

func stereo(gr int, header Header, is *[2][2][iblen]float32) {
	if header.Mode == modeJoinStereo {
		if header.ModeExtension&intensityStereo == intensityStereo {
			// TODO
		}
		if header.ModeExtension&msStereo == msStereo {
			for sample := 0; sample < iblen; sample++ {
				mid := is[gr][0][sample]
				side := is[gr][1][sample]
				is[gr][0][sample] = (mid + side) / math.Sqrt2
				is[gr][1][sample] = (mid - side) / math.Sqrt2
			}
		}
	}
}

func aliasReduction(gr, ch int, sideInfo SideInformation, is *[2][2][iblen]float32) {
	if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == blockShort && sideInfo.MixedBlockFlag[gr][ch] == 0 {
		return
	}

	nsb := 32
	if sideInfo.MixedBlockFlag[gr][ch] == 1 {
		nsb = 2
	}

	for sb := 1; sb < nsb; sb++ {
		for i := 0; i < 8; i++ {
			//xar[18*sb-1-i] = xr[18*sb-1-i]Cs[i] - xr[18*sb+i]Ca[i]
			//xar[18*sb+i] = xr[18*sb+i]Cs[i] + xr[18*sb-1-i]Ca[i]
			li := 18*sb - 1 - i
			ui := 18*sb + i
			sample1 := is[gr][ch][li]
			sample2 := is[gr][ch][ui]
			is[gr][ch][li] = sample1*cs[i] - sample2*cs[i]
			is[gr][ch][ui] = sample2*cs[i] + sample1*cs[i]
		}
	}
}

func imdct(gr, ch int, blockType byte, is *[2][2][iblen]float32) {
	n := 36
	if blockType == blockShort {
		n = 12
	}
	halfN := n / 2
	samplesBlock := make([]float32, 36)

	nwin := 1
	if blockType == blockShort {
		nwin = 3
	}

	for block := 0; block < 32; block++ {
		for window := 0; window < nwin; window++ {
			for i := 0; i < n; i++ {
				xi := float32(0.)
				for k := 0; k < halfN; k++ {
					s := is[gr][ch][block*18+halfN*window+k]
					xi += s * float32(math.Cos(math.Pi/float64(2*n)*float64(2*i+1+halfN)*float64(2*k+1)))
				}
				samplesBlock[window*n+i] = xi * imdctWinData[blockType][i]
			}
		}
		// TODO
	}
}

func frequencyInversion() {

}

func synthFilterbank() {

}
