package mpeg

import (
	"awCodec/utils"
)

// These bits are used in joint_stereo mode
var bits = [15]int{
	0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, // 16 invalid
}

// (2^nb / (2^nb - 1))
var linear = [15]float32{
	1.3333333333333333, 1.1428571428571428, 1.0666666666666667, 1.032258064516129, 1.0158730158730158, 1.0078740157480315, 1.003921568627451,
	1.0019569471624266, 1.0009775171065494, 1.0004885197850513, 1.0002442002442002, 1.0001220852154804, 1.0000610388817677, 1.000030518509476,
}

// 2^(-nb + 1)
var linear2 = [15]float32{
	0.5, 0.25, 0.125, 0.0625, 0.03125, 0.015625, 0.0078125, 0.00390625, 0.001953125, 0.0009765625, 0.00048828125, 0.000244140625, 0.0001220703125, 6.103515625e-05,
}

func requantization(s, nb int) float32 {
	s ^= 1 << (nb - 1)
	s |= -(s & (1 << (nb - 1))) // s'''
	shift := float32(int(1) << uint(nb-1))
	s3 := float32(s) / shift
	// requantize --------------------------------------------------
	// s'' = (2^nb / (2^nb - 1)) * (s''' + 2^(-nb + 1))
	return linear[nb-2] * (s3 + linear2[nb-2]) // s''
}

func synthSubbandFilter(samples []float32, ch int, vVec *[1024]float32, pcm_ []float32) {
	// Shift V vector --------------------------------------------------
	copy(vVec[64:], vVec[0:960]) // 1024-64

	// Matrix --------------------------------------------------
	for i := 0; i < 64; i++ {
		vVec[i] = 0.0
		for k := 0; k < 32; k++ {
			vVec[i] += samples[k] * synthN[i][k]
		}
	}

	// Build the U vector --------------------------------------------------
	uVec := [512]float32{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 32; j++ {
			uVec[i*64+j] = vVec[i*128+j]
			uVec[i*64+32+j] = vVec[i*128+96+j]
		}
	}

	// Build the W vector --------------------------------------------------
	wVec := [512]float32{}
	for i := 0; i < 512; i++ {
		wVec[i] = uVec[i] * synthD[i]
	}

	// Calc 32 samples --------------------------------------------------
	for i := 0; i < 32; i++ {
		var S float32
		for j := 0; j < 512; j += 32 {
			S += wVec[j+i]
		}

		idx := 2 * i
		if ch == 0 {
			pcm_[idx] = S
		} else {
			pcm_[idx+1] = S
		}
	}
}

func decodeLayer1(br *utils.BitReader, nch int, bound int) (int, [2][32 * 12]float32) {
	// allocation --------------------------------------------------
	allocation := [2][32]int{}
	for sb := 0; sb < bound; sb++ {
		for ch := 0; ch < nch; ch++ {
			// TODO add error if value equal 15 (invalid value) (if really needed error)
			allocation[ch][sb] = bits[br.ReadBits(4)]
		}
	}
	for sb := bound; sb < 32; sb++ {
		allocation[0][sb] = bits[br.ReadBits(4)]
		allocation[1][sb] = allocation[0][sb]
	}

	// scalefactor --------------------------------------------------
	scaleFactor := [2][32]float32{}
	for sb := 0; sb < 32; sb++ {
		for ch := 0; ch < nch; ch++ {
			if allocation[ch][sb] != 0 {
				//scaleFactor[ch][sb] = float32(br.ReadBits(6))
				scaleFactor[ch][sb] = requantizeFactor[br.ReadBits(6)]
			}
		}
	}

	// samples --------------------------------------------------
	samples := [2][32 * 12]float32{}
	for s := 0; s < 12; s++ {
		for sb := 0; sb < bound; sb++ {
			for ch := 0; ch < nch; ch++ {
				if allocation[ch][sb] != 0 {
					nb := allocation[ch][sb]
					samples[ch][32*s+sb] = requantization(br.ReadBits(nb), nb) * scaleFactor[ch][sb] // s' = factor * s''
				}
			}
		}
		for sb := bound; sb < 32; sb++ {
			if allocation[0][sb] != 0 {
				nb := allocation[0][sb]
				samples[0][32*s+sb] = requantization(br.ReadBits(nb), nb) * scaleFactor[1][sb] // s' = factor * s''
				samples[1][32*s+sb] = samples[0][32*s+sb]                                      // s' = factor * s''
			}
		}
	}
	return 0, samples
}
