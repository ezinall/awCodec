package g711

import (
	"awCodec/pcm"
	"time"
)

var muLawExponentTable = [256]int16{
	0, 0, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
}

const (
	muLawClip = 0x7F7B
	muLawBias = 0x84
)

// MuLawEncode encodes a 16bit LPCM sample to u-law PCM.
func MuLawEncode(sample int16) uint8 {
	sign := (sample >> 8) & 0x80
	if sign != 0 {
		sample = -sample
	}
	if sample > muLawClip {
		sample = muLawClip
	}

	sample = sample + muLawBias
	exponent := muLawExponentTable[(sample>>7)&0xFF]
	mantissa := (sample >> (exponent + 3)) & 0x0F
	return uint8(^(sign | (exponent << 4) | mantissa))
}

type MuLaw struct {
	pcm     []uint8
	context pcm.Context
}

func (muLaw *MuLaw) Pcm() interface{} {
	return muLaw.pcm
}

func (muLaw *MuLaw) Len() int {
	return len(muLaw.pcm)
}

func (muLaw *MuLaw) BitPerSample() int {
	return 8
}

func (muLaw *MuLaw) Duration() time.Duration {
	duration := len(muLaw.pcm) / muLaw.context.Channels / muLaw.context.SampleRate
	return time.Duration(duration) * time.Second
}

func (muLaw *MuLaw) Context() *pcm.Context {
	return &muLaw.context
}

func ToMuLaw(s pcm.Samples) *MuLaw {
	muLaw := &MuLaw{}
	muLaw.Context().SampleRate = s.Context().SampleRate
	muLaw.Context().Channels = s.Context().Channels

	switch t := s.Pcm().(type) {
	case []int16:
		muLaw.pcm = make([]uint8, len(t))
		for i, v := range t {
			s := MuLawEncode(v)

			muLaw.pcm[i] = s
		}
	}

	return muLaw
}
