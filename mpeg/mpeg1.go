package mpeg

import (
	"awCodec/id3"
	"awCodec/pcm"
	"awCodec/utils"
	"bytes"
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
	mpeg1 = 0b1
	mpeg2 = 0b0

	// Header.Layer
	layer1        = 0b11
	layer2        = 0b10
	layer3        = 0b01
	layerReserved = 0b00

	// Header.ProtectionBit
	protected = 0b0

	// Header.PaddingBit
	padding = 0b1

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

const iblen = 576 // Frequency lines of each granule

type sideInformation struct {
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
	// sideInformation.BlockType
	blockReserved = 0b00 // Reserved block
	blockStart    = 0b01 // Start block
	blockShort    = 0b10 // 3 short block
	blockEnd      = 0b11 // End block
)

type Scalefac struct {
	L [2][2][22]byte
	S [2][2][13][3]byte // gr/ch/sfb/window
}

var scalefacCompress = [16][2]int{
	{0, 0}, {0, 1}, {0, 2}, {0, 3}, {3, 0}, {1, 1}, {1, 2}, {1, 3},
	{2, 1}, {2, 2}, {2, 3}, {3, 1}, {3, 2}, {3, 3}, {4, 2}, {4, 3},
}

var bandIndex = [3][2][]int{
	{
		{0, 4, 8, 12, 16, 20, 24, 30, 36, 44, 52, 62, 74, 90, 110, 134, 162, 196, 238, 288, 342, 418, 576},
		{0, 4, 8, 12, 16, 22, 30, 40, 52, 66, 84, 106, 136, 192},
	},
	{
		{0, 4, 8, 12, 16, 20, 24, 30, 36, 42, 50, 60, 72, 88, 106, 128, 156, 190, 230, 276, 330, 384, 576},
		{0, 4, 8, 12, 16, 22, 28, 38, 50, 64, 80, 100, 126, 192},
	},
	{
		{0, 4, 8, 12, 16, 20, 24, 30, 36, 44, 54, 66, 82, 102, 126, 156, 194, 240, 296, 364, 448, 550, 576},
		{0, 4, 8, 12, 16, 22, 30, 42, 58, 78, 104, 138, 180, 192},
	},
}

// If it is set, then the values of the Table 11 are added to the scale factors.
//Preflag table only for sideInformation.BlockType blockShort windows.
var pretab = [22]int{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 3, 3, 3, 2, 0,
}

// Decode MPEG1/MPEG2 format.
func mpeg(file *bytes.Reader) (*pcm.F32LE, error) {
	var out = &pcm.F32LE{}

	id3.ReadID3(file)

	var prevData []byte
	prevSamples := [2][32][18]float32{}
	vVec := [2][1024]float32{}
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
		if header.PaddingBit == padding {
			frameLength += 1
		}
		//fmt.Println(bitrate/1000, samplingFrequency, frameLength)

		if out.Context().SampleRate == 0 {
			out.Context().SampleRate = samplingFrequency
		}

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

		if out.Context().Channels == 0 {
			out.Context().Channels = nch
		}

		// Side Information ==========================================================================================
		sideInformationLength := 32
		if header.Mode == modeSingleChannel {
			sideInformationLength = 17
		}

		buf = make([]uint8, sideInformationLength)
		file.Read(buf)
		//fmt.Printf("%08b\n", bitReader.bytes)

		sideInfoBitReader := utils.NewBitReader(buf)

		sideInfo := sideInformation{}
		sideInfo.MainDataBegin = uint16(sideInfoBitReader.ReadBits(9)) // main_data_begin

		if header.Mode == modeSingleChannel {
			sideInfo.PrivateBits = byte(sideInfoBitReader.ReadBits(5)) // private_bits
		} else {
			sideInfo.PrivateBits = byte(sideInfoBitReader.ReadBits(3)) // private_bits
		}

		for ch := 0; ch < nch; ch++ {
			for band := 0; band < 4; band++ {
				sideInfo.Scfsi[ch][band] = byte(sideInfoBitReader.ReadBits(1)) // scfsi[ch][scfsi_band]
			}
		}

		for gr := 0; gr < 2; gr++ { // 2 granules for MPEG1, 1 granules for MPEG2
			for ch := 0; ch < nch; ch++ {
				sideInfo.Part23Length[gr][ch] = uint16(sideInfoBitReader.ReadBits(12))      // part2_3_length[gr][ch]
				sideInfo.BigValues[gr][ch] = uint16(sideInfoBitReader.ReadBits(9))          // big_values[gr][ch]
				sideInfo.GlobalGain[gr][ch] = uint8(sideInfoBitReader.ReadBits(8))          // global_gain[gr][ch]
				sideInfo.ScalefacCompress[gr][ch] = byte(sideInfoBitReader.ReadBits(4))     // scalefac_compress[gr][ch]
				sideInfo.WindowsSwitchingFlag[gr][ch] = byte(sideInfoBitReader.ReadBits(1)) // window_switching_flag[gr][ch]

				if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 {
					sideInfo.BlockType[gr][ch] = byte(sideInfoBitReader.ReadBits(2))       // block_type[gr][ch]
					sideInfo.MixedBlockFlag[gr][ch] = uint8(sideInfoBitReader.ReadBits(1)) // mixed_block_flag[gr][ch]

					for region := 0; region < 2; region++ {
						sideInfo.TableSelect[gr][ch][region] = byte(sideInfoBitReader.ReadBits(5)) // table_select[gr][ch][region]
					}

					for window := 0; window < 3; window++ {
						sideInfo.SubblockGain[gr][ch][window] = uint8(sideInfoBitReader.ReadBits(3)) // subblock_gain[gr][ch][window]
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
						sideInfo.TableSelect[gr][ch][region] = byte(sideInfoBitReader.ReadBits(5)) // table_select[gr][ch][region]
					}

					sideInfo.Region0Count[gr][ch] = byte(sideInfoBitReader.ReadBits(4)) // region0_count[gr][ch]
					sideInfo.Region1Count[gr][ch] = byte(sideInfoBitReader.ReadBits(3)) // region1_count[gr][ch]
				}

				sideInfo.Preflag[gr][ch] = byte(sideInfoBitReader.ReadBits(1))           // preflag[gr][ch]
				sideInfo.ScalfacScale[gr][ch] = byte(sideInfoBitReader.ReadBits(1))      // scalefac_scale[gr][ch]
				sideInfo.Count1tableSelect[gr][ch] = byte(sideInfoBitReader.ReadBits(1)) // count1table_select[gr][ch]
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

		mainDataBitReader := utils.NewBitReader(mainData)

		scalefac := Scalefac{}
		is := [2][2][iblen]float32{}
		countValues := [2][2]int{}

		for gr := 0; gr < 2; gr++ {
			for ch := 0; ch < nch; ch++ {
				mainDataBitReader.Counter = 0

				slen1 := scalefacCompress[sideInfo.ScalefacCompress[gr][ch]][0]
				slen2 := scalefacCompress[sideInfo.ScalefacCompress[gr][ch]][1]

				// Scalefactor ========================================================================================
				//var part2Length int // Number of bits used for scalefactors
				if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == blockShort {
					if sideInfo.MixedBlockFlag[gr][ch] == 1 { // Mixed blocks
						//part2Length = 17*slen1 + 18*slen2 // part2_length all bit length

						for sfb := 0; sfb < 8; sfb++ { // scalefactors bands
							scalefac.L[gr][ch][sfb] = byte(mainDataBitReader.ReadBits(slen1))
						}
						for sfb := 3; sfb < 6; sfb++ {
							for window := 0; window < 3; window++ {
								scalefac.S[gr][ch][sfb][window] = byte(mainDataBitReader.ReadBits(slen1))
							}
						}

					} else { // Short blocks
						//part2Length = 18*slen1 + 18*slen2 // part2_length all bit length

						for sfb := 0; sfb < 6; sfb++ {
							for window := 0; window < 3; window++ {
								scalefac.S[gr][ch][sfb][window] = byte(mainDataBitReader.ReadBits(slen1))
							}
						}
					}

					for sfb := 6; sfb < 12; sfb++ {
						for window := 0; window < 3; window++ {
							scalefac.S[gr][ch][sfb][window] = byte(mainDataBitReader.ReadBits(slen2))
						}
					}

				} else { // Long blocks
					//part2Length = 11*slen1 + 10*slen2 // part2_length all bit length

					if gr == 0 {
						for sfb := 0; sfb < 11; sfb++ {
							scalefac.L[gr][ch][sfb] = byte(mainDataBitReader.ReadBits(slen1))
						}
						for sfb := 11; sfb < 21; sfb++ {
							scalefac.L[gr][ch][sfb] = byte(mainDataBitReader.ReadBits(slen2))
						}

					} else {
						for sfb := 0; sfb < 6; sfb++ {
							if sideInfo.Scfsi[ch][0] == 0 {
								scalefac.L[gr][ch][sfb] = byte(mainDataBitReader.ReadBits(slen1))
							} else {
								scalefac.L[gr][ch][sfb] = scalefac.L[0][ch][sfb]
							}
						}
						for sfb := 6; sfb < 11; sfb++ {
							if sideInfo.Scfsi[ch][1] == 0 {
								scalefac.L[gr][ch][sfb] = byte(mainDataBitReader.ReadBits(slen1))
							} else {
								scalefac.L[gr][ch][sfb] = scalefac.L[0][ch][sfb]
							}
						}
						for sfb := 11; sfb < 16; sfb++ {
							if sideInfo.Scfsi[ch][2] == 0 {
								scalefac.L[gr][ch][sfb] = byte(mainDataBitReader.ReadBits(slen2))
							} else {
								scalefac.L[gr][ch][sfb] = scalefac.L[0][ch][sfb]
							}
						}
						for sfb := 16; sfb < 21; sfb++ {
							if sideInfo.Scfsi[ch][3] == 0 {
								scalefac.L[gr][ch][sfb] = byte(mainDataBitReader.ReadBits(slen2))
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
					region1 = iblen
				} else {
					region0 = bandIndex[header.SamplingFrequency][0][sideInfo.Region0Count[gr][ch]+1]
					region1 = bandIndex[header.SamplingFrequency][0][sideInfo.Region0Count[gr][ch]+1+sideInfo.Region1Count[gr][ch]+1]
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

					x, y := decodeHuffman(mainDataBitReader, tableNum)
					is[gr][ch][sample] = float32(x)
					is[gr][ch][sample+1] = float32(y)
				}

				count1 := 0
				for ; sample+4 <= iblen && mainDataBitReader.Counter < int(sideInfo.Part23Length[gr][ch]); sample += 4 {
					count1++
					var v, w, x, y int
					if sideInfo.Count1tableSelect[gr][ch] == 1 {
						v, w, x, y = decodeHuffmanB(mainDataBitReader)
					} else {
						v, w, x, y = decodeHuffmanA(mainDataBitReader)
					}

					is[gr][ch][sample] = float32(v)
					is[gr][ch][sample+1] = float32(w)
					is[gr][ch][sample+2] = float32(x)
					is[gr][ch][sample+3] = float32(y)
				}

				countValues[gr][ch] = int(sideInfo.BigValues[gr][ch])*2 + count1*4

				//fmt.Println(sideInfo.BigValues[gr][ch]*2, count1*4, part23Length, iblen)
			}
		}

		// ??? ================================================================================================
		pcm := make([]float32, iblen*2*2) // iblen * number granules * byte count per sample
		for gr := 0; gr < 2; gr++ {
			for ch := 0; ch < nch; ch++ {
				requantize(gr, ch, header, sideInfo, scalefac, &is, countValues)
				reorder(gr, ch, header, sideInfo, &is, countValues)
			}
			stereo(gr, header, &is)
			for ch := 0; ch < nch; ch++ {
				aliasReduction(gr, ch, sideInfo, &is)
				imdct(gr, ch, sideInfo.BlockType[gr][ch], &is, &prevSamples)
				frequencyInversion(gr, ch, &is)
				synthFilterbank(gr, ch, &is, &vVec, pcm[iblen*gr*2:])
			}
		}
		//out.SetPcm(append(out.Pcm(), pcm...))
		out.Append(pcm)
	}

	return out, nil
}

func requantize(gr, ch int, header Header, sideInfo sideInformation, scalefac Scalefac, is *[2][2][iblen]float32, countValues [2][2]int) {
	var (
		A, B float64
	)
	//sfb := 0
	//window := 0

	scalefacMultiplier := 0.5
	if sideInfo.ScalfacScale[gr][ch] == 1 {
		scalefacMultiplier = 1.0
	}

	if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == blockShort { // Short blocks
		if sideInfo.MixedBlockFlag[gr][ch] == 1 { // 2 long sb first
			sfb := 0
			nextSfb := bandIndex[header.SamplingFrequency][0][sfb+1]
			for i := 0; i < 36; i++ {
				if i == nextSfb {
					sfb++
					nextSfb = bandIndex[header.SamplingFrequency][0][sfb+1]
				}
				A = float64(sideInfo.GlobalGain[gr][ch]) - 210
				B = scalefacMultiplier * float64(int(scalefac.L[gr][ch][sfb])+int(sideInfo.Preflag[gr][ch])*pretab[sfb])

				sign := math.Copysign(1, float64(is[gr][ch][i]))
				is[gr][ch][i] = float32(sign * math.Pow(math.Abs(float64(is[gr][ch][i])), 4.0/3.0) *
					math.Pow(2, A/4.0) * math.Pow(2, -B))
			}
			sfb = 3
			nextSfb = bandIndex[header.SamplingFrequency][1][sfb+1] * 3
			windowLen := bandIndex[header.SamplingFrequency][1][sfb+1] - bandIndex[header.SamplingFrequency][1][sfb]
			for i := 36; i < countValues[gr][ch]; {
				if i == nextSfb {
					sfb++
					nextSfb = bandIndex[header.SamplingFrequency][1][sfb+1] * 3
					windowLen = bandIndex[header.SamplingFrequency][1][sfb+1] - bandIndex[header.SamplingFrequency][1][sfb]
				}
				for window := 0; window < 3; window++ {
					for j := 0; j < windowLen; j++ {

						A = float64(sideInfo.GlobalGain[gr][ch]) - 210.0 - 8.0*float64(sideInfo.SubblockGain[gr][ch][window])
						B = scalefacMultiplier * float64(scalefac.S[gr][ch][sfb][window])
						C := A/4.0 + -B

						sign := math.Copysign(1, float64(is[gr][ch][i]))
						is[gr][ch][i] = float32(sign * math.Pow(math.Abs(float64(is[gr][ch][i])), 4.0/3.0) *
							math.Pow(2, C))

						i++
					}
				}
			}
		} else { // Only short blocks
			sfb := 0
			nextSfb := bandIndex[header.SamplingFrequency][1][sfb+1] * 3
			windowLen := bandIndex[header.SamplingFrequency][1][sfb+1] - bandIndex[header.SamplingFrequency][1][sfb]
			for i := 0; i < countValues[gr][ch]; {
				if i == nextSfb {
					sfb++
					nextSfb = bandIndex[header.SamplingFrequency][1][sfb+1] * 3
					windowLen = bandIndex[header.SamplingFrequency][1][sfb+1] - bandIndex[header.SamplingFrequency][1][sfb]
				}
				for window := 0; window < 3; window++ {
					for j := 0; j < windowLen; j++ {
						A = float64(sideInfo.GlobalGain[gr][ch]) - 210.0 - 8.0*float64(sideInfo.SubblockGain[gr][ch][window])
						B = scalefacMultiplier * float64(scalefac.S[gr][ch][sfb][window])
						C := A/4.0 + -B

						sign := math.Copysign(1, float64(is[gr][ch][i]))
						is[gr][ch][i] = float32(sign * math.Pow(math.Abs(float64(is[gr][ch][i])), 4.0/3.0) *
							math.Pow(2, C))
						i++
					}
				}
			}
		}
	} else { // Only long blocks
		sfb := 0
		nextSfb := bandIndex[header.SamplingFrequency][0][sfb+1]
		for i := 0; i < countValues[gr][ch]; i++ {
			if i == nextSfb {
				sfb++
				nextSfb = bandIndex[header.SamplingFrequency][0][sfb+1]
			}
			A = float64(sideInfo.GlobalGain[gr][ch]) - 210
			B = scalefacMultiplier * float64(int(scalefac.L[gr][ch][sfb])+int(sideInfo.Preflag[gr][ch])*pretab[sfb])

			sign := math.Copysign(1, float64(is[gr][ch][i]))
			is[gr][ch][i] = float32(sign * math.Pow(math.Abs(float64(is[gr][ch][i])), 4.0/3.0) *
				math.Pow(2, A/4.0) * math.Pow(2, -B))
		}
	}

	//for sample, i := 0, 0; sample < iblen; sample, i = sample+1, i+1 {
	//	if sideInfo.BlockType[gr][ch] == blockShort || sideInfo.MixedBlockFlag[gr][ch] == 1 && sfb >= 8 { // Short blocks
	//		if i == bandIndex[header.SamplingFrequency][1][sfb] {
	//			i = 0
	//			if window == 2 {
	//				window = 0
	//				sfb++
	//			} else {
	//				window++
	//			}
	//		}
	//
	//		A = float64(sideInfo.GlobalGain[gr][ch]) - 210.0 - 8.0*float64(sideInfo.SubblockGain[gr][ch][window])
	//		B = scalefacMultiplier * float64(scalefac.S[gr][ch][sfb][window])
	//
	//	} else { // Long blocks
	//		if sample == bandIndex[header.SamplingFrequency][0][sfb+1] {
	//			sfb++
	//		}
	//
	//		A = float64(sideInfo.GlobalGain[gr][ch]) - 210.0
	//		B = scalefacMultiplier * float64(scalefac.L[gr][ch][sfb]+sideInfo.Preflag[gr][ch]*byte(pretab[sfb]))
	//	}
	//
	//	sign := math.Copysign(1, float64(is[gr][ch][sample]))
	//	C := A/4.0 + -B
	//	//fmt.Println(C)
	//	is[gr][ch][sample] = float32(sign * math.Pow(math.Abs(float64(is[gr][ch][sample])), 4.0/3.0) *
	//		math.Pow(2, C))
	//}
}

func reorder(gr, ch int, header Header, sideInfo sideInformation, is *[2][2][iblen]float32, countValues [2][2]int) {
	samplesBuf := make([]float32, iblen)
	shortBand := bandIndex[header.SamplingFrequency][1]

	// Only reorder short blocks
	if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == blockShort {
		sfb := 0
		if sideInfo.MixedBlockFlag[gr][ch] == 1 {
			sfb = 3
		}

		nextSfb := shortBand[sfb+1] * 3
		windowLen := shortBand[sfb+1] - shortBand[sfb]

		i := 36
		if sfb == 0 {
			i = 0
		}

		for i < iblen {
			if i == nextSfb {
				j := shortBand[sfb] * 3
				copy(is[gr][ch][j:j+3*windowLen], samplesBuf[0:3*windowLen])

				if i >= countValues[gr][ch] {
					return
				}

				sfb++
				nextSfb = shortBand[sfb+1] * 3
				windowLen = shortBand[sfb+1] - shortBand[sfb]
			}
			for window := 0; window < 3; window++ {
				for j := 0; j < windowLen; j++ {
					samplesBuf[j*3+window] = is[gr][ch][i]
					i++
				}
			}
		}
		j := 3 * shortBand[12]
		copy(is[gr][ch][j:j+3*windowLen], samplesBuf[0:3*windowLen])
	}
}

func stereo(gr int, header Header, is *[2][2][iblen]float32) {
	if header.Mode == modeJoinStereo {
		if header.ModeExtension&intensityStereo == intensityStereo {
			// TODO
		}
		if header.ModeExtension&msStereo == msStereo {
			for sample := 0; sample < iblen; sample++ {
				m := is[gr][0][sample] // mid
				s := is[gr][1][sample] // side
				is[gr][0][sample] = (m + s) / math.Sqrt2
				is[gr][1][sample] = (m - s) / math.Sqrt2
			}
		}
	}
}

func aliasReduction(gr, ch int, sideInfo sideInformation, is *[2][2][iblen]float32) {
	if sideInfo.WindowsSwitchingFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == blockShort && sideInfo.MixedBlockFlag[gr][ch] == 0 {
		return
	}

	nsb := 32
	if sideInfo.MixedBlockFlag[gr][ch] == 1 {
		nsb = 2
	}

	for sb := 1; sb < nsb; sb++ {
		for i := 0; i < 8; i++ {
			li := 18*sb - 1 - i
			ui := 18*sb + i
			is[gr][ch][li] = is[gr][ch][li]*cs[i] - is[gr][ch][ui]*ca[i]
			is[gr][ch][ui] = is[gr][ch][ui]*cs[i] + is[gr][ch][li]*ca[i]
		}
	}
}

func imdct(gr, ch int, blockType byte, is *[2][2][iblen]float32, prevSamples *[2][32][18]float32) {
	n := 36
	if blockType == blockShort {
		n = 12
	}
	halfN := n / 2

	nWin := 1
	if blockType == blockShort {
		nWin = 3
	}

	for block := 0; block < 32; block++ {
		samplesBlock := make([]float32, 36)
		if blockType == blockShort {
			for window := 0; window < nWin; window++ {
				for i := 0; i < n; i++ {
					xi := float32(0.0)
					for k := 0; k < halfN; k++ {
						s := is[gr][ch][block*18+nWin*k+window]
						xi += s * float32(math.Cos(math.Pi/float64(2*n)*float64(2*i+1+halfN)*float64(2*k+1)))
					}
					samplesBlock[6*window+i+6] += xi * winShape[blockType][i]
				}
			}
		} else {
			// nWin = 1
			for i := 0; i < n; i++ {
				xi := float32(0.0)
				for k := 0; k < halfN; k++ {
					s := is[gr][ch][block*18+k]
					xi += s * float32(math.Cos(math.Pi/float64(2*n)*float64(2*i+1+halfN)*float64(2*k+1)))
				}
				samplesBlock[i] = xi * winShape[blockType][i]
			}
		}

		// Overlapping
		for i := 0; i < 18; i++ {
			is[gr][ch][block*18+i] = samplesBlock[i] + prevSamples[ch][block][i]
			prevSamples[ch][block][i] = samplesBlock[i+18]
		}
	}
}

func frequencyInversion(gr, ch int, is *[2][2][iblen]float32) {
	for sb := 1; sb < 32; sb += 2 {
		for i := 1; i < 18; i += 2 {
			is[gr][ch][sb*18+i] *= -1
		}
	}
}

func synthFilterbank(gr, ch int, is *[2][2][iblen]float32, vVec *[2][1024]float32, pcm []float32) {
	uVec := [512]float32{}
	wVec := [512]float32{}
	//pcm := [iblen]float32{}

	tempVec := [32]float32{}

	for sb := 0; sb < 18; sb++ { // Loop through 18 samples in 32 subband blocks
		copy(vVec[ch][64:1024], vVec[ch][0:1024-64])

		// Offset --------------------------------------------------
		for i := 0; i < 32; i++ {
			tempVec[i] = is[gr][ch][i*18+sb]
		}

		// Matrix --------------------------------------------------
		for i := 0; i < 64; i++ {
			vVec[ch][i] = 0.0
			for j := 0; j < 32; j++ {
				vVec[ch][i] += tempVec[j] * synthN[i][j]
			}
		}

		// Build the U vector --------------------------------------------------
		for i := 0; i < 8; i++ {
			for j := 0; j < 32; j++ {
				uVec[i*64+j] = vVec[ch][i*128+j]
				uVec[i*64+32+j] = vVec[ch][i*128+96+j]
			}
		}

		// Build the V vector --------------------------------------------------
		for i := 0; i < 512; i++ {
			wVec[i] = uVec[i] * synthD[i]
		}

		// Calc 32 samples --------------------------------------------------
		for i := 0; i < 32; i++ {
			S := float32(0)
			for j := 0; j < 512; j += 32 {
				S += wVec[j+i]
			}
			//sample := int(S * math.MaxInt16)
			//if sample > math.MaxInt16 {
			//	sample = math.MaxInt16
			//} else if sample < -math.MaxInt16 {
			//	sample = -math.MaxInt16
			//}

			//s := int16(sample)
			idx := 2 * (32*sb + i)
			if ch == 0 {
				pcm[idx] = S
			} else {
				pcm[idx+1] = S
			}
		}
	}
}

var Mp1 = mpeg
var Mp2 = mpeg
var Mp3 = mpeg
