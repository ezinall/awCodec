package id3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/language"
	"golang.org/x/text/transform"
	"io"
)

var tag = []byte{'T', 'A', 'G'}
var fileIdentifier = []byte{'I', 'D', '3'}

// Flags
const (
	unSync = 0x01 // Unsynchronisation
	ext    = 0x02 // Extended header
	exp    = 0x04 // Experimental indicator
	footer = 0x06 // Footer present
)

const (
	// ISO-8859-1. Terminated with $00.
	encodingISO8859 byte = 0x00

	// UTF-16 encoded Unicode with BOM. All strings in the same frame SHALL have the same byte order.
	// Terminated with $00 00.
	encodingUTF16WithBOM byte = 0x01

	// UTF-16BE encoded Unicode without BOM. Terminated with $00 00.
	encodingUTF16 byte = 0x02

	// UTF-8 encoded Unicode. Terminated with $00.
	encodingUTF8 byte = 0x03
)

type id3v2Header struct {
	FileIdentifier [3]byte
	Version        uint16
	Flags          uint8
	Size           [4]byte
}

type extHeader struct {
	HeaderSize int32
	//Flags      int16
	//Padding    int32
}

type frameHeader struct {
	FrameID [4]byte
	Size    uint32
	Flags   uint16
}

// Decode tag bytes to string.
func decodeByte(enc byte, b []byte) string {
	var r []byte

	switch enc {
	case encodingISO8859:
		r, _ = charmap.ISO8859_1.NewDecoder().Bytes(b)

	case encodingUTF16WithBOM:
		encoder := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
		r, _, _ = transform.Bytes(encoder.NewDecoder(), b)

	case encodingUTF16:
		encoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		r, _, _ = transform.Bytes(encoder.NewDecoder(), b)

	case encodingUTF8:
		// Do nothing.

	default:
		r, _ = charmap.ISO8859_1.NewDecoder().Bytes(b)
	}

	return string(r)
}

func findLenWithTerm(enc byte, b []byte) (l int) {
	if enc == encodingUTF16WithBOM || enc == encodingUTF16 {
		l = bytes.Index(b, []byte{0x00, 0x00})
		if l != -1 {
			l += 2
		}
	} else {
		l = bytes.Index(b, []byte{0x00})
		if l != -1 {
			l += 1
		}
	}
	return
}

// ReadID3 ...
func ReadID3(file []byte) ([]byte, error) {
	if bytes.Equal(file[:3], fileIdentifier) {
		r := bytes.NewReader(file)

		header := id3v2Header{}
		if err := binary.Read(r, binary.BigEndian, &header); err != nil {
			return file, err
		}
		//fmt.Printf("%+v\n", header)

		size := (uint32(header.Size[0]) << 21) | (uint32(header.Size[1]) << 14) |
			(uint32(header.Size[2]) << 7) | uint32(header.Size[3])

		id3Tag := bytes.NewBuffer(make([]byte, 0, size))
		if _, err := id3Tag.ReadFrom(io.LimitReader(r, int64(size))); err != nil {
			return nil, err
		}

		if header.Flags&ext == ext {
			extHeader := extHeader{}
			_ = binary.Read(id3Tag, binary.BigEndian, &extHeader)
			id3Tag.Next(int(extHeader.HeaderSize))
		}

		for id3Tag.Len() != 0 {
			fh := frameHeader{}
			_ = binary.Read(id3Tag, binary.BigEndian, &fh)

			if fh.Size == 0 {
				break
			}

			//fmt.Printf("%s, %d, %d\n", fh.FrameID, fh.Size, fh.Flags)

			frameId := string(fh.FrameID[:])

			buf := make([]byte, fh.Size)
			n, _ := id3Tag.Read(buf)
			if n != int(fh.Size) {
				break
			}
			frame := bytes.NewBuffer(buf)

			if frameId == "APIC" {
				enc, _ := frame.ReadByte() // Text encoding

				buf, _ := frame.ReadBytes(0x00)
				mimeType := decodeByte(enc, buf) // MIME type

				picType, _ := frame.ReadByte() // Picture type

				l := findLenWithTerm(enc, frame.Bytes())
				desc := decodeByte(enc, frame.Next(l)) // Description

				_ = frame.Next(int(fh.Size) - 1 - len(buf) - l) // Picture data

				fmt.Println(frameId, mimeType, picType, desc)

			} else if frameId == "COMM" {
				enc, _ := frame.ReadByte() // Text encoding

				buf := frame.Next(3)
				base, _ := language.ParseBase(string(buf))
				lang := base.ISO3() // Language

				l := findLenWithTerm(enc, frame.Bytes())
				desc := decodeByte(enc, frame.Next(l)) // Short content descrip.

				val := decodeByte(enc, frame.Next(int(fh.Size)-4-l)) // The actual text

				fmt.Printf("%s lang:%s desc:%s val:%s\n", frameId, lang, desc, val)

			} else if frameId[0] == 'T' {
				enc, _ := frame.ReadByte() // Text encoding

				buf := frame.Next(int(fh.Size - 1))

				t := decodeByte(enc, buf)

				fmt.Println(frameId, t)

			} else {
				_ = frame.Next(int(fh.Size)) // frameData
			}
		}

		// 10 byte header length
		return file[size+10:], nil

	} else if bytes.Equal(file[len(file)-128:len(file)-125], tag) {
		frame := bytes.NewBuffer(file[len(file)-128:])
		songName := frame.Next(30)
		artist := frame.Next(30)
		albumName := frame.Next(30)
		year := frame.Next(4)
		comment := frame.Next(30)
		songGenreIdentifier := frame.Next(1)
		fmt.Println(songName, artist, albumName, year, comment, songGenreIdentifier)

		return file[:len(file)-128], nil
	}
	return file, nil
}
