#include <stdint.h>
#include <math.h>
#include <complex.h>

extern float complex convolute(float_t * coeffs, float complex * state, uint32_t * lc, uint32_t * local_frame);

float complex * filter(float_t * coeffs, float complex * state, const uint32_t lc, float complex * input, const uint32_t li, uint32_t * local_frame);

//nicked directly from Volklib to make sure that the errors aren't an SSE issue
static inline void volk_32fc_32f_dot_prod_32fc_a_generic(lv_32fc_t* result, const lv_32fc_t* input, const float * taps, unsigned int num_points) {

  float *realpt = &((float *)result)[0], *imagpt = &((float *)result)[1];
  const float* aPtr = (float*)input;
  const float* bPtr=  taps;
  unsigned int number = 0;

  *realpt = 0;
  *imagpt = 0;

  for(number = 0; number < num_points; number++){
    *realpt += ((*aPtr++) * (*bPtr));
    *imagpt += ((*aPtr++) * (*bPtr++));
  }

}
