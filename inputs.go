package main

import (
    "encoding/binary"
    "log"
    "io"
    "time"
)

func  fileInput (input_file io.Reader, c chan []complex64) {
    var cpx_buffer = make([]complex64,100000)
    for {
        err := binary.Read(input_file, binary.LittleEndian, cpx_buffer)
        if err == io.EOF {
            c <- cpx_buffer
            time.Sleep(2)
            break
        }
        if err != nil {
            log.Fatal("fileInput: ",err)
        }
        c <- cpx_buffer
    }
}
