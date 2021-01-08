package riff

var (
	riffId   = [4]byte{'R', 'I', 'F', 'F'}
	riffList = [4]byte{'L', 'I', 'S', 'T'}

	formatWave = [4]byte{'W', 'A', 'V', 'E'}
	formatAvi  = [4]byte{'A', 'V', 'I', ' '}

	listHdrl = [4]byte{'h', 'd', 'r', 'l'}
	listStrl = [4]byte{'s', 't', 'r', 'l'}
	listMovi = [4]byte{'m', 'o', 'v', 'i'}
	listInfo = [4]byte{'I', 'N', 'F', 'O'}

	chunkFmt  = [4]byte{'f', 'm', 't', ' '}
	chunkAvih = [4]byte{'a', 'v', 'i', 'h'}
	chunkStrh = [4]byte{'s', 't', 'r', 'h'}
	chunkStrf = [4]byte{'s', 't', 'r', 'f'}
	chunkIdx1 = [4]byte{'i', 'd', 'x', '1'}
	chunkData = [4]byte{'d', 'a', 't', 'a'}
	chunkJunk = [4]byte{'j', 'u', 'n', 'k'}
)

type list struct {
	DwList   [4]byte
	DwSize   uint32
	DwFourCC [4]byte
}

type chunk struct {
	DwFourCC [4]byte
	DwSize   uint32
}

type WaveFormat struct {
	WFormatTag      uint16 // Format type
	NChannels       uint16 // Number of channels (i.e. mono, stereo...)
	NSamplesPerSec  uint32 // Sample rate
	NAvgBytesPerSec uint32 // For buffer estimation
	NBlockAlign     uint16 // Block size of data
}

type PcmWaveFormat struct {
	Wf             WaveFormat
	WBitsPerSample uint16 // Number of bits per sample
}