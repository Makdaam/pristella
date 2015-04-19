package rtl_input

//#cgo LDFLAGS: -lrtlsdr
//#include "rtl-input.h"
import "C"
import (
    "unsafe"
    "log"
)

type Rtlsdr_device struct {
    dev_id uint32
    dev_ptr unsafe.Pointer
}

var open_devices []Rtlsdr_device

func Get_device_names() []string {
    var device_names []string
    single_name := make([]uint8,100,100)
    var name_length uint32

    dev_count := uint32(C.rtlinput_get_device_count())

    for i:=uint32(0); i<dev_count; i++ {
        C.rtlinput_get_name(C.uint32_t(i), (*C.uint32_t)(&name_length), (*C.uint8_t)(&single_name[0]))
        if name_length < 100 {
            var name string
            for _, k:= range single_name {
                if k==0 {
                    break
                }
                name+=string(k)
            }
            device_names = append(device_names, name)
        } else {
            log.Panic("Device name too long:", i, name_length)
        }
    }
    return device_names
}

func Open_rtlsdr(dev_id uint32) Rtlsdr_device {
    var dev Rtlsdr_device
    dev.dev_id = dev_id 
   
    dev_count := uint32(C.rtlsdr_get_device_count())
    if dev_count == 0 {
        log.Panic("No devices connected")
    }
    if dev_id >= uint32(dev_count) {
        log.Panic("Wrong device id, please input a number from 0 to",dev_count - 1)
    }
    for _,item := range open_devices {
        if item.dev_id == dev_id {
            log.Println("Device",dev_id,"already open")
            panic("Device already open")
        }
    }
    retval := int32(C.rtlinput_open(&dev.dev_ptr, (C.uint32_t)(dev_id)))
    if retval <0 {
        log.Println("Error while opening: rtlsdr_open:", retval)
    }
    open_devices = append(open_devices, dev)
    return dev
}

func (dev *Rtlsdr_device) Close() {
    retval := int32(C.rtlinput_close(dev.dev_ptr))
    if retval < 0 {
        log.Panic("Closing of device ",dev.dev_id,"returned ",retval)
    }
}

func (dev *Rtlsdr_device) SetCenterFreq(freq uint32) {
    retval := int32(C.rtlinput_set_center_freq(dev.dev_ptr, (C.uint32_t)(freq)))
    if retval < 0 {
        log.Panic("Set center freq of device ",dev.dev_id,"returned ",retval)
    }
}

func (dev *Rtlsdr_device) SetSampleRate(s_rate uint32) {
    retval := int32(C.rtlinput_set_sample_rate(dev.dev_ptr, (C.uint32_t)(s_rate)))
    if retval < 0 {
        log.Panic("Set sample rate of device ",dev.dev_id,"returned ",retval)
    }
}
