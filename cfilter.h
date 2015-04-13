#include <stdint.h>
#include <math.h>
#include <complex.h>

extern float complex convolute(float_t * coeffs, float complex * state, uint32_t * lc, uint32_t * local_frame);

float complex * filter(float_t * coeffs, float complex * state, const uint32_t lc, float complex * input, const uint32_t li, uint32_t * local_frame);
