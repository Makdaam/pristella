#include "rtl-input.h"
#include <string.h>
#include <stdio.h>


uint32_t rtlinput_get_device_count(void) {
    return rtlsdr_get_device_count();
}

int32_t rtlinput_get_name(uint32_t dev_id, uint32_t* name_length, uint8_t * name)
{
    uint32_t i;
    const char * name_from_rtlsdr;
    name_from_rtlsdr = rtlsdr_get_device_name(dev_id);

    i = 0;
    while (name_from_rtlsdr[i] && i<100) {
        name[i] = name_from_rtlsdr[i];
        i++;
    }
    (*name_length) = i;
    return 0;
}

int32_t rtlinput_open(void ** dev_ptr, uint32_t dev_id) {
    return rtlsdr_open((rtlsdr_dev_t **) dev_ptr, dev_id);
}

int32_t rtlinput_close(void * dev_ptr) {
    return rtlsdr_close((rtlsdr_dev_t *) dev_ptr);
}

int32_t rtlinput_get_serial(void * dev_ptr, uint8_t * serial) {
    return rtlsdr_get_usb_strings((rtlsdr_dev_t *)dev_ptr, NULL, NULL, (char *)serial);
}

int32_t rtlinput_get_index_from_serial(uint8_t * serial) {
    return  rtlsdr_get_index_by_serial((char *)serial);
}

//frequency stuff
int32_t rtlinput_set_center_freq(void * dev_ptr, uint32_t freq) {
    return rtlsdr_set_center_freq((rtlsdr_dev_t *) dev_ptr, freq);
}

uint32_t rtlinput_get_center_freq(void * dev_ptr) {
    return rtlsdr_get_center_freq((rtlsdr_dev_t *) dev_ptr);
}

int32_t rtlinput_set_sample_rate(void * dev_ptr, uint32_t s_rate) {
    return rtlsdr_set_sample_rate((rtlsdr_dev_t *) dev_ptr, s_rate);
}

uint32_t rtlinput_get_sample_rate(void * dev_ptr) {
    return rtlsdr_get_sample_rate((rtlsdr_dev_t *) dev_ptr);
}

int32_t rtlinput_set_freq_correction(void * dev_ptr, int32_t ppm) {
    return rtlsdr_set_freq_correction((rtlsdr_dev_t *) dev_ptr, ppm);
}

int32_t rtlinput_get_freq_correction(void * dev_ptr) {
    return rtlsdr_get_freq_correction((rtlsdr_dev_t *) dev_ptr);
}

int32_t rtlinput_get_tuner_type(void * dev_ptr) {
    return rtlsdr_get_tuner_type((rtlsdr_dev_t *) dev_ptr);
}

int32_t rtlinput_set_xtal_freq(void * dev_ptr, uint32_t rtl_freq,
				    uint32_t tuner_freq) {
    return rtlsdr_set_xtal_freq((rtlsdr_dev_t *) dev_ptr, rtl_freq, tuner_freq);
}

int32_t rtlinput_get_xtal_freq(void * dev_ptr, uint32_t *rtl_freq,
				    uint32_t *tuner_freq) {
    return rtlsdr_get_xtal_freq((rtlsdr_dev_t *) dev_ptr, rtl_freq, tuner_freq);
}

//gainorama
int32_t rtlinput_get_tuner_gains(void * dev_ptr, int32_t *gains, int32_t *count) {
    *count = rtlsdr_get_tuner_gains((rtlsdr_dev_t *) dev_ptr, NULL);
    if (*count > 255) {
        //go code doesn't expect so many gain values
        return -255;
    }
    return rtlsdr_get_tuner_gains((rtlsdr_dev_t *) dev_ptr, gains);
}

int32_t rtlinput_set_tuner_gain(void * dev_ptr, int32_t gain) {
    return rtlsdr_set_tuner_gain((rtlsdr_dev_t *) dev_ptr, gain);
}

int32_t rtlinput_get_tuner_gain(void * dev_ptr) {
    return rtlsdr_get_tuner_gain((rtlsdr_dev_t *) dev_ptr);
}

int32_t rtlinput_set_tuner_if_gain(void * dev_ptr, int32_t stage, int32_t gain) {
    return rtlsdr_set_tuner_if_gain((rtlsdr_dev_t *) dev_ptr, stage, gain);
}

int32_t rtlinput_set_tuner_gain_mode(void * dev_ptr, int32_t manual) {
    return rtlsdr_set_tuner_gain_mode((rtlsdr_dev_t *) dev_ptr, manual);
}

int32_t rtlinput_set_agc_mode(void * dev_ptr, uint32_t agc_on) {
    return rtlsdr_set_agc_mode((rtlsdr_dev_t *) dev_ptr, agc_on);
}

int32_t rtlinput_set_direct_sampling(void * dev_ptr, int32_t direct_on) {
    return rtlsdr_set_direct_sampling((rtlsdr_dev_t *) dev_ptr, direct_on);
}

int32_t rtlinput_get_direct_sampling(void * dev_ptr) {
    return rtlsdr_get_direct_sampling((rtlsdr_dev_t *) dev_ptr);
}

int32_t rtlinput_reset_buffer(void * dev_ptr) {
    return rtlsdr_reset_buffer((rtlsdr_dev_t *) dev_ptr);
}
int32_t rtlinput_read_sync(void * dev_ptr, void *buf, int len, int *n_read) {
    return rtlsdr_read_sync((rtlsdr_dev_t *) dev_ptr, buf, len, n_read);
}
