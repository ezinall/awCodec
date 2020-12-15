package main

import (
	"awCodec/g711"
	"awCodec/mpeg"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	content, err := ioutil.ReadFile("Lumen - Государство.mp3")
	if err != nil {
		log.Fatal(err)
	}
	file := bytes.NewReader(content)
	out, _ := mpeg.Mp3(file)
	ulawOut := make([]uint8, len(out.Pcm))
	for i, v := range out.Pcm {
		ulawOut[i] = g711.ULawEncode(v)
	}

	outFile, _ := os.Create("outMulaw.raw")
	for _, v := range ulawOut {

		if err := binary.Write(outFile, binary.LittleEndian, v); err != nil {
			fmt.Println(err)
		}
	}

	outFile2, _ := os.Create("outs16le.raw")
	for _, v := range out.Pcm {
		if err := binary.Write(outFile2, binary.LittleEndian, v); err != nil {
			fmt.Println(err)
		}
	}
}
