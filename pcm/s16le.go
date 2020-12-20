package pcm

import (
	"time"
)

type S16LE struct {
	pcm     []int16
	context Context
}

func (s16le *S16LE) Append(elements interface{}) {
	switch n := elements.(type) {
	case []int16:
		s16le.pcm = append(s16le.pcm, n...)
	}
}

func (s16le *S16LE) Pcm() interface{} {
	return s16le.pcm
}

// Len return total samples count.
func (s16le *S16LE) Len() int {
	return len(s16le.pcm)
}

func (s16le *S16LE) BitPerSample() int {
	return 16
}

func (s16le *S16LE) Context() *Context {
	return &s16le.context
}

func (s16le *S16LE) Duration() time.Duration {
	duration := len(s16le.pcm) / s16le.context.Channels / s16le.context.SampleRate
	return time.Duration(duration) * time.Second
}
