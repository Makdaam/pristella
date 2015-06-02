package main

import (
	"fmt"
	"log"
	"math"
)

//FIR coefficients for different filters
//all fixed point operations made in Q16.15 format
//so 32768 (1 << 15) represents 1

const twopi = 205887
const One = 32768

var intsin []int32
var intcos []int32


var low_divider uint32
var high_divider uint32
var average float64
var target_average float64
var decim_counter uint32
var step uint32

func init_decimate(in_sample_rate uint32, out_sample_rate uint32) {
	target_average = float64(in_sample_rate) / float64(out_sample_rate)
	low_divider = uint32(math.Floor(target_average))
	high_divider = uint32(math.Ceil(target_average))
	average = target_average
	step = high_divider
}

func decimate(in []complex64) []complex64 {
	out := make([]complex64, (len(in) / int(math.Floor(target_average))))
	j := uint32(0)
	for i := uint32(0); i < uint32(len(in)); i += step {
		//        fmt.Println("ST",step, average, target_average)
		out[j] = in[i]
		j += 1
		average = 0.99*average + 0.01*float64(step)
		if average > target_average {
			step = low_divider
		} else {
			step = high_divider
		}
	}
	return out[0:j]
}

type FreqShifter struct {
	exp  int32
	frame int32
}

func newFreqShifter(freqOffset int, sampleRate int) FreqShifter {
    var out FreqShifter
    fmt.Println("####",freqOffset)
    out.exp = int32(32768.0 * (2*math.Pi) * float64(-freqOffset) / float64(sampleRate)) % twopi
    for out.exp<0 {
        out.exp += twopi
    }
    out.frame = 0
    return out
}

func (fs *FreqShifter) freqShift(in []int32) []int32 {
	out := make([]int32, len(in))

	for i:=int32(0); i < int32(len(in)); i+=2 {
        framepos := fs.frame+(i/2) % twopi
        trigarg := int32((int64(framepos) * int64(fs.exp)) % twopi)
        out[i] = ((intcos[trigarg] * in[i]) /32768) - ((intsin[trigarg] * in[i+1]) /32768)
        out[i+1] = ((intcos[trigarg] * in[i+1]) /32768) + ((intsin[trigarg] * in[i]) /32768)
	}
    fs.frame = (fs.frame + int32(len(in)/2)) % twopi
	return out
}

func toComplex(in []int32) []complex64 {
    out := make([]complex64,len(in)/2,len(in)/2)
    for i := 0; i<len(in); i+=2 {
        out[i/2] = complex(float32(in[i])/32768.0,float32(in[i+1])/32768.0)
    }
    return out   
}

func initSinCos() {
    intsin = make([]int32,twopi,twopi)
    intcos = make([]int32,twopi,twopi)
    for i:= 0; i<twopi; i++ {
        intsin[i]=int32(math.Sin(float64(i)/32768)*32768)
        intcos[i]=int32(math.Cos(float64(i)/32768)*32768)
    }
}

func dsp(in chan []int32, out chan []complex64) {
	var current_coeffs []float32
	var current_fir_state []complex64
	var sample_rate_ok bool

	init_coeffs()

    fmt.Println("Initializing sin cos")
    initSinCos()
    fmt.Println("Sin cos done")
	current_coeffs, sample_rate_ok = fir_coeffs[sample_rate]
	if !sample_rate_ok {
		log.Fatal("Missing FIR coefficients for this sample_rate, see fir_coeffs.go for details")
	}
	current_fir_state = make([]complex64, len(current_coeffs))
	fmt.Println(current_coeffs[0])
	initFilterV(current_coeffs)
	init_decimate(uint32(sample_rate), 36000)
	fmt.Println("FILTER")
	t := 0
	f := 0
    fs := newFreqShifter(offset_freq,sample_rate)
	for {
		buf1 := <-in
		//select the filtering function (from filter.go) over here
		//out <- decimate(firFilterV(freqShift(buf1, &t), &f, current_coeffs, current_fir_state))
		//out <- firFilterV(freqShift(buf1, &t), &f, current_coeffs, current_fir_state)
//        out <- toComplex(buf1)
        out <- toComplex(fs.freqShift(buf1))
		if t < -1000 && f < -1000 {
			fmt.Println("Unexpected", len(current_fir_state))
		}
		//out <- buf1
		//out <- freqShift(buf1, &t)
	}
}
