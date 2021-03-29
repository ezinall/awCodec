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
	HandlerType uint32    // unsigned int(32)
	Reserved    [3]uint32 // unsigned int(32)[3]
	//Name        string
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
	//minfType = [4]byte{'m', 'i', 'n', 'f'}
)

//const (
//	ftypType uint32 = 0x66747970
//	moovType uint32 = 0x6D6F6F76
//)

func mpeg12(file *bytes.Buffer) {
	var moov *bytes.Buffer
	var trackList []*bytes.Buffer
	var mdat []*bytes.Buffer

	fileTypeBox := FileTypeBox{}
	if err := binary.Read(file, binary.BigEndian, &fileTypeBox); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	compatibleBrands := make([][4]byte, (fileTypeBox.Size-uint32(unsafe.Sizeof(FileTypeBox{})))/4)
	if err := binary.Read(file, binary.BigEndian, &compatibleBrands); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Printf("%+v CompatibleBrands:%s\n", fileTypeBox, compatibleBrands)

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

		boxData := bytes.NewBuffer(file.Next(int(size - uint64(unsafe.Sizeof(Box{})))))

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
	fmt.Printf("%+v\n", movieHeaderBox)
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
			trackList = append(trackList, boxData)
		}
	}

	fmt.Println()
	for _, trak := range trackList {
		trackHeaderBox := TrackHeaderBox{}
		binary.Read(trak, binary.BigEndian, &trackHeaderBox)
		fmt.Printf("%+v\n", trackHeaderBox)

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
		fmt.Printf("\t%+v\n", mediaHeaderBox)
		//fmt.Printf("%016b\n", mediaHeaderBox.Language)
		//fmt.Printf("%x, %x, %x\n", 0b10101+0x60, 0b01110+0x60, 0b00100+0x60) // und 0 10101 01110 00100

		handlerBox := HandlerBox{}
		binary.Read(mdia, binary.BigEndian, &handlerBox)
		name, _ := (mdia).ReadString('\x00')
		fmt.Printf("\t%+v Name:%s\n", handlerBox, name)

		box := Box{}
		binary.Read(mdia, binary.BigEndian, &box)
		fmt.Printf("\t{Size:%d Type:%s}\n", box.Size, box.Type)
		minf := bytes.NewBuffer((mdia).Next(int(box.Size - uint32(unsafe.Sizeof(Box{})))))
		for minf.Len() != 0 {
			binary.Read(minf, binary.BigEndian, &box)
			fmt.Printf("\t\t{Size:%d Type:%s}\n", box.Size, box.Type)
			minf.Next(int(box.Size - uint32(unsafe.Sizeof(Box{}))))
		}
		fmt.Println()
	}
}

var Mp4 = mpeg12
