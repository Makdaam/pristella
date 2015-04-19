#ifndef __RTL_INPUT_H
#define __RTL_INPUT_H

#include <rtl-sdr.h>

uint32_t rtlinput_get_device_count(void);
int32_t rtlinput_get_name(uint32_t dev_id, uint32_t* name_length, uint8_t * name);
int32_t rtlinput_open(void ** dev_ptr, uint32_t dev_id);
int32_t rtlinput_close(void * dev_ptr);

int32_t rtlinput_set_center_freq(void * dev_ptr, uint32_t freq);
int32_t rtlinput_set_sample_rate(void * dev_ptr, uint32_t freq);

#endif /* __RTL_INPUT_H */
