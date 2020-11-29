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
	offset  int
	bytes   []byte
	counter int
}

func NewBitReader(b []byte) *BitReader {
	return &BitReader{0, b, 0}
}

func (bitReader *BitReader) ReadBits(n int) int {
	buf := make([]byte, 4)

	offset := bitReader.offset / 8
	copy(buf, bitReader.bytes[offset:])
	//fmt.Printf("%08b", buf)

	r := uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	r = r >> (32 - (n + bitReader.offset%8)) & (0xFFFFFFFF >> (32 - n))

	bitReader.offset += n
	bitReader.counter += n
	return int(r)
}

func (bitReader *BitReader) Seek(n int) {
	bitReader.offset += n
	bitReader.counter += n
}
