package mpeg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
	"unsafe"
)

const (
	trackIsAudio = 0x0100
)

var UTS = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC)

type Box struct {
	Size uint32
	Type [4]byte
}

type FullBox struct {
	Box
	Version uint8
	Flags   [3]byte
}

type FileTypeBox struct {
	Box
	MajorBrand   [4]byte
	MinorVersion [4]byte
	//CompatibleBrands [][4]byte
}

type MovieHeaderBox struct {
	FullBox
	CreationTime     uint32    // unsigned int(32) / unsigned int(64)
	ModificationTime uint32    // unsigned int(32) / unsigned int(64)
	Timescale        uint32    // unsigned int(32)
	Duration         uint32    // unsigned int(32) / unsigned int(64)
	Rate             int32     // template int(32)
	Volume           int16     // template int(16)
	Reserved         uint16    // const bit(16)
	Reserved_        [2]uint32 // const unsigned int(32)[2]
	Matrix           [9]int32  // template int(32)[9]
	PreDefined       [6]int32  // bit(32)[6]
	NextTrackID      uint32    // unsigned int(32)
}

type TrackHeaderBox struct {
	FullBox
	CreationTime     uint32    // unsigned int(32) / unsigned int(64)
	ModificationTime uint32    // unsigned int(32) / unsigned int(64)
	TrackID          uint32    // unsigned int(32)
	Reserved         uint32    // unsigned int(32)
	Duration         uint32    // unsigned int(32) / unsigned int(64)
	Reserved_        [2]uint32 // unsigned int(32)[2]
	Layer            int16     // template int(16)
	AlternateGroup   int16     // template int(16)
	Volume           int16     // template int(16)
	Reserved__       uint16    // unsigned int(16)
	Matrix           [9]uint32 // template int(32)[9]
	Width            uint32    // unsigned int(32)
	Height           uint32    // unsigned int(32)
}

type MediaHeaderBox struct {
	FullBox
	CreationTime     uint32 // unsigned int(32) / unsigned int(64)
	ModificationTime uint32 // unsigned int(32) / unsigned int(64)
	Timescale        uint32 // unsigned int(32)
	Duration         uint32 // unsigned int(32) / unsigned int(64)
	Language         uint16 // 1 bit pad and unsigned int(5)[3]
	PreDefined       uint16 // unsigned int(16)
}

type HandlerBox struct {
	FullBox
	PreDefined  uint32    // unsigned int(32)
	HandlerType [4]byte   // unsigned int(32)
	Reserved    [3]uint32 // unsigned int(32)[3]
	//Name        string
}

type VideoMediaHeaderBox struct {
	GraphicsMode uint16    // unsigned int(16)
	OpColor      [3]uint16 // unsigned int(16)[3]
}

type SoundMediaHeaderBox struct {
	Balance  int16  // int(16)
	Reserved uint16 // unsigned int(16)
}

type HintMediaHeaderBox struct {
	MaxPDUSize uint16 // unsigned int(16)
	AvgPDUSize uint16 // unsigned int(16)
	MaxBitRate uint32 // unsigned int(32)
	AvgBitRate uint32 // unsigned int(32)
	Reserved   uint32 // unsigned int(32)
}

type NullMediaHeaderBox struct {
}

type DataReferenceBox struct {
	FullBox
	EntryCount uint32 // unsigned int(32)
}

type SampleTableBox struct {
	Box
}

type SampleDescriptionBox struct {
	FullBox
	//I          int32  // int
	EntryCount uint32 // unsigned int(32)
}

type SampleEntry struct {
	Box
	Reserved           [6]uint8 // unsigned int(8)[6]
	DataReferenceIndex uint16   // unsigned int(16)
}

type HintSampleEntry struct {
	SampleEntry
}

type VisualSampleEntry struct {
	SampleEntry
	PreDefined      uint16    // unsigned int(16)
	Reserved        uint16    // unsigned int(16)
	PreDefined_     [3]uint32 // unsigned int(32)[3]
	Width           uint16    // unsigned int(16)
	Height          uint16    // unsigned int(16)
	HorizResolution uint32    // unsigned int(32)
	VertResolution  uint32    // unsigned int(32)
	Reserved_       uint32    // unsigned int(32)
	FrameCount      uint16    // unsigned int(16)
	CompressorName  [32]byte  // string[32]
	Depth           uint16    // unsigned int(16)
	PreDefined__    int16     // int(16)
}

type AudioSampleEntry struct {
	SampleEntry
	Reserved     [2]uint32
	ChannelCount uint16
	SampleSize   uint16
	PreDefined   uint16
	Reserved_    uint16
	SampleRate   uint32
}

// ISO/IEC 14496-15

type AVCDecoderConfigurationRecord struct {
	ConfigurationVersion       uint8
	AVCProfileIndication       uint8
	ProfileCompatibility       uint8
	AVCLevelIndication         uint8
	LengthSizeMinusOne         uint8 // bit(6) reserved and unsigned int(2) lengthSizeMinusOne
	NumOfSequenceParameterSets uint8 // bit(3) reserved and unsigned int(5) numOfSequenceParameterSets
}

type AVCConfigurationBox struct {
	Box
	AVCDecoderConfigurationRecord // AVCConfig
}

// AVCConfigurationBoxExt non standard struct
type AVCConfigurationBoxExt struct {
	ChromaFormat                 uint8 // bit(6) reserved and unsigned int(2) chroma_format
	BitDepthLumaMinus8           uint8 // bit(5) reserved and unsigned int(3) bit_depth_luma_minus8
	BitDepthChromaMinus8         uint8 // bit(5) reserved and unsigned int(3) bit_depth_chroma_minus8
	NumOfSequenceParameterSetExt uint8 // unsigned int(8)
}

type SampleSizeBox struct {
	SampleSize  uint32
	SampleCount uint32
}

type CompactSampleSizeBox struct {
	Reserved    uint16 // unsigned int(24)
	FieldSize   uint16 // unisgned int(8)
	SampleCount uint32 // unsigned int(32)
}

var (
	//ftypType = [4]byte{'f', 't', 'y', 'p'}
	moovType = [4]byte{'m', 'o', 'o', 'v'}
	mdatType = [4]byte{'m', 'd', 'a', 't'}
	//mvhdType = [4]byte{'m', 'v', 'h', 'd'}
	trakType = [4]byte{'t', 'r', 'a', 'k'}
	//tkhdType = [4]byte{'t', 'k', 'h', 'd'}
	mdiaType = [4]byte{'m', 'd', 'i', 'a'}
	//mdhdType = [4]byte{'m', 'd', 'h', 'd'}
	//hdlrType = [4]byte{'h', 'd', 'l', 'r'}
	handlerVide = [4]byte{'v', 'i', 'd', 'e'}
	handlerSoun = [4]byte{'s', 'o', 'u', 'n'}
	handlerHint = [4]byte{'h', 'i', 'n', 't'}
	//minfType = [4]byte{'m', 'i', 'n', 'f'}

	vmhdType = [4]byte{'v', 'm', 'h', 'd'}
	smhdType = [4]byte{'s', 'm', 'h', 'd'}
	hmhdType = [4]byte{'h', 'm', 'h', 'd'}
	nmhdType = [4]byte{'n', 'm', 'h', 'd'}

	//urlType = [4]byte{'u', 'r', 'l', ' '}
	urnType = [4]byte{'u', 'r', 'n', ' '}

	sttsType = [4]byte{'s', 't', 't', 's'}
	stscType = [4]byte{'s', 't', 's', 'c'}
	stcoType = [4]byte{'s', 't', 'c', 'o'}
	co64Type = [4]byte{'c', '0', '6', '4'}
	stszType = [4]byte{'s', 't', 's', 'z'}
	stz2Type = [4]byte{'s', 't', 'z', '2'}
)

//const (
//	ftypType uint32 = 0x66747970
//	moovType uint32 = 0x6D6F6F76
//)

func mpeg12(file *bytes.Reader) {
	var moov *bytes.Buffer
	var trakList []*bytes.Buffer
	var mdat []*bytes.Buffer

	fileTypeBox := FileTypeBox{}
	if err := binary.Read(file, binary.BigEndian, &fileTypeBox); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	compatibleBrands := make([][4]byte, (fileTypeBox.Size-uint32(unsafe.Sizeof(FileTypeBox{})))/4)
	if err := binary.Read(file, binary.BigEndian, &compatibleBrands); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	fmt.Printf("{Box:{Size:%d Type:%s} MajorBrand:%s MinorVersion:%s} CompatibleBrands:%s\n",
		fileTypeBox.Size, fileTypeBox.Type, fileTypeBox.MajorBrand, fileTypeBox.MinorVersion, compatibleBrands)

	for file.Len() != 0 {
		box := Box{}
		if err := binary.Read(file, binary.BigEndian, &box); err != nil {
			fmt.Println("binary.Read failed:", err)
		}

		fmt.Printf("{Size:%d Type:%s}\n", box.Size, box.Type)

		size := uint64(box.Size)
		if size == 1 {
			if err := binary.Read(file, binary.BigEndian, &size); err != nil {
				fmt.Println("binary.Read failed:", err)
			}
			size -= 8 // self size
		} else if size == 0 {
			size = uint64(file.Len())
		}

		buf := make([]byte, int(size-uint64(unsafe.Sizeof(Box{}))))
		file.Read(buf)
		boxData := bytes.NewBuffer(buf)

		//boxData := bytes.NewBuffer(file.Next(int(size - uint64(unsafe.Sizeof(Box{})))))

		if box.Type == moovType {
			moov = boxData

		} else if box.Type == mdatType {
			mdat = append(mdat, boxData)
		}
	}

	fmt.Println()

	movieHeaderBox := MovieHeaderBox{}
	if err := binary.Read(moov, binary.BigEndian, &movieHeaderBox); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	fmt.Printf("MovieHeaderBox %+v\n", movieHeaderBox)

	for (moov).Len() != 0 {
		box := Box{}
		if err := binary.Read(moov, binary.BigEndian, &box); err != nil {
			fmt.Println("binary.Read failed:", err)
		}

		fmt.Printf("{Size:%d Type:%s}\n", box.Size, box.Type)

		size := uint64(box.Size)
		if size == 1 {
			if err := binary.Read(moov, binary.BigEndian, &box); err != nil {
				fmt.Println("binary.Read failed:", err)
			}
			size -= 8
		} else if size == 0 {
			size = uint64((moov).Len())
		}

		boxData := bytes.NewBuffer((moov).Next(int(size - uint64(unsafe.Sizeof(Box{})))))

		if box.Type == trakType {
			trakList = append(trakList, boxData)
		}
	}

	fmt.Println()
	for _, trak := range trakList {
		trackHeaderBox := TrackHeaderBox{}
		binary.Read(trak, binary.BigEndian, &trackHeaderBox)

		fmt.Printf("TrackHeaderBox %+v\n", trackHeaderBox)

		var mdia *bytes.Buffer

		for trak.Len() != 0 {
			box := Box{}
			if err := binary.Read(trak, binary.BigEndian, &box); err != nil {
				fmt.Println("binary.Read failed:", err)
			}

			fmt.Printf("{Size:%d Type:%s}\n", box.Size, box.Type)

			if box.Type == mdiaType {
				mdia = bytes.NewBuffer(trak.Next(int(box.Size - uint32(unsafe.Sizeof(Box{})))))
			} else {
				trak.Next(int(box.Size - uint32(unsafe.Sizeof(Box{}))))
			}
		}

		mediaHeaderBox := MediaHeaderBox{}
		binary.Read(mdia, binary.BigEndian, &mediaHeaderBox)

		fmt.Printf("\tMediaHeaderBox {FullBox:{Box:{Size:%d Type:%s} Version:%d Flags:%d} CreationTime:%d ModificationTime:%d Timescale:%d Duration:%d Language:%d PreDefined:%d}\n",
			mediaHeaderBox.Size, mediaHeaderBox.Type, mediaHeaderBox.Version, mediaHeaderBox.Flags, mediaHeaderBox.CreationTime, mediaHeaderBox.ModificationTime, mediaHeaderBox.Timescale, mediaHeaderBox.Duration, mediaHeaderBox.Language, mediaHeaderBox.PreDefined)
		//fmt.Printf("%016b\n", mediaHeaderBox.Language)
		//fmt.Printf("%x, %x, %x\n", 0b10101+0x60, 0b01110+0x60, 0b00100+0x60) // und 0 10101 01110 00100

		handlerBox := HandlerBox{}
		binary.Read(mdia, binary.BigEndian, &handlerBox)
		name, _ := (mdia).ReadString('\x00')

		fmt.Printf("\tHandlerBox {FullBox:{Box:{Size:%d Type:%s} Version:%d Flags:%d} PreDefined:%d HandlerType:%s Reserved:%d} Name:%s\n",
			handlerBox.Size, handlerBox.Type, handlerBox.Version, handlerBox.Flags, handlerBox.PreDefined, handlerBox.HandlerType, handlerBox.Reserved, name)

		box := Box{}
		binary.Read(mdia, binary.BigEndian, &box)

		fmt.Printf("\t{Size:%d Type:%s}\n", box.Size, box.Type)

		minf := bytes.NewBuffer((mdia).Next(int(box.Size - uint32(unsafe.Sizeof(Box{})))))

		fullBox := FullBox{}
		binary.Read(minf, binary.BigEndian, &fullBox)

		fmt.Printf("\t\ttypeMediaHeaderBox {Box:{Size:%d Type:%s} Version:%d Flags:%d} ", fullBox.Size, fullBox.Type, fullBox.Version, fullBox.Flags)

		var typeMediaHeaderBox interface{}
		if fullBox.Type == vmhdType {
			header := VideoMediaHeaderBox{}
			binary.Read(minf, binary.BigEndian, &header)
			typeMediaHeaderBox = header
		} else if fullBox.Type == smhdType {
			header := SoundMediaHeaderBox{}
			binary.Read(minf, binary.BigEndian, &header)
			typeMediaHeaderBox = header
		} else if fullBox.Type == hmhdType {
			header := HintMediaHeaderBox{}
			binary.Read(minf, binary.BigEndian, &header)
			typeMediaHeaderBox = header
		} else if fullBox.Type == nmhdType {
			typeMediaHeaderBox = NullMediaHeaderBox{}
		}

		fmt.Printf("%+v\n", typeMediaHeaderBox)

		binary.Read(minf, binary.BigEndian, &box)

		fmt.Printf("\t\t{Size:%d Type:%s}\n", box.Size, box.Type)

		dataReferenceBox := DataReferenceBox{}
		binary.Read(minf, binary.BigEndian, &dataReferenceBox)

		fmt.Printf("\t\t\tDataReferenceBox {FullBox:{Box:{Size:%d Type:%s} Version:%d Flags:%d} EntryCount:%d}\n",
			dataReferenceBox.Size, dataReferenceBox.Type, dataReferenceBox.Version, dataReferenceBox.Flags, dataReferenceBox.EntryCount)

		for i := 1; i <= int(dataReferenceBox.EntryCount); i++ {
			binary.Read(minf, binary.BigEndian, &fullBox)

			fmt.Printf("\t\t\t{Box:{Size:%d Type:%s} Version:%d Flags:%d} ", fullBox.Size, fullBox.Type, fullBox.Version, fullBox.Flags)

			if fullBox.Size > uint32(unsafe.Sizeof(FullBox{})) {
				var name string
				if fullBox.Type == urnType {
					name, _ = minf.ReadString('\x00')
				}
				location, _ := minf.ReadString('\x00')

				fmt.Printf("Name:%s Location:%s\n", name, location)
			} else {
				fmt.Println()
			}
		}

		sampleTableBox := SampleTableBox{}
		binary.Read(minf, binary.BigEndian, &sampleTableBox)

		fmt.Printf("\t\t{Size:%d Type:%s}\n", sampleTableBox.Size, sampleTableBox.Type)

		sampleDescriptionBox := SampleDescriptionBox{}
		binary.Read(minf, binary.BigEndian, &sampleDescriptionBox)

		fmt.Printf("\t\t\tSampleDescriptionBox {FullBox:{Box:{Size:%d Type:%s} Version:%d Flags:%d} EntryCount:%d}\n",
			sampleDescriptionBox.Size, sampleDescriptionBox.Type, sampleDescriptionBox.Version, sampleDescriptionBox.Flags, sampleDescriptionBox.EntryCount)

		for i := 1; i <= int(sampleDescriptionBox.EntryCount); i++ {
			switch handlerBox.HandlerType {
			case handlerSoun:
				audioSampleEntry := AudioSampleEntry{}
				binary.Read(minf, binary.BigEndian, &audioSampleEntry)
				fmt.Printf("\t\t\tAudioSampleEntry {SampleEntry:{Box:{Size:%d Type:%s} Reserved:%d DataReferenceIndex:%d} Reserved:%d ChannelCount:%d SampleSize:%d PreDefined:%d Reserved_:%d SampleRate:%d}\n",
					audioSampleEntry.Size, audioSampleEntry.Type, audioSampleEntry.Reserved, audioSampleEntry.DataReferenceIndex, audioSampleEntry.Reserved, audioSampleEntry.ChannelCount, audioSampleEntry.SampleSize, audioSampleEntry.PreDefined, audioSampleEntry.Reserved_, audioSampleEntry.SampleRate)

			case handlerVide:
				visualSampleEntry := VisualSampleEntry{}
				binary.Read(minf, binary.BigEndian, &visualSampleEntry)

				fmt.Printf("\t\t\tVisualSampleEntry {SampleEntry:{Box:{Size:%d Type:%s} Reserved:%d DataReferenceIndex:%d} PreDefined:%d Reserved:%d PreDefined_:%d Width:%d Height:%d HorizResolution:%d VertResolution:%d Reserved_:%d FrameCount:%d CompressorName:%s Depth:%d PreDefined__:%d}\n",
					visualSampleEntry.Size, visualSampleEntry.Type, visualSampleEntry.Reserved, visualSampleEntry.DataReferenceIndex, visualSampleEntry.PreDefined, visualSampleEntry.Reserved, visualSampleEntry.PreDefined_, visualSampleEntry.Width, visualSampleEntry.Height, visualSampleEntry.HorizResolution, visualSampleEntry.VertResolution, visualSampleEntry.Reserved_, visualSampleEntry.FrameCount, visualSampleEntry.CompressorName, visualSampleEntry.Depth, visualSampleEntry.PreDefined__)

				if visualSampleEntry.Type == [4]byte{'a', 'v', 'c', '1'} {
					avcConfigurationBox := AVCConfigurationBox{}
					binary.Read(minf, binary.BigEndian, &avcConfigurationBox)

					fmt.Printf("\t\t\tAVCConfigurationBox {Box:{Size:%d Type:%s} AVCConfig:{ConfigurationVersion:%d AVCProfileIndication:%d ProfileCompatibility:%d AVCLevelIndication:%d LengthSizeMinusOne:%d NumOfSequenceParameterSets:%d}}\n",
						avcConfigurationBox.Size, avcConfigurationBox.Type, avcConfigurationBox.ConfigurationVersion, avcConfigurationBox.AVCProfileIndication, avcConfigurationBox.AVCProfileIndication, avcConfigurationBox.ProfileCompatibility, avcConfigurationBox.LengthSizeMinusOne&0b00000011, avcConfigurationBox.NumOfSequenceParameterSets&0b00011111)

					for i = 0; i < int(avcConfigurationBox.NumOfSequenceParameterSets&0b00011111); i++ {
						var sequenceParameterSetLength uint16
						binary.Read(minf, binary.BigEndian, &sequenceParameterSetLength)
						sequenceParameterSetNALUnit := minf.Next(int(sequenceParameterSetLength)) // sequenceParameterSetNALUnit
						fmt.Printf("\t\t\tsequenceParameterSetNALUnit %v\n", sequenceParameterSetNALUnit)
					}

					var numOfPictureParameterSets uint8
					binary.Read(minf, binary.BigEndian, &numOfPictureParameterSets)
					for i = 0; i < int(numOfPictureParameterSets); i++ {
						var pictureParameterSetLength uint16
						binary.Read(minf, binary.BigEndian, &pictureParameterSetLength)
						pictureParameterSetNALUnit := minf.Next(int(pictureParameterSetLength)) // pictureParameterSetNALUnit
						fmt.Printf("\t\t\tpictureParameterSetNALUnit %v\n", pictureParameterSetNALUnit)
					}
					//if avcConfigurationBox.AVCProfileIndication == 100 || avcConfigurationBox.AVCProfileIndication == 110 ||
					//	avcConfigurationBox.AVCProfileIndication == 122 || avcConfigurationBox.AVCProfileIndication == 144 {
					//	avcConfigurationBoxExt := AVCConfigurationBoxExt{}
					//	binary.Read(minf, binary.BigEndian, &avcConfigurationBoxExt)
					//
					//	fmt.Printf("\t\t\t{ChromaFormat:%d BitDepthLumaMinus8:%d BitDepthChromaMinus8:%d NumOfSequenceParameterSetExt:%d}\n",
					//		avcConfigurationBoxExt.ChromaFormat&0b00000011, avcConfigurationBoxExt.BitDepthLumaMinus8&0b00000111, avcConfigurationBoxExt.BitDepthChromaMinus8&0b00000111, avcConfigurationBoxExt.NumOfSequenceParameterSetExt)
					//
					//	for i = 0; i < 1; i++ {
					//		var sequenceParameterSetExtLength uint16
					//		binary.Read(minf, binary.BigEndian, &sequenceParameterSetExtLength)
					//		fmt.Println("sequenceParameterSetExtLength", sequenceParameterSetExtLength)
					//		if sequenceParameterSetExtLength != 0 {
					//			sequenceParameterSetExtNALUnit := minf.Next(int(sequenceParameterSetExtLength))
					//			fmt.Printf("%s\n", sequenceParameterSetExtNALUnit)
					//		}
					//	}
					//}
				}
			case handlerHint:

			}
		}

		for minf.Len() != 0 {
			binary.Read(minf, binary.BigEndian, &fullBox)

			fmt.Printf("\t\t\t{Size:%d Type:%s} ", fullBox.Size, fullBox.Type)

			if fullBox.Type == sttsType {
				var entryCount uint32
				binary.Read(minf, binary.BigEndian, &entryCount)

				fmt.Printf("EntryCount:%d\n", entryCount)

				for i := 1; i <= int(entryCount); i++ {
					var sampleCount, sampleDelta uint32
					binary.Read(minf, binary.BigEndian, &sampleCount)
					binary.Read(minf, binary.BigEndian, &sampleDelta)
					fmt.Printf("\t\t\tsampleCount %d sampleDelta %d\n", sampleCount, sampleDelta)
				}
			} else if fullBox.Type == stscType {
				var entryCount uint32
				binary.Read(minf, binary.BigEndian, &entryCount)

				fmt.Printf("EntryCount:%d\n", entryCount)

				for i := 1; i <= int(entryCount); i++ {
					var firstChunk, samplePerChunk, sampleDescriptionIndex uint32
					binary.Read(minf, binary.BigEndian, &firstChunk)
					binary.Read(minf, binary.BigEndian, &samplePerChunk)
					binary.Read(minf, binary.BigEndian, &sampleDescriptionIndex)
					//fmt.Printf("\t\t\tfirstChunk %d samplePerChunk %d sampleDescriptionIndex %d\n",
					//	firstChunk, samplePerChunk, sampleDescriptionIndex)
				}

				// Exactly one variant 'stsz' or 'stz2' must be present
			} else if fullBox.Type == stszType {
				sampleSizeBox := SampleSizeBox{}
				binary.Read(minf, binary.BigEndian, &sampleSizeBox)

				fmt.Printf("SampleSizeBox:%+v\n", sampleSizeBox)

				if sampleSizeBox.SampleSize == 0 {
					for i := 1; i <= int(sampleSizeBox.SampleCount); i++ {
						var entrySize uint32
						binary.Read(minf, binary.BigEndian, &entrySize)
					}
				}

			} else if fullBox.Type == stz2Type {
				compactSampleSizeBox := CompactSampleSizeBox{}
				binary.Read(minf, binary.BigEndian, &compactSampleSizeBox)

				fmt.Printf("CompactSampleSizeBox:%+v\n", compactSampleSizeBox)

				for i := 1; i <= int(compactSampleSizeBox.SampleCount); i++ {
					var entrySize = 1 << (compactSampleSizeBox.FieldSize & 0xFF)
					binary.Read(minf, binary.BigEndian, &entrySize)
				}

				// Exactly one variant 'stco' or 'co64' must be present
			} else if fullBox.Type == stcoType || fullBox.Type == co64Type {
				var entryCount uint32
				binary.Read(minf, binary.BigEndian, &entryCount)
				fmt.Printf("EntryCount:%d\n", entryCount)

				for i := 1; i <= int(entryCount); i++ {
					if fullBox.Type == stcoType {
						var chunkOffset uint32
						binary.Read(minf, binary.BigEndian, &chunkOffset)
						//fmt.Printf("\t\t\tchunkOffset %d\n", chunkOffset)
					}
				}
			} else {
				fmt.Println()
				minf.Next(int(fullBox.Size - uint32(unsafe.Sizeof(FullBox{}))))
			}
		}
		fmt.Println()
	}
}

var Mp4 = mpeg12
