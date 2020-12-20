package utils

type BitReader struct {
	offset  int
	bytes   []byte
	Counter int
}

func NewBitReader(b []byte) *BitReader {
	return &BitReader{0, b, 0}
}

func (bitReader *BitReader) ReadBits(n int) int {
	buf := make([]byte, 4)

	offset := bitReader.offset / 8
	copy(buf, bitReader.bytes[offset:])

	r := uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	r = r >> (32 - (n + bitReader.offset%8)) & (0xFFFFFFFF >> (32 - n))

	bitReader.offset += n
	bitReader.Counter += n
	return int(r)
}

func (bitReader *BitReader) Seek(n int) {
	bitReader.offset += n
	bitReader.Counter += n
}
