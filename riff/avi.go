package riff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"unsafe"
)

const (
	// avih dwFlags
	AvifHasindex       = 0x00000010
	AvifMustuseindex   = 0x00000020
	AvifIsinterleaved  = 0x00000100
	AvifWascapturefile = 0x00010000
	AvifCopyrighted    = 0x00020000
	AvifTrustcktype    = 0x00000800

	// strh dwFlags
	AvisfDisabled        = 0x00000001
	AvisfVideoPalchanges = 0x00010000

	// bIndexType
	AviIndexOfIndexes = 0x00
	AviIndexOfChunks  = 0x01
	AviIndexIsData    = 0x80

	// bIndexSubtype
	AviIndex2field = 0x01

	// idx1 chunk dwFlags
	AviifList      = 0x00000001
	AviifKeyframe  = 0x00000010
	AviifFirstpart = 0x00000020
	AviifLastpart  = 0x00000040
	AviifNotime    = 0x00000100
)

type List struct {
	DwList   [4]byte
	DwSize   uint32
	DwFourCC [4]byte
	Data     []byte
}

type Chunk struct {
	DwFourCC [4]byte
	DwSize   uint32
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

type AVIIndexEntry struct {
	Ckid          [4]byte
	DwFlags       uint32
	DwChunkOffset uint32
	DwChunkLength uint32
}

// DecodeAvi ...
func DecodeAvi(file *bytes.Reader) {
	riff := list{}
	if err := binary.Read(file, binary.LittleEndian, &riff); err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	//fmt.Printf("%+v\n", riff)

	var hdrl bytes.Reader
	var movi bytes.Reader
	var info bytes.Reader
	//var idx1 []byte

	avih := MainAVIHeader{}
	var strl []bytes.Reader
	var vids []struct {
		strh AVIStreamHeader
		strf BitMapInfoHeader
	}
	var auds []struct {
		strh AVIStreamHeader
		strf WaveFormat
	}
	var txts [][]byte

	infoData := make(map[[4]byte]interface{})

	for file.Len() != 0 {
		var fourCC [4]byte
		_ = binary.Read(file, binary.LittleEndian, &fourCC)
		var size uint32
		_ = binary.Read(file, binary.LittleEndian, &size)
		//fmt.Printf("%s\n", fourCC)

		if fourCC == riffList { // LIST
			var type_ [4]byte
			_, _ = file.Read(type_[:])
			//fmt.Printf("\t%s\n", type_)

			data := make([]byte, size-4)
			_, _ = file.Read(data)

			switch type_ {
			case listHdrl: // hdrl
				hdrl = *bytes.NewReader(data)
			case listMovi: // movi
				movi = *bytes.NewReader(data)
			case listInfo: // INFO
				info = *bytes.NewReader(data)
			}

		} else if fourCC == chunkIdx1 { // idx1
			_, _ = io.CopyN(ioutil.Discard, file, int64(size))

		} else if fourCC == chunkJunk { // JUNK
			_, _ = io.CopyN(ioutil.Discard, file, int64(size))
		}
	}

	//fmt.Printf("\n%s\n", hdrl)
	for hdrl.Len() != 0 {
		var fourCC [4]byte
		_ = binary.Read(&hdrl, binary.LittleEndian, &fourCC)
		var size uint32
		_ = binary.Read(&hdrl, binary.LittleEndian, &size)
		//fmt.Printf("%s\n", fourCC)

		if fourCC == chunkAvih { // avih
			if err := binary.Read(&hdrl, binary.LittleEndian, &avih); err != nil {
				fmt.Println("binary.Read failed:", err)
			}

		} else if fourCC == riffList { // LIST
			var type_ [4]byte
			_, _ = hdrl.Read(type_[:])
			fmt.Printf("\t%s\n", type_)

			if type_ == listStrl { // strl
				data := make([]byte, size-4)
				_, _ = hdrl.Read(data)
				fmt.Println(string(data[:]))

				strl = append(strl, *bytes.NewReader(data))

			} else { // may be odml
				_, _ = io.CopyN(ioutil.Discard, &hdrl, int64(size-4))
			}

		} else if fourCC == chunkJunk { // JUNK
			_, _ = io.CopyN(ioutil.Discard, &hdrl, int64(size))
		}
	}
	fmt.Printf("avih %+v\n", avih)

	for _, stream := range strl {
		var fourCC [4]byte
		_ = binary.Read(&stream, binary.LittleEndian, &fourCC)
		var size uint32
		_ = binary.Read(&stream, binary.LittleEndian, &size)
		//fmt.Printf("%s\n", fourCC)

		strh := AVIStreamHeader{}
		if err := binary.Read(&stream, binary.LittleEndian, &strh); err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		//fmt.Printf("strh %+v\n", strh)

		switch strh.FccType {
		case [4]byte{'v', 'i', 'd', 's'}:
			data := struct {
				strh AVIStreamHeader
				strf BitMapInfoHeader
			}{strh: strh}

			for stream.Len() != 0 {
				_ = binary.Read(&stream, binary.LittleEndian, &fourCC)
				_ = binary.Read(&stream, binary.LittleEndian, &size)
				fmt.Println(string(fourCC[:]), size)

				if fourCC == chunkStrf {
					if err := binary.Read(&stream, binary.LittleEndian, &data.strf); err != nil {
						fmt.Println("binary.Read failed:", err)
					}

				} else {
					data := make([]byte, size)
					_, _ = stream.Read(data)
					//fmt.Println(string(other[:]))
				}
			}

			vids = append(vids, data)

		case [4]byte{'a', 'u', 'd', 's'}:
			data := struct {
				strh AVIStreamHeader
				strf WaveFormat
			}{strh: strh}

			for stream.Len() != 0 {
				_ = binary.Read(&stream, binary.LittleEndian, &fourCC)
				_ = binary.Read(&stream, binary.LittleEndian, &size)
				fmt.Println(string(fourCC[:]), size)

				if fourCC == chunkStrf {
					if err := binary.Read(&stream, binary.LittleEndian, &data.strf); err != nil {
						fmt.Println("binary.Read failed:", err)
					}

				} else {
					data := make([]byte, size)
					_, _ = stream.Read(data)
					//fmt.Println(string(other[:]))
				}
			}

			auds = append(auds, data)

		case [4]byte{'t', 'x', 't', 's'}:
			data := make([]byte, size)
			stream.Read(data)
			txts = append(txts, data)
		}
	}
	fmt.Printf("vids %+v\n", vids)
	fmt.Printf("auds %+v\n", auds)
	fmt.Printf("txts %s\n", txts)

	fmt.Println("movi", movi.Len())

	//fmt.Printf("\n%s\n", info)
	for info.Len() != 0 {
		var fourCC [4]byte
		_ = binary.Read(&info, binary.LittleEndian, &fourCC)
		var size uint32
		_ = binary.Read(&info, binary.LittleEndian, &size)
		//fmt.Printf("%s\n", fourCC)

		if fourCC == chunkJunk { // JUNK
			_, _ = io.CopyN(ioutil.Discard, &info, int64(size))

		} else {
			data := make([]byte, size)
			if err := binary.Read(&info, binary.LittleEndian, &data); err != nil {
				fmt.Println("binary.Read failed:", err)
			}
			infoData[fourCC] = data
			//fmt.Printf("\t%s\n", data)
		}
	}
	fmt.Printf("info %s\n", infoData)
}

// Deprecated: test implementation
func OldDecodeAvi(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	AVI := List{}
	copy(AVI.DwList[:], content[:4])
	AVI.DwSize = binary.LittleEndian.Uint32(content[4:8])
	copy(AVI.DwFourCC[:], content[8:12])
	AVI.Data = content[12:]

	if AVI.DwFourCC != formatAvi {
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
			dwFourCC := [4]byte{}
			copy(dwFourCC[:], AVI.Data[index+8:index+12])

			if dwFourCC == listHdrl { // hdrl
				copy(hdrl.DwList[:], dwList)
				hdrl.DwSize = dwSize
				hdrl.DwFourCC = dwFourCC
				hdrl.Data = AVI.Data[index+12 : index+12+int(dwSize)-4]
			} else if dwFourCC == listMovi { // movi
				copy(movi.DwList[:], dwList)
				movi.DwSize = dwSize
				movi.DwFourCC = dwFourCC
				movi.Data = AVI.Data[index+12 : index+12+int(dwSize)-4]
			} else if dwFourCC == listInfo { // INFO
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
		switch strh.FccType {
		case [4]byte{'v', 'i', 'd', 's'}:
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
		case [4]byte{'a', 'u', 'd', 's'}:
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
		case [4]byte{'t', 'x', 't', 's'}:
			fmt.Printf("%+v\n", strh)
		case [4]byte{'m', 'i', 'd', 's'}:
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
