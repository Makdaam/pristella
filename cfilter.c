// 


/* 
            var j uint32
//        for j=0; j< lc; j++ {
            index1 := j+local_frame
            index2 := (index1)%lc
//            sum += coeffs[j]*state[index2]
            re := real(state[index2]) * coeffs[j]
            im := imag(state[index2]) * coeffs[j]
            sum += complex(re,im)
        }

*/

#include <stdint.h>
#include <math.h>
#include <complex.h>

float complex convolute(float_t * coeffs, float complex * state, uint32_t * lc, uint32_t * local_frame) {
    int32_t j;
    float complex retval;
    retval = 0.0;
    for(j=0;j<*lc;j++) {
        retval = retval + (coeffs[j] * state[((j+*local_frame)%*lc)]);
    }
    return retval;
}
