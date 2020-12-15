package mpeg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
)

type atom struct {
	AtomSize  uint32
	AtomTType [4]byte
}

func mpeg14(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	index := 0
	for index != len(content) {
		atom := atom{}
		r := bytes.NewReader(content[index : index+8])
		if err := binary.Read(r, binary.BigEndian, &atom); err != nil {
			fmt.Println("binary.Read failed:", err)
		}

		data := content[index+8 : index+int(atom.AtomSize)]

		index += int(atom.AtomSize)

		fmt.Printf("%d %s\n", atom.AtomSize, atom.AtomTType)

		if string(atom.AtomTType[:]) == "ftyp" {
			fmt.Println(string(data))
		} else if string(atom.AtomTType[:]) == "moov" {
			fmt.Println(string(data))
		}
	}
}

var Mp4 = mpeg14
