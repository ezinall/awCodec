package pcm

import (
	"time"
)

type F32LE struct {
	pcm     []float32
	context Context
}

func (f32le *F32LE) Append(elements interface{}) {
	switch n := elements.(type) {
	case []float32:
		f32le.pcm = append(f32le.pcm, n...)
	}
}

func (f32le *F32LE) Pcm() interface{} {
	return f32le.pcm
}

// Len return total samples count.
func (f32le *F32LE) Len() int {
	return len(f32le.pcm)
}

func (f32le *F32LE) BitPerSample() int {
	return 32
}

func (f32le *F32LE) Context() *Context {
	return &f32le.context
}

func (f32le *F32LE) Duration() time.Duration {
	duration := len(f32le.pcm) / f32le.context.Channels / f32le.context.SampleRate
	return time.Duration(duration) * time.Second
}
