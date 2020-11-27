package main

import (
	"io"
)

type PlainBitReader struct {
	reader io.ByteReader
	byte   byte
	offset uint8
}

// Simple method to return a pointer to new instance
func NewPlainBitReader(reader io.ByteReader) *PlainBitReader {
	return &PlainBitReader{reader, 0, 0}
}

// Read one bit and return boolean result and error.
func (bitReader *PlainBitReader) ReadBit() (bool, error) {
	if bitReader.offset == 8 {
		bitReader.offset = 0
	}

	// Get next byte
	if bitReader.offset == 0 {
		var err error

		bitReader.byte, err = bitReader.reader.ReadByte()

		if err != nil {
			return false, err
		}
	}

	// Compare current byte to 1000 0000 shifted right bitReader.offset times
	bit := bitReader.byte & (0x80 >> bitReader.offset)
	// Increment our offset
	bitReader.offset++
	// Comparison will turn byte to boolean, and no error is returned
	return bit != 0, nil
}

// Read bits and return uint64 value and error.
func (bitReader *PlainBitReader) ReadBits(bits int64) (uint64, error) {
	var bitRange uint64

	for i := bits - 1; i >= 0; i-- {
		bit, err := bitReader.ReadBit()

		if err != nil {
			return 0, err
		}

		if bit {
			bitRange |= 1 << uint64(i)
		}
	}

	return bitRange, nil
}

func (bitReader *PlainBitReader) MustReadBits(bits int64) uint64 {
	var bitRange uint64

	for i := bits - 1; i >= 0; i-- {
		bit, err := bitReader.ReadBit()

		if err != nil {
			return 0
		}

		if bit {
			bitRange |= 1 << uint64(i)
		}
	}

	return bitRange
}

// Old
type BitReader struct {
	offset int
	bytes  []byte
}

func NewBitReader(b []byte) *BitReader {
	return &BitReader{0, b}
}

func (si *BitReader) ReadBits(n int) int {
	buf := make([]byte, 5)

	offset := si.offset / 8
	copy(buf, si.bytes[offset:])
	//fmt.Printf("%08b", buf)

	r := uint64(buf[0])<<32 | uint64(buf[1])<<24 | uint64(buf[2])<<16 | uint64(buf[3])<<8 | uint64(buf[4])
	r = r >> (40 - (n + si.offset%8)) & (0xFFFFFFFFFFFFFFFF >> (64 - n))

	si.offset += n

	return int(r)
}
