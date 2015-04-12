package main

//#include "cfilter.h"
import "C"


func firFilterC (in []complex64, fir_frame *int, coeffs[]float32, state[]complex64) []complex64{
    out := make([]complex64, len(in))
    var lc uint32
    lc = uint32(len(coeffs))
    var local_frame uint32
    local_frame = uint32(*fir_frame)
    var li uint32
    li = len(in)
    
    for i:=0; i < li; i++ {
        //place value in state
        state[local_frame]=in[i]
        //calculate output
        out[i] = complex64(C.convolute((*C.float_t)(&coeffs[0]),(*C.complexfloat)(&state[0]),(*C.uint32_t)(&lc),(*C.uint32_t)(&local_frame)))
        //shift fir_frame
        local_frame = (local_frame + 1) % lc
    }
    *fir_frame = int(local_frame)
    return out
}

func firFilter (in []complex64, fir_frame *int, coeffs[]float32, state[]complex64) []complex64{
    out := make([]complex64, len(in))
    var lc uint32
    lc = uint32(len(coeffs))
    var local_frame uint32
    local_frame = uint32(*fir_frame)
    var sum complex64
    for i:=0; i < len(in); i++ {
        //place value in state
        state[local_frame]=in[i]
        //calculate output
        sum = 0
        var j uint32
        for j=0; j< lc; j++ {
            index1 := j+local_frame
            index2 := (index1)%lc
//            sum += coeffs[j]*state[index2]
            re := real(state[index2]) * coeffs[j]
            im := imag(state[index2]) * coeffs[j]
            sum += complex(re,im)
        }
        out[i] = sum
        //shift fir_frame
        local_frame = (local_frame + 1) % lc
    }
    *fir_frame = int(local_frame)
    return out
}
