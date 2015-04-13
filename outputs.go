package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

func fileOutput(output_file io.Writer, c chan []complex64) {
	var cpx_buffer []complex64
	for {
		cpx_buffer = <-c
		err := binary.Write(output_file, binary.LittleEndian, cpx_buffer)
		if err != nil {
			fmt.Println("fileOutput: ", err)
		}
	}
}
