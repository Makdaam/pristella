package rtl_input

//#cgo LDFLAGS: -lrtlsdr
//#include "rtl-input.h"
import "C"
import (
	"log"
	"unsafe"
)

type RtlsdrDevice struct {
	devId  uint32
	devPtr unsafe.Pointer
}

const (
	RTLSDR_TUNER_UNKNOWN = iota
	RTLSDR_TUNER_E4000
	RTLSDR_TUNER_FC0012
	RTLSDR_TUNER_FC0013
	RTLSDR_TUNER_FC2580
	RTLSDR_TUNER_R820T
	RTLSDR_TUNER_R828D
)

var openDevices []RtlsdrDevice

func GetRtlsdrNames() []string {
	var deviceNames []string
	singleName := make([]uint8, 100, 100)
	var nameLength uint32

	devCount := uint32(C.rtlinput_get_device_count())

	for i := uint32(0); i < devCount; i++ {
		C.rtlinput_get_name(C.uint32_t(i), (*C.uint32_t)(&nameLength), (*C.uint8_t)(&singleName[0]))
		if nameLength < 100 {
			var name string
			for _, k := range singleName {
				if k == 0 {
					break
				}
				name += string(k)
			}
			deviceNames = append(deviceNames, name)
		} else {
			log.Panic("Device name too long:", i, nameLength)
		}
	}
	return deviceNames
}

func GetRtlsdrIdFromSerial(serial []uint8) int32 {
    zeroTerminated := false
    for _,i := range serial {
        if i == 0 {
            zeroTerminated = true
            break
        }
    }
    if zeroTerminated {
        return int32(C.rtlinput_get_index_from_serial((*C.uint8_t) (&serial[0])))
    } else {
        log.Panic("GetIdBySerial failed, serial not a zero terminated C string")
    }
    return -1
}

func OpenRtlsdr(devId uint32) RtlsdrDevice {
	var dev RtlsdrDevice
	dev.devId = devId

	devCount := uint32(C.rtlsdr_get_device_count())
	if devCount == 0 {
		log.Panic("No devices connected")
	}
	if devId >= uint32(devCount) {
		log.Panic("Wrong device id, please input a number from 0 to", devCount-1)
	}
	for _, item := range openDevices {
        //TODO this needs to be changed to check serials, ids can change
		if item.devId == devId {
			log.Println("Device", devId, "already open")
			panic("Device already open")
		}
	}
	retval := int32(C.rtlinput_open(&dev.devPtr, (C.uint32_t)(devId)))
	if retval < 0 {
		log.Println("Error while opening: rtlsdrOpen:", retval)
	}
	openDevices = append(openDevices, dev)
	return dev
}

func (dev *RtlsdrDevice) Close() {
	retval := int32(C.rtlinput_close(dev.devPtr))
	if retval < 0 {
		log.Panic("Closing of device ", dev.devId, "returned ", retval)
	}
}

func (dev *RtlsdrDevice) GetTunerType() int32 {
	retval := int32(C.rtlinput_get_tuner_type(dev.devPtr))
	if retval == 0 {
		log.Panic("Get tuner type of device ", dev.devId, "returned ", retval)
	}
	return retval
}

func (dev *RtlsdrDevice) GetSerial() []uint8 {
	byteBuffer := make([]uint8, 257, 257)
	C.rtlinput_get_serial(dev.devPtr, (*C.uint8_t)(&byteBuffer[0]))
	for k, i := range byteBuffer {
		if i == 0 {
            return byteBuffer[:k+1]
		}
	}
    return make([]uint8,0,0)
}
