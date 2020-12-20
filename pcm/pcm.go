package pcm

import (
	"math"
	"time"
)

type Context struct {
	SampleRate int // sampling frequency
	Channels   int // numbers of channels
}

// Samples ...
//
// Len return pcm total count.
type Samples interface {
	Pcm() interface{}
	Len() int
	BitPerSample() int

	Context() *Context

	Duration() time.Duration
}

func ToS16LE(s Samples) *S16LE {
	s16le := &S16LE{}
	s16le.Context().SampleRate = s.Context().SampleRate
	s16le.Context().Channels = s.Context().Channels

	switch t := s.Pcm().(type) {
	case []float32:
		s16le.pcm = make([]int16, len(t))
		for i, v := range t {
			s := int(v * math.MaxInt16)
			if s > math.MaxInt16 {
				s = math.MaxInt16
			} else if s < -math.MaxInt16 {
				s = -math.MaxInt16
			}

			s16le.pcm[i] = int16(s)
		}
	}

	return s16le
}
