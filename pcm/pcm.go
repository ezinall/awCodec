package pcm

import (
	"math"
	"time"
)

// Context contains information about the audio stream.
type Context struct {
	SampleRate int // Sampling frequency defines how many times per second a sound is sampled.
	Channels   int // Numbers of channels.
}

// Copy Context.
func (c *Context) Copy(src Context) {
	*c = src
}

// Samples ...
type Samples interface {
	// Pcm ...
	Pcm() interface{}

	// Len returns pcm total count samples.
	Len() int

	// BitPerSample returns counts bits per one sample.
	BitPerSample() int

	// Context returns pointer to stream context.
	Context() *Context

	// Duration represents the total time in seconds as time.Duration.
	Duration() time.Duration
}

func ToS16LE(s Samples) *S16LE {
	s16le := &S16LE{}
	s16le.Context().Copy(*s.Context())

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
