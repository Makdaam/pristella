package rtl_input

//#cgo LDFLAGS: -lrtlsdr
//#include "rtl-input.h"
import "C"
import (
	"log"
)

//freq stuff
func (dev *RtlsdrDevice) SetCenterFreq(freq uint32) {
	retval := int32(C.rtlinput_set_center_freq(dev.devPtr, (C.uint32_t)(freq)))
	if retval < 0 {
		log.Panic("Set center freq of device ", dev.devId, "returned ", retval)
	}
}

func (dev *RtlsdrDevice) GetCenterFreq() uint32 {
	retval := uint32(C.rtlinput_get_center_freq(dev.devPtr))
	if retval == 0 {
		log.Panic("Get center freq of device ", dev.devId, "returned ", retval)
	}
	return retval
}

func (dev *RtlsdrDevice) SetFrequencyCorrection(ppm int32) {
	retval := int32(C.rtlinput_set_freq_correction(dev.devPtr, (C.int32_t)(ppm)))
	if retval != 0 {
		log.Panic("Set freq correction of device ", dev.devId, "returned ", retval)
	}
}

func (dev *RtlsdrDevice) GetFrequencyCorrection() int32 {
	return int32(C.rtlinput_get_freq_correction(dev.devPtr))
}

func (dev *RtlsdrDevice) SetSampleRate(sampleRate uint32) {
	retval := int32(C.rtlinput_set_sample_rate(dev.devPtr, (C.uint32_t)(sampleRate)))
	if retval < 0 {
		log.Panic("Set sample rate of device ", dev.devId, "returned ", retval)
	}
}

func (dev *RtlsdrDevice) GetSampleRate() uint32 {
	retval := uint32(C.rtlinput_get_sample_rate(dev.devPtr))
	if retval == 0 {
		log.Panic("Get sample rate of device ", dev.devId, "returned ", retval)
	}
	return retval
}
