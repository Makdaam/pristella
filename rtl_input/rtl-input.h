#ifndef __RTL_INPUT_H
#define __RTL_INPUT_H

#include <rtl-sdr.h>

uint32_t rtlinput_get_device_count(void);
int32_t rtlinput_get_name(uint32_t dev_id, uint32_t* name_length, uint8_t * name);
int32_t rtlinput_open(void ** dev_ptr, uint32_t dev_id);
int32_t rtlinput_close(void * dev_ptr);

int32_t rtlinput_get_serial(void * dev_ptr, uint8_t * serial);
int32_t rtlinput_get_index_from_serial(uint8_t * serial);

//frequency stuff
int32_t rtlinput_set_center_freq(void * dev_ptr, uint32_t freq);
uint32_t rtlinput_get_center_freq(void * dev_ptr);

int32_t rtlinput_set_freq_correction(void * dev_ptr, int32_t ppm);
int32_t rtlinput_get_freq_correction(void * dev_ptr);

/*
enum rtlinput_tuner {
	RTLSDR_TUNER_UNKNOWN = 0,
	RTLSDR_TUNER_E4000,
	RTLSDR_TUNER_FC0012,
	RTLSDR_TUNER_FC0013,
	RTLSDR_TUNER_FC2580,
	RTLSDR_TUNER_R820T,
	RTLSDR_TUNER_R828D
};

retval from gettuner
*/
int32_t rtlinput_get_tuner_type(void * dev_ptr);


int32_t rtlinput_set_sample_rate(void * dev_ptr, uint32_t freq);
uint32_t rtlinput_get_sample_rate(void * dev_ptr);

int32_t rtlinput_set_xtal_freq(void * dev_ptr, uint32_t rtl_freq,
				    uint32_t tuner_freq); //not used
int32_t rtlinput_get_xtal_freq(void * dev_ptr, uint32_t *rtl_freq,
				    uint32_t *tuner_freq); //not used

//gainorama
int32_t rtlinput_get_tuner_gains(void * dev_ptr, int32_t *gains, int32_t *count);
int32_t rtlinput_set_tuner_gain(void * dev_ptr, int32_t gain);
int32_t rtlinput_get_tuner_gain(void * dev_ptr);
int32_t rtlinput_set_tuner_if_gain(void * dev_ptr, int32_t stage, int32_t gain);
int32_t rtlinput_set_tuner_gain_mode(void * dev_ptr, int32_t if_manual);
int32_t rtlinput_set_agc_mode(void * dev_ptr, uint32_t agc_on);
int32_t rtlinput_set_direct_sampling(void * dev_ptr, int32_t direct_on);
int32_t rtlinput_get_direct_sampling(void * dev_ptr);


/* streaming functions */
int32_t rtlinput_reset_buffer(void * dev_ptr); //
int32_t rtlinput_read_sync(void * dev_ptr, void *buf, int len, int *n_read); //

#endif
