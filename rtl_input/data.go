package rtl_input

//#cgo LDFLAGS: -lrtlsdr
//#include "rtl-input.h"
import "C"
import (
	"log"
)

func (dev *RtlsdrDevice) ResetBuffer() {
    //return value not documented, assuming sucess
    C.rtlinput_set_agc_mode(dev.devPtr)
}

//TODO func (dev *RtlsdrDevice) GetSamples() []int32 {
