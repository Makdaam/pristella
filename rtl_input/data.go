package rtl_input

//#cgo LDFLAGS: -lrtlsdr
//#include "rtl-input.h"
import "C"
import (
	"log"
    "unsafe"
)

func (dev *RtlsdrDevice) ResetBuffer() {
    //return value not documented, assuming sucess
    C.rtlinput_reset_buffer(dev.devPtr)
}

func (dev *RtlsdrDevice) GetSamples() []int32 {
    bufferLength := 32768
    read := int32(0)
    byteBuffer := make([]uint8, bufferLength*2,bufferLength*2)
    intBuffer := make([]int32, bufferLength*2,bufferLength*2)
    ret := C.rtlinput_read_sync(dev.devPtr, (unsafe.Pointer)(&byteBuffer[0]), (C.int32_t)(bufferLength*2), (*C.int32_t)(&read))
    if ret < 0 {
        log.Panic("rtlsdr_read_sync returned",ret)
    }
    //changing to interop fixed point
    for i:= int32(0); i<read; i++ {
        intBuffer[i]=(int32(byteBuffer[i]) << 8) - 32614 // shift for fixed point and -127.4
    }
    return intBuffer
}
