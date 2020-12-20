package main

import (
	"awCodec/g711"
	"awCodec/mpeg"
	"awCodec/pcm"
	"awCodec/riff"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	content, err := ioutil.ReadFile("Встань, страх преодолей.mp3")
	if err != nil {
		log.Fatal(err)
	}
	file := bytes.NewReader(content)
	out, _ := mpeg.Mp3(file)

	fmt.Println(out.Duration())

	out2 := pcm.ToS16LE(out)
	out3 := g711.ToMuLaw(out2)

	//riff.EncodeWav(out, riff.WaveFormatIeeeFloat)
	//riff.EncodeWav(out2, riff.WaveFormatPcm)
	riff.EncodeWav(out3, riff.WaveFormatMulaw)
}
