package rtl_input

//#cgo LDFLAGS: -lrtlsdr
//#include "rtl-input.h"
import "C"
import (
	"log"
)

func (dev *RtlsdrDevice) SetTunerGainMode(manualGain bool) {
    var manualInt int32
    if manualGain {
        manualInt = 1
    } else {
        manualInt = 0
    }
    retval := int32(C.rtlinput_set_tuner_gain_mode(dev.devPtr, (C.int32_t)(manualInt)))
	if retval != 0 {
		log.Panic("Set gain mode for dev ", dev.devId, "returned ", retval)
	}
}

func (dev *RtlsdrDevice) SetAGCMode(enabled bool) {
    var enabledInt uint32
    if enabled {
        enabledInt = 1
    } else {
        enabledInt = 0
    }
    retval := int32(C.rtlinput_set_agc_mode(dev.devPtr, (C.uint32_t)(enabledInt)))
	if retval != 0 {
		log.Panic("Set AGC mode for dev ", dev.devId, "returned ", retval)
	}
}

func (dev *RtlsdrDevice) SetDirectSampling(directI bool, directQ bool) {
    if (directI && directQ) {
        log.Panic("Can't enable direct sampling both on I and Q channels")
    }
    var enabledInt uint32
    if directI {
        enabledInt = 1
    }
    if directI {
        enabledInt = 2
    }
	retval := int32(C.rtlinput_set_agc_mode(dev.devPtr, (C.uint32_t)(enabledInt)))
	if retval != 0 {
		log.Panic("Set direct mode for dev ", dev.devId, "returned ", retval)
	}
}

func (dev *RtlsdrDevice) GetDirectSampling() (directI bool,directQ bool) {
	retval := int32(C.rtlinput_get_direct_sampling(dev.devPtr))
    if retval == 0 {
        return false, false
    }
    if retval == 1 {
        return true, false
    }
    if retval == 2 {
        return false, true
    }
    log.Panic("Get direct mode of device ", dev.devId, "returned ", retval)
	return false,false
}

func getSquareDistance(a int32, b int32) int32 {
    if a < 32000 && b < 32000 {
        return (a-b) *(a-b)
    } else {
        log.Panic("Integer values too large for comparison")
        return 0
    }
}

//accepts gain in tenths of dB, so value of 30 is 3dB
func (dev *RtlsdrDevice) SetTunerGain(gain int32) (actualGain int32) {
    possibleGains := make([]int32,256,256)
    var possibleGainsCount int32
	retval := int32(C.rtlinput_get_tuner_gains(dev.devPtr,(*C.int32_t)(&possibleGains[0]),(*C.int32_t)(&possibleGainsCount)))
    if retval == -255 {
        log.Panic("This device supports many more gain values than I expected, please contact the developer with your device model and number of gains =",possibleGainsCount)
    }
    if retval != 0 {
        log.Panic("Get possible gain values of device ", dev.devId, "returned ", retval)
    }
    possibleGains = possibleGains[:possibleGainsCount]
    actualGain = possibleGains[0]
    for _,g := range possibleGains {
        if getSquareDistance(actualGain,gain) > getSquareDistance(g,gain) {
            actualGain = g
        }
    }
    //actualGain is the nearest value from possible gains to the gain specified by the user    
	retval = int32(C.rtlinput_set_tuner_gain(dev.devPtr,(C.int32_t)(actualGain)))
    if retval != 0 {
        log.Panic("Set tuner gain of device ", dev.devId, "returned ", retval)
    }
    return actualGain
}

func (dev *RtlsdrDevice) GetTunerGain() int32 {
	retval := int32(C.rtlinput_get_tuner_gain(dev.devPtr))
    if retval == 0 {
        log.Panic("Get tuner gain of device ", dev.devId, "returned ", retval)
    }
    return retval
}

func IfGainsE4000(gain int32) (result []int32) {
//this function's code is based on OsmoSDRs implementation on setting IF gain for E4000
//you can find the original in the OsmoSDR project in lib/rtl/rtl_source_c.cc
//it has some minor changes based on experience with my RTLSDR and this might not apply to yours
    roundedGain := gain - gain % 30

    result = make([]int32,6,6)
    result[0] = -30 //stage 1 is noisy if set to 90, so not using it at all

    sum := int32(-30+30+30) //lowest values for all gain stages
    //gain stage 6 3.0dB - 15.0dB, step 3.0, the minimal value is already accounted for in the sum
    for i:=int32(0); i<=(150-30); i+=30 {
        if getSquareDistance(sum + result[5],roundedGain) > getSquareDistance(sum + i,roundedGain) {
            result[5] = i
        }
    }
    sum += result[5]
    //gain stage 5 3.0dB - 15.0dB, step 3.0, the minimal value is already accounted for in the sum
    for i:=int32(0); i<=(150-30); i+=30 {
        if getSquareDistance(sum + result[4],roundedGain) > getSquareDistance(sum + i,roundedGain) {
            result[4] = i
        }
    }
    sum += result[4]
    //gain stage 3 0.0dB - 9.0dB, step 3.0
    for i:=int32(0); i<=90; i+=30 {
        if getSquareDistance(sum + result[2],roundedGain) > getSquareDistance(sum + i,roundedGain) {
            result[2] = i
        }
    }
    sum+=result[2]
    //gain stage 2 0.0dB - 9.0dB, step 3.0
    for i:=int32(0); i<=90; i+=30 {
        if getSquareDistance(sum + result[1],roundedGain) > getSquareDistance(sum + i,roundedGain) {
            result[1] = i
        }
    }
    sum+=result[1]
    //gain stage 4 0.0dB - 2.0dB, step 1.0, precision stage
    for i:=int32(0); i<=20; i+=10 {
        if getSquareDistance(sum + result[2],gain) > getSquareDistance(sum + i,gain) {
            result[2] = i
        }
    }


    //correcting for minimal values 
    result[4] = result[4] + 30
    result[5] = result[5] + 30
    return result
}

func (dev *RtlsdrDevice) SetIFGain(gain int32) (actualGain int32) {
    tunerType := dev.GetTunerType()
    var gains []int32
    switch tunerType {
        case 0: gains = IfGainsE4000(gain)
        default: return 0;
    }
    for i,g := range gains {
        retval := int32(C.rtlinput_set_tuner_if_gain(dev.devPtr, (C.int32_t)(i+1),(C.int32_t)(g)))
        if retval != 0 {
            log.Panic("Set IF gain of device ", dev.devId, ", stage ", i+1 ,"returned ", retval)
        }
        actualGain = actualGain + g
    }
    return actualGain
}



