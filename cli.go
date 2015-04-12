package main

import (
	"flag"
	"fmt"
    "os"
    "log"
    "runtime/pprof"
)

var dev_number int
var center_freq int
var offset_freq int
var sample_rate int
var input_filename string
var output_filename string
var cpu_profile string

func init() {
	//flag definition (CLI arguments)
	flag.IntVar(&dev_number, "d", 0, "number of RTL-SDR device (counting from 0)")
	flag.StringVar(&output_filename, "o", "/tmp/fifo1", "file name of the output fifo")
	flag.StringVar(&input_filename, "i", "", "input file name - for debug only")
	flag.StringVar(&cpu_profile, "p", "", "cpu_profile file - for debug only")
	flag.IntVar(&center_freq, "f", 429000, "center frequency in kHz")
	flag.IntVar(&offset_freq, "a", 134000, "offset frequency in Hz")
	flag.IntVar(&sample_rate, "s", 1000000, "sample rate in samples per second")
	flag.Parse() //call after all flags are defined
}

func main() {
	fmt.Println("Pristella Tetra - a lightweight TELive radio frontent")
    if cpu_profile != "" {
        f, err := os.Create(cpu_profile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    fmt.Println(" - opening output")
    output_file, err := os.OpenFile(output_filename,os.O_WRONLY, 0200)
    if err!=nil {
        log.Fatal("Can't open output file: ",err)
    }
    defer output_file.Close()

    output_chan := make(chan []complex64)
    input_chan := make(chan []complex64)

    go fileOutput(output_file, output_chan)

    fmt.Println(" - preparing filters")
    go dsp(input_chan,output_chan)
    
    fmt.Println(" - opening input")
	if input_filename != "" {
		fmt.Println("WARNING! DEBUG MODE")
        input_file, err := os.Open(input_filename)
        if err!=nil {
            log.Fatal("Can't open input file: ",err)
        }
        defer input_file.Close()
        fileInput(input_file, input_chan)
        
	} else {
		fmt.Println("Opening RTL-SDR")
		//run a device input
	}
}
