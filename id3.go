package main

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
	encodingISO8859      byte = 0x0
	encodingUTF16WithBOM byte = 0x1
	encodingUTF16        byte = 0x2
	encodingUTF8         byte = 0x3
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
			frameData := make([]byte, frame.Size)
			id3Tag.Read(frameData)

		case frameId[0] == 'T':
			enc, _ := id3Tag.ReadByte() // text encoding

			var reader io.Reader
			switch enc {
			case encodingISO8859:
				reader = charmap.ISO8859_1.NewDecoder().Reader(io.LimitReader(id3Tag, int64(frame.Size)-1))
			case encodingUTF16WithBOM:
				win16be := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
				reader = transform.NewReader(io.LimitReader(id3Tag, int64(frame.Size)-1), win16be.NewDecoder())
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
