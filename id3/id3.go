package id3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
)

const (
	// ISO-8859-1. Terminated with $00.
	encodingISO8859 byte = 0x0

	// UTF-16 encoded Unicode with BOM. All strings in the same frame SHALL have the same byte order.
	//Terminated with $00 00.
	encodingUTF16WithBOM byte = 0x1

	// UTF-16BE encoded Unicode without BOM. Terminated with $00 00.
	encodingUTF16 byte = 0x2

	// UTF-8 encoded Unicode. Terminated with $00.
	encodingUTF8 byte = 0x3
)

type ID3Header struct {
	FileIdentifier [3]byte
	Version        uint16
	Flags          uint8
	Size           [4]byte
}

type ID3Frame struct {
	FrameID [4]byte
	Size    uint32
	Flags   uint16
}

func decodeBuffer(tag *bytes.Buffer, n int64) string {
	enc, _ := tag.ReadByte() // text encoding

	var reader io.Reader
	switch enc {
	case encodingISO8859:
		reader = charmap.ISO8859_1.NewDecoder().Reader(io.LimitReader(tag, n-1))

	case encodingUTF16WithBOM:
		encoder := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
		reader = transform.NewReader(io.LimitReader(tag, n-1), encoder.NewDecoder())

	case encodingUTF16:
		encoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		reader = transform.NewReader(io.LimitReader(tag, n-1), encoder.NewDecoder())

	case encodingUTF8:
		// Do nothing.

	default:
		reader = charmap.ISO8859_1.NewDecoder().Reader(io.LimitReader(tag, n-1))
	}

	t, _ := ioutil.ReadAll(reader)

	return string(t)
}

// Decode tag bytes to UTF-8.
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

// ReadID3 ...
func ReadID3(buffer *bytes.Reader) {
	id3Header := ID3Header{}
	if err := binary.Read(buffer, binary.BigEndian, &id3Header); err != nil {
		log.Println(err)
	}
	fmt.Printf("%+v\n", id3Header)

	size := (uint32(id3Header.Size[0]) << 21) | (uint32(id3Header.Size[1]) << 14) |
		(uint32(id3Header.Size[2]) << 7) | uint32(id3Header.Size[3])

	id3Tag := bytes.NewBuffer(make([]byte, 0, size))
	if _, err := id3Tag.ReadFrom(io.LimitReader(buffer, int64(size))); err != nil {
		log.Println(err)
	}

	for id3Tag.Len() != 0 {
		frame := ID3Frame{}
		if err := binary.Read(id3Tag, binary.BigEndian, &frame); err != nil {
			log.Println(err)
		}

		if frame.Size == 0 {
			break
		}

		fmt.Printf("%s, %d, %d\n", frame.FrameID, frame.Size, frame.Flags)

		frameId := string(frame.FrameID[:])
		switch {
		case frameId == "APIC":
			frameData := make([]byte, frame.Size)
			id3Tag.Read(frameData)

		case frameId == "COMM":
			enc, _ := id3Tag.ReadByte() // text encoding

			buf := make([]byte, 3)
			id3Tag.Read(buf)
			lang := decodeByte(enc, buf)

			buf, _ = id3Tag.ReadBytes(0x00)
			desc := decodeByte(enc, buf)

			buf = make([]byte, int(frame.Size)-4+len(desc))
			id3Tag.Read(buf)
			val := decodeByte(enc, buf)

			fmt.Printf("%s %s %s\n", lang, desc, val)

		case frameId[0] == 'T':
			enc, _ := id3Tag.ReadByte() // text encoding

			var reader io.Reader
			switch enc {
			case encodingISO8859:
				reader = charmap.ISO8859_1.NewDecoder().Reader(io.LimitReader(id3Tag, int64(frame.Size)-1))

			case encodingUTF16WithBOM, encodingUTF16:
				encoder := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
				reader = transform.NewReader(io.LimitReader(id3Tag, int64(frame.Size)-1), encoder.NewDecoder())

			case encodingUTF8:
				// Do nothing.

			default:
				reader = charmap.ISO8859_1.NewDecoder().Reader(io.LimitReader(id3Tag, int64(frame.Size)-1))
			}

			t, _ := ioutil.ReadAll(reader)
			fmt.Println(string(t))

		default:
			frameData := make([]byte, frame.Size)
			id3Tag.Read(frameData)
		}
	}
}
