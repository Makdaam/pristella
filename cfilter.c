// 
/*
    for i:=0; i < li; i++ {
        //place value in state
        state[local_frame]=in[i]
        //calculate output
        out[i] = complex64(C.convolute((*C.float_t)(&coeffs[0]),(*C.complexfloat)(&state[0]),(*C.uint32_t)(&lc),(*C.uint32_t)(&local_frame)))
        //shift fir_frame
        local_frame = (local_frame + 1) % lc
    }

*/

#include <stdint.h>
#include <math.h>
#include <complex.h>

float complex convolute(const float_t * coeffs, float complex * state, uint32_t * lc, uint32_t * local_frame) {
    int32_t j;
    float complex retval;
    retval = 0.0;
    for(j=0;j<*lc;j++) {
        retval = retval + (coeffs[j] * state[((j+*local_frame)%*lc)]);
    }
    return retval;
}

//calculates the convolution of a whole input buffer in place
float complex * filter(float_t * coeffs, float complex * state, const uint32_t lc, float complex * input, const uint32_t li, uint32_t * local_frame) {
    int32_t i;
    int32_t j;

    for(i=0; i < li; i++) {
        //for each sample in input
        //place the sample in state
        state[*local_frame]=input[i];
        //calculate the output
        input[i] = 0.0;
        for(j=0;j< lc;j++) {
            input[i] = input[i] + (coeffs[j] * state[((j+*local_frame)%lc)]);
        }
        //shift the local_frame (move one ahead in the state circular buffer)
        *local_frame = (*local_frame + 1) % lc;
    }

    return input;
}
