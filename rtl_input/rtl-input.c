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

int32_t rtlinput_set_center_freq(void * dev_ptr, uint32_t freq) {
    return rtlsdr_set_center_freq((rtlsdr_dev_t *) dev_ptr, freq);
}

int32_t rtlinput_set_sample_rate(void * dev_ptr, uint32_t s_rate) {
    return rtlsdr_set_sample_rate((rtlsdr_dev_t *) dev_ptr, s_rate);
}
