package riff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"unsafe"
)

const (
	// avih dwFlags
	AVIF_HASINDEX       = 0x00000010
	AVIF_MUSTUSEINDEX   = 0x00000020
	AVIF_ISINTERLEAVED  = 0x00000100
	AVIF_WASCAPTUREFILE = 0x00010000
	AVIF_COPYRIGHTED    = 0x00020000
	AVIF_TRUSTCKTYPE    = 0x00000800

	// strh dwFlags
	AVISF_DISABLED         = 0x00000001
	AVISF_VIDEO_PALCHANGES = 0x00010000

	// bIndexType
	AVI_INDEX_OF_INDEXES = 0x00
	AVI_INDEX_OF_CHUNKS  = 0x01
	AVI_INDEX_IS_DATA    = 0x80

	// bIndexSubtype
	AVI_INDEX_2FIELD = 0x01

	// idx1 chunk dwFlags
	AVIIF_LIST      = 0x00000001
	AVIIF_KEYFRAME  = 0x00000010
	AVIIF_FIRSTPART = 0x00000020
	AVIIF_LASTPART  = 0x00000040
	AVIIF_NOTIME    = 0x00000100
)

type Chunk struct {
	DwFourCC [4]byte
	DwSize   uint32
	Data     []byte
}

type List struct {
	DwList   [4]byte
	DwSize   uint32
	DwFourCC [4]byte
	Data     []byte
}

type MainAVIHeader struct {
	DwMicroSecPerFrame    uint32 // frame display rate (or 0)
	DwMaxBytesPerSec      uint32 // max. transfer rate
	DwPaddingGranularity  uint32 // pad to multiples of this size
	DwFlags               uint32 // the ever-present flags
	DwTotalFrames         uint32 // # frames in file
	DwInitialFrames       uint32
	DwStreams             uint32
	DwSuggestedBufferSize uint32
	DwWidth               uint32
	DwHeight              uint32
	//DwReserved            [4]uint32
	DwScale  uint32
	DwRate   uint32
	DwStart  uint32
	DwLength uint32
}

type AVIStreamHeader struct {
	FccType               [4]byte // string
	FccHandler            [4]byte // string
	DwFlags               uint32
	WPriority             uint16
	WLanguage             uint16
	DwInitialFrames       uint32
	DwScale               uint32
	DwRate                uint32 // dwRate / dwScale == samples/second
	DwStart               uint32
	DwLength              uint32 // In units above...
	DwSuggestedBufferSize uint32
	DwQuality             uint32
	DwSampleSize          uint32
	RcFrame               struct {
		Left   uint16
		Top    uint16
		Right  uint16
		Bottom uint16
	}
}

type BitMapInfoHeader struct {
	BiSize          uint32
	BiWidth         uint32
	BiHeight        uint32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   [4]byte // string
	BiSizeImage     uint32
	BiXPelsPerMeter uint32
	BiYPelsPerMeter uint32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type WaveFormat struct {
	WFormatTag      uint16
	NChannels       uint16
	NSamplesPerSec  uint32
	NAvgBytesPerSec uint32
	NBlockAlign     uint16
}

type AVIIndexEntry struct {
	Ckid          [4]byte
	DwFlags       uint32
	DwChunkOffset uint32
	DwChunkLength uint32
}

func Avi(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	AVI := List{}
	copy(AVI.DwList[:], content[:4])
	AVI.DwSize = binary.LittleEndian.Uint32(content[4:8])
	copy(AVI.DwFourCC[:], content[8:12])
	AVI.Data = content[12:]

	if string(AVI.DwFourCC[:]) != "AVI " {
		return
	}

	hdrl := List{}
	avih := MainAVIHeader{}
	var strl []List
	var vids []struct {
		strh AVIStreamHeader
		strf BitMapInfoHeader
	}
	var auds []struct {
		strh AVIStreamHeader
		strf WaveFormat
	}
	//var txts []struct {
	//	AVIStreamHeader
	//}
	movi := List{}
	// info := new(map[string]string)
	idx1 := Chunk{}

	index := 0
	for index != len(AVI.Data) {
		if dwList := AVI.Data[index : index+4]; string(dwList) == "LIST" {
			dwSize := binary.LittleEndian.Uint32(AVI.Data[index+4 : index+8])
			dwFourCC := AVI.Data[index+8 : index+12]

			if string(dwFourCC) == "hdrl" {
				copy(hdrl.DwList[:], dwList)
				hdrl.DwSize = dwSize
				copy(hdrl.DwFourCC[:], dwFourCC)
				hdrl.Data = AVI.Data[index+12 : index+12+int(dwSize)-4]
			} else if string(dwFourCC) == "movi" {
				copy(movi.DwList[:], dwList)
				movi.DwSize = dwSize
				copy(movi.DwFourCC[:], dwFourCC)
				movi.Data = AVI.Data[index+12 : index+12+int(dwSize)-4]
			} else if string(dwFourCC) == "INFO" {
				// past to info
			}

			index += int(dwSize) + 8
		} else if dwFourCC := AVI.Data[index : index+4]; string(dwFourCC) == "idx1" {
			copy(idx1.DwFourCC[:], AVI.Data[index:index+4])
			idx1.DwSize = binary.LittleEndian.Uint32(AVI.Data[index+4 : index+8])
			idx1.Data = AVI.Data[index+8 : index+12+int(idx1.DwSize)-4]

			index += int(idx1.DwSize) + 8
		} else {
			dwSize := binary.LittleEndian.Uint32(AVI.Data[index+4 : index+8])

			index += int(dwSize) + 8
		}
	}

	index = 0
	for index != len(hdrl.Data) {
		if dwFourCC := hdrl.Data[index : index+4]; string(dwFourCC) == "avih" {
			dwSize := binary.LittleEndian.Uint32(hdrl.Data[index+4 : index+8])

			r := bytes.NewReader(hdrl.Data[index+8 : index+8+int(dwSize)])
			if err := binary.Read(r, binary.LittleEndian, &avih); err != nil {
				fmt.Println("binary.Read failed:", err)
			}

			index += int(dwSize) + 8
		} else if dwList := hdrl.Data[index : index+4]; string(dwList) == "LIST" {
			list := List{}
			copy(list.DwList[:], dwList)
			list.DwSize = binary.LittleEndian.Uint32(hdrl.Data[index+4 : index+8])
			copy(list.DwFourCC[:], hdrl.Data[index+8:index+12])
			list.Data = hdrl.Data[index+12 : index+12+int(list.DwSize)-4]

			strl = append(strl, list)

			index += int(list.DwSize) + 8
		} else {
			dwSize := binary.LittleEndian.Uint32(hdrl.Data[index+4 : index+8])

			index += int(dwSize) + 8
		}
	}

	for _, list := range strl {
		dwSize := binary.LittleEndian.Uint32(list.Data[4:8])

		strh := AVIStreamHeader{}
		r := bytes.NewReader(list.Data[8 : 8+int(dwSize)])
		if err := binary.Read(r, binary.LittleEndian, &strh); err != nil {
			fmt.Println("binary.Read failed:", err)
		}

		index := 8 + int(dwSize)
		switch string(strh.FccType[:]) {
		case "vids":
			dwSize := binary.LittleEndian.Uint32(list.Data[index+4 : index+8])

			strf := BitMapInfoHeader{}
			r := bytes.NewReader(list.Data[index+8 : index+12+int(dwSize)])
			if err := binary.Read(r, binary.LittleEndian, &strf); err != nil {
				fmt.Println("binary.Read failed:", err)
			}

			vids = append(vids, struct {
				strh AVIStreamHeader
				strf BitMapInfoHeader
			}{strh, strf})
		case "auds":
			dwSize := binary.LittleEndian.Uint32(list.Data[index+4 : index+8])

			strf := WaveFormat{} //TODO реализовать смену структуры в зависимости от размера
			r := bytes.NewReader(list.Data[index+8 : index+8+int(dwSize)])
			if err := binary.Read(r, binary.LittleEndian, &strf); err != nil {
				fmt.Println("binary.Read failed:", err)
			}

			auds = append(auds, struct {
				strh AVIStreamHeader
				strf WaveFormat
			}{strh, strf})
		case "txts":
			fmt.Printf("%+v\n", strh)
		case "mids":
		}
	}

	index = 0
	for index != len(movi.Data) {
		if dwFourCC := movi.Data[index : index+4]; strings.Contains(string(dwFourCC), "wb") {
			dwSize := binary.LittleEndian.Uint32(movi.Data[index+4 : index+8])
			_ = movi.Data[index+8 : index+8+int(dwSize)]

			// fmt.Printf("\t\t%s %d\n", dwFourCC, dwSize)

			index += int(dwSize) + 8 + int(dwSize%2)
		} else if dwFourCC := movi.Data[index : index+4]; strings.Contains(string(dwFourCC), "db") {
			dwSize := binary.LittleEndian.Uint32(movi.Data[index+4 : index+8])
			_ = movi.Data[index+8 : index+8+int(dwSize)]

			// fmt.Printf("\t\t%s %d\n", dwFourCC, dwSize)

			index += int(dwSize) + 8
		} else if dwFourCC := movi.Data[index : index+4]; strings.Contains(string(dwFourCC), "dc") {
			dwSize := binary.LittleEndian.Uint32(movi.Data[index+4 : index+8])
			_ = movi.Data[index+8 : index+8+int(dwSize)]

			// fmt.Printf("\t\t%s %d\n", dwFourCC, dwSize)

			// moveData = append(moveData, movi.Data[index+8:index+8+int(dwSize)])

			index += int(dwSize) + 8 + int(dwSize%2)
		} else if dwFourCC := movi.Data[index : index+4]; strings.Contains(string(dwFourCC), "tx") {
			dwSize := binary.LittleEndian.Uint32(movi.Data[index+4 : index+8])
			_ = movi.Data[index+8 : index+8+int(dwSize)]

			// fmt.Printf("\t\t%s %d\n", dwFourCC, dwSize)

			// fmt.Printf("%s\n", movi.Data[index+8:index+8+int(dwSize)])

			index += int(dwSize) + 8
		} else if dwFourCC := movi.Data[index : index+4]; strings.Contains(string(dwFourCC), "ix") {
			dwSize := binary.LittleEndian.Uint32(movi.Data[index+4 : index+8])
			_ = movi.Data[index+8 : index+8+int(dwSize)]

			// fmt.Printf("\t\t%s %d\n", dwFourCC, dwSize)

			index += int(dwSize) + 8
		}
	}

	index = 0
	for index != len(idx1.Data) {
		chunk := AVIIndexEntry{}
		r := bytes.NewReader(idx1.Data[index : index+16])
		if err := binary.Read(r, binary.LittleEndian, &chunk); err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		// fmt.Printf("\t\tCkid: %s DwFlags: %d DwChunkOffset: %d DwChunkLength:%d\n", chunk.Ckid, chunk.DwFlags, chunk.DwChunkOffset, chunk.DwChunkLength)
		// fmt.Printf("\t\t%+v\n", chunk)

		index += 16
	}

	fmt.Printf("%s %d %s\n", AVI.DwList, AVI.DwSize, AVI.DwFourCC)
	fmt.Printf("\t%s %d %s\n", hdrl.DwList, hdrl.DwSize, hdrl.DwFourCC)
	fmt.Printf("\t\tavih %d %+v\n", unsafe.Sizeof(avih), avih)
	// for _, list := range strl {
	// fmt.Printf("\t\t%s %d %s\n", list.DwList, list.DwSize, list.DwFourCC)
	// }
	for _, chunk := range vids {
		fmt.Printf("\t\tstrh %d %+v\n", unsafe.Sizeof(chunk.strh), chunk.strh)
		fmt.Printf("\t\tstrf %d %+v\n", unsafe.Sizeof(chunk.strf), chunk.strf)
	}
	for _, chunk := range auds {
		fmt.Printf("\t\tstrh %d %+v\n", unsafe.Sizeof(chunk.strh), chunk.strh)
		fmt.Printf("\t\tstrf %d %+v\n", unsafe.Sizeof(chunk.strf), chunk.strf)
	}
	fmt.Printf("\t%s %d %s\n", movi.DwList, movi.DwSize, movi.DwFourCC)
	fmt.Printf("\t%s %d\n", idx1.DwFourCC, idx1.DwSize)
}
