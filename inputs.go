package main

//#cgo LDFLAGS: -lrtlsdr
//#include <rtl-sdr.h>
//
//rtlsdr_dev_t * local_rtl_device;
//
import "C"
import (
	"encoding/binary"
	"io"
	"log"
	"time"
    "fmt"
    "unsafe"
)

func fileInput(input_file io.Reader, c chan []complex64) {
	var cpx_buffer = make([]complex64, 100000)
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

func rtlsdrInput(dev_id int, c_freq int, s_rate int, c chan []complex64) {
    //check if there's a device by this id
    var dev_count uint32
    dev_count = uint32(C.rtlsdr_get_device_count())
    if dev_id <0 || dev_id >= int(dev_count) {
        log.Fatal("Wrong device id, please input a number from 0 to",dev_count - 1)
    }
    fmt.Println("Opening device",dev_id,string(*C.rtlsdr_get_device_name(C.uint32_t(dev_id))))
    ret := C.rtlsdr_open(&C.local_rtl_device, C.uint32_t(dev_id))
    if ret < 0 {
        log.Fatal("rtlsdr_open returned",ret)
    }
    defer C.rtlsdr_close(C.local_rtl_device)
    //setting center freq
    if c_freq < 0 {
        log.Fatal("Negative freqs not supported")
    }
    ret = C.rtlsdr_set_center_freq(C.local_rtl_device,C.uint32_t(c_freq))
    if ret < 0 {
        log.Fatal("rtlsdr_set_center_freq returned",ret)
    }
    //setting sample rate
    if s_rate <=0 {
        log.Fatal("Negative or 0 sample rates not supported")
    }
    ret = C.rtlsdr_set_sample_rate(C.local_rtl_device,C.uint32_t(s_rate))
    if ret < 0 {
        log.Fatal("rtlsdr_set_sample_rate returned",ret)
    }
/// quick init
    C.rtlsdr_set_tuner_gain_mode(C.local_rtl_device, 1) //manual tuner gain
    C.rtlsdr_set_agc_mode(C.local_rtl_device, 0) //no AGC on rtl
    C.rtlsdr_set_direct_sampling(C.local_rtl_device, 0) //no direct sampling, we want E4000
    C.rtlsdr_set_offset_tuning(C.local_rtl_device, 0) //offset tuning off
    C.rtlsdr_set_tuner_gain(C.local_rtl_device, 290) //+29dB
    //IF gain total +20dB 15+3+2
    C.rtlsdr_set_tuner_if_gain(C.local_rtl_device, 6,150)
    C.rtlsdr_set_tuner_if_gain(C.local_rtl_device, 5,30)
    C.rtlsdr_set_tuner_if_gain(C.local_rtl_device, 4,20)
    C.rtlsdr_set_tuner_if_gain(C.local_rtl_device, 3,0)    
    C.rtlsdr_set_tuner_if_gain(C.local_rtl_device, 2,0)
    C.rtlsdr_set_tuner_if_gain(C.local_rtl_device, 1,0)

    buffer_cpx := make([]complex64,512*64)
    buffer_byte := make([]uint8, 512*64*2)
    read := int32(0)
    C.rtlsdr_reset_buffer(C.local_rtl_device)
    for {
        ret = C.rtlsdr_read_sync(C.local_rtl_device, (unsafe.Pointer)(&buffer_byte[0]), 512*64*2, (*C.int)(&read))
        if ret < 0 && ret != -9 {
            log.Fatal("rtlsdr_read_sync returned",ret)
        }
        fmt.Println("Read samples:",read/2)
        for i:=int32(0); i<read/2; i++ {
            buffer_cpx[i] = complex((float32(buffer_byte[2*i])-127.4)/128,(float32(buffer_byte[2*i+1])-127.4)/128)
            //fmt.Println(buffer_cpx[i])
        }
        c <- buffer_cpx[:read/2]
    }
}
