/*
 * MIT License
 *
 * Copyright (c) 2022 Crimson Technologies LLC. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package internal

import "math"

const (
	signMask                uint64 = 9223372036854775808
	exponentMask            uint64 = 9218868437227405312
	significandMask         uint64 = 4503599627370495    // excludes hidden bit.
	hiddenBit                      = significandMask + 1 // with hidden bit.
	physicalSignificandSize int    = 52
	exponentBias                   = 0x3FF + physicalSignificandSize
	denormalExponent               = -exponentBias + 1
	infinity                uint64 = 9218868437227405312
)

type F64Wrapper uint64

// WrapF64 wraps n into a new F64Wrapper.
func WrapF64(n float64) F64Wrapper {
	return F64Wrapper(math.Float64bits(n))
}

func (wrapper F64Wrapper) Sign() int {
	if uint64(wrapper)&signMask == 0 {
		return 1
	} else {
		return -1
	}
}

func (wrapper F64Wrapper) Significand() uint64 {
	significand := uint64(wrapper) & significandMask
	if !wrapper.IsDenormal() {
		return significand + hiddenBit
	} else {
		return significand
	}
}

// Next gets the next greatest float64.
// Returns positive infinity on positive infinity input.
func (wrapper F64Wrapper) Next() float64 {
	if uint64(wrapper) == infinity {
		return math.Float64frombits(infinity)
	}

	if wrapper.Sign() < 0 && wrapper.Significand() == 0 {
		return 0
	}

	if wrapper.Sign() < 0 {
		return math.Float64frombits(uint64(wrapper) - 1)
	} else {
		return math.Float64frombits(uint64(wrapper) + 1)
	}
}

func (wrapper F64Wrapper) IsDenormal() bool {
	return uint64(wrapper)&exponentMask == 0
}

func (wrapper F64Wrapper) Exponent() int {
	if wrapper.IsDenormal() {
		return denormalExponent
	}

	biased := int((uint64(wrapper) & exponentMask) >> physicalSignificandSize)

	return biased - exponentBias
}
