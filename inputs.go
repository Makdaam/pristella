package main

import "C"
import (
	"encoding/binary"
	"fmt"
	"github.com/Makdaam/pristella/rtl_input"
	"io"
	"log"
	"time"
)

func fileInput(input_file io.Reader, c chan []int32) {
	var cpx_buffer = make([]int32, 100000)
	for {
		err := binary.Read(input_file, binary.LittleEndian, cpx_buffer)
		if err == io.EOF {
			c <- cpx_buffer
			time.Sleep(2)
			break
		}
		if err != nil {
			log.Fatal("fileInput: ", err)
		}
		c <- cpx_buffer
	}
}

func rtlsdrInput(dev_id int, c_freq int, s_rate int, c chan []int32, corr int) {
	names := rtl_input.GetRtlsdrNames()
	var dev rtl_input.RtlsdrDevice
	fmt.Println(">>>>>>>>>>>>>>>>")
	for i, j := range names {
		fmt.Println(i, j)
	}
	fmt.Println(">>>>>>>>>>>>>>>>")
	//make sure id is positive
	if dev_id < 0 {
		log.Panic("Device id has to be >= 0")
	}
	log.Println("Opening device", dev_id)

	dev = rtl_input.OpenRtlsdr(uint32(dev_id))
	defer dev.Close()

	//setting center freq
	if c_freq < 0 {
		log.Fatal("Negative freqs not supported")
	}
	dev.SetCenterFreq(uint32(c_freq))
    //setting sample rate
    if s_rate <=0 {
       log.Fatal("Negative or 0 sample rates not supported")
    }
    dev.SetSampleRate(uint32(s_rate))
    dev.SetFrequencyCorrection(int32(corr))
    /// quick and dirty init
    dev.SetTunerGainMode(true) //manual tuner gain
    dev.SetAGCMode(false) //no AGC on rtl
    dev.SetDirectSampling(false, false) //no direct sampling, we want E4000
    dev.SetTunerGain(290) //+29dB
    dev.SetIFGain(200) //+29dB
    dev.ResetBuffer()
    for {
        c <- dev.GetSamples()
    }
}
