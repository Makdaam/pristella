package main

//#cgo LDFLAGS: -lvolk -L/usr/local/lib/
//#define LV_HAVE_SSE 1
//#include <volk/volk_typedefs.h>
//#include <volk/volk_32fc_32f_dot_prod_32fc_a.h>
//#include "cfilter.h"
//
//
//void volk_dot(float complex * res, float complex * state, float_t * coeffs, uint32_t lc) {
//  volk_32fc_32f_dot_prod_32fc_a_generic((lv_32fc_t* )res,(lv_32fc_t* )state,(float *)coeffs,(unsigned int)lc);
//}
import "C"
import "log"


var vcoeffs map[uint32][]float32

//FIR filter with multiplication moved to C
func firFilterC(in []complex64, fir_frame *int, coeffs []float32, state []complex64) []complex64 {
	var lc uint32
	lc = uint32(len(coeffs))
	var local_frame uint32
	local_frame = uint32(*fir_frame)
	var li uint32
	li = uint32(len(in))

	C.filter((*C.float_t)(&coeffs[0]), (*C.complexfloat)(&state[0]), (C.uint32_t)(lc), (*C.complexfloat)(&in[0]), (C.uint32_t)(li), (*C.uint32_t)(&local_frame))

	*fir_frame = int(local_frame)
	return in
}

//init function for memory aligned coefficients slices
func initFilterV(coeffs []float32) {
    lc:=uint32(len(coeffs))
    vcoeffs = make(map[uint32][]float32)
    for j:=uint32(0); j<4; j++ {
        vcoeffs[j] = make([]float32,2*lc)
        for i:=uint32(0); i<lc; i++ {
            vcoeffs[j][(lc-j+i)%lc] = coeffs[i]
            vcoeffs[j][lc+((lc-j+i)%lc)] = coeffs[i]
        }
    }
}
    
//FIR filter with vector multiplication with mem aligned coefficients
//there's something wrong here
func firFilterV(in []complex64, fir_frame *int, coeffs []float32, state []complex64) []complex64 {
    aligner := uint32(4)
	out := make([]complex64, len(in))
	var lc uint32
	lc = uint32(len(coeffs))
    if lc % 4 != 0 {
        log.Fatal("FIR tap count has to be divisible by 4")
    }
    if len(in) % 4 != 0 {
        log.Fatal("Sample buffer not divisible by 4")
    }
    //number of samples should be divisible by 4
	var local_frame uint32
	local_frame = uint32(*fir_frame)
	var sum complex64
    long_loop := uint32(len(in))/aligner    
	for i := uint32(0); i < long_loop; i++ {
        for j:= uint32(0); j<aligner; j++ {
		    //place value in state
		    state[local_frame+j] = in[(i*aligner)+j]
		    //calculate output
            C.volk_dot((*C.complexfloat)(&sum), (*C.complexfloat)(&state[0]),(*C.float_t)(&vcoeffs[j][local_frame]), (C.uint32_t)(lc))
		    out[(i*aligner)+j] = sum
		    //shift fir_frame
        }
	    local_frame = (local_frame + aligner) % (lc - aligner)
	}
	*fir_frame = int(local_frame)
	return out
}

//FIR filter, trivial but slow implementation
func firFilter(in []complex64, fir_frame *int, coeffs []float32, state []complex64) []complex64 {
	out := make([]complex64, len(in))
	var lc uint32
	lc = uint32(len(coeffs))
	var local_frame uint32
	local_frame = uint32(*fir_frame)
	var sum complex64
	for i := 0; i < len(in); i++ {
		//place value in state
		state[local_frame] = in[i]
		//calculate output
		sum = 0
		var j uint32
		for j = 0; j < lc; j++ {
			index1 := j + local_frame
			index2 := (index1) % lc
			//            sum += coeffs[j]*state[index2]
			re := real(state[index2]) * coeffs[j]
			im := imag(state[index2]) * coeffs[j]
			sum += complex(re, im)
		}
		out[i] = sum
		//shift fir_frame
		local_frame = (local_frame + 1) % lc
	}
	*fir_frame = int(local_frame)
	return out
}
