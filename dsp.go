package main

import (
	"fmt"
	"log"
	"math"
)

//FIR coefficients for different filters

// sin coefficients
var _sin = [...]float32{
	1.58962301576546568060E-10, // 0x3de5d8fd1fd19ccd
	-2.50507477628578072866E-8, // 0xbe5ae5e5a9291f5d
	2.75573136213857245213E-6,  // 0x3ec71de3567d48a1
	-1.98412698295895385996E-4, // 0xbf2a01a019bfdf03
	8.33333333332211858878E-3,  // 0x3f8111111110f7d0
	-1.66666666666666307295E-1, // 0xbfc5555555555548
}

// cos coefficients
var _cos = [...]float32{
	-1.13585365213876817300E-11, // 0xbda8fa49a0861a9b
	2.08757008419747316778E-9,   // 0x3e21ee9d7b4e3f05
	-2.75573141792967388112E-7,  // 0xbe927e4f7eac4bc6
	2.48015872888517045348E-5,   // 0x3efa01a019c844f5
	-1.38888888888730564116E-3,  // 0xbf56c16c16c14f91
	4.16666666666665929218E-2,   // 0x3fa555555555554b
}

//this one is ripped from math.Sincos, but is faster by about 15%
func sincos32(x float32) (sin, cos float32) {
	const (
		PI4A = 7.85398125648498535156E-1                             // 0x3fe921fb40000000, Pi/4 split into three parts
		PI4B = 3.77489470793079817668E-8                             // 0x3e64442d00000000,
		PI4C = 2.69515142907905952645E-15                            // 0x3ce8469898cc5170,
		M4PI = 1.273239544735162542821171882678754627704620361328125 // 4/pi
	)
	// special cases
	switch {
	case x == 0:
		return x, 1 // return ±0.0, 1.0
	}

	// make argument positive
	sinSign, cosSign := false, false
	if x < 0 {
		x = -x
		sinSign = true
	}

	j := int32(x * M4PI) // integer part of x/(Pi/4), as integer for tests on the phase angle
	y := float32(j)      // integer part of x/(Pi/4), as float

	if j&1 == 1 { // map zeros to origin
		j += 1
		y += 1
	}
	j &= 7     // octant modulo 2Pi radians (360 degrees)
	if j > 3 { // reflect in x axis
		j -= 4
		sinSign, cosSign = !sinSign, !cosSign
	}
	if j > 1 {
		cosSign = !cosSign
	}

	z := ((x - y*PI4A) - y*PI4B) - y*PI4C // Extended precision modular arithmetic
	zz := z * z
	cos = 1.0 - 0.5*zz + zz*zz*((((((_cos[0]*zz)+_cos[1])*zz+_cos[2])*zz+_cos[3])*zz+_cos[4])*zz+_cos[5])
	sin = z + z*zz*((((((_sin[0]*zz)+_sin[1])*zz+_sin[2])*zz+_sin[3])*zz+_sin[4])*zz+_sin[5])
	if j == 1 || j == 2 {
		sin, cos = cos, sin
	}
	if cosSign {
		cos = -cos
	}
	if sinSign {
		sin = -sin
	}
	return
}

func freqShift(in []complex64, frame *int) []complex64 {
	out := make([]complex64, len(in))
	exp := float32(2*math.Pi) * float32(-offset_freq) / float32(sample_rate)
	for key, value := range in {
		floatframe := float32(*frame + key)
		sin, cos := sincos32(exp * floatframe)
		out[key] = value * complex64(complex(cos, sin))
	}
	return out
}

func dsp(in chan []complex64, out chan []complex64) {
	var current_coeffs []float32
	var current_fir_state []complex64;
	var sample_rate_ok bool

	init_coeffs()
	current_coeffs, sample_rate_ok = fir_coeffs[sample_rate]
	if !sample_rate_ok {
		log.Fatal("Missing FIR coefficients for this sample_rate, see fir_coeffs.go for details")
	}
	current_fir_state = make([]complex64, len(current_coeffs))
	fmt.Println(current_coeffs[0])
    initFilterV(current_coeffs)
	fmt.Println("FILTER")
	t := 0
    f := 0
	for {
		buf1 := <-in
        //select the filtering function (from filter.go) over here
		out <- firFilterC(freqShift(buf1, &t), &f, current_coeffs, current_fir_state)
		//out <- freqShift(buf1, &t)
	}
}
