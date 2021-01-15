package main

import (
	"awCodec/riff"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	content, err := ioutil.ReadFile("_dacha.avi")
	if err != nil {
		log.Fatal(err)
	}
	file := bytes.NewReader(content)
	//out, _ := mpeg.DecodeMp3(file)

	//fmt.Println(out.Duration())

	//out2 := pcm.ToS16LE(out)
	//out3 := g711.ToMuLaw(out2)

	//riff.EncodeWave(riff.WaveFormatIeeeFloat, out)
	//riff.EncodeWave(out2, riff.WaveFormatPcm)
	//riff.EncodeWave(out3, riff.WaveFormatMulaw)

	out := riff.DecodeAvi(file)
	fmt.Println(out.Duration(), out.Context())
	riff.EncodeWave(riff.WaveFormatPcm, &out)
}
