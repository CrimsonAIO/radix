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

package radix

import (
	"github.com/CrimsonAIO/radix/internal"
	"math"
)

const (
	// characters used for character conversion.
	chars = "0123456789abcdefghijklmnopqrstuvwxyz"

	// buffer array size.
	bufferSize = 2200

	// MinRadix is the minimum radix that can be specified in ToString.
	MinRadix = 2
	// MaxRadix is the maximum radix that can be specified in ToString.
	MaxRadix = 36
)

// ToString converts n to radix.
//
// The specified radix must be between MinRadix and MaxRadix inclusively.
func ToString(n float64, radix int) string {
	// If n is NaN then the result is always NaN.
	if math.IsNaN(n) {
		return "NaN"
	}

	// If n is 0, then the result is always zero.
	if n == 0 {
		return "0"
	}

	// If n is positive or negative infinity, return infinity in string representation.
	if math.IsInf(n, 1) {
		return "Infinity"
	} else if math.IsInf(n, -1) {
		return "-Infinity"
	}

	// Buffer for the result. We start with the decimal point in the
	// middle and write to the left for the integer part and to the right for the
	// fractional part. 1024 characters for the exponent and 52 for the mantissa
	// either way, with additional space for sign, decimal point and string
	// termination should be sufficient.
	var buffer [bufferSize]rune
	intCursor := bufferSize / 2
	fractionCursor := intCursor

	// If n is negative then flip n.
	negative := n < 0
	if negative {
		n = -n
	}

	// Split the value into an integer part and a fractional part.
	integer := math.Floor(n)
	fraction := n - integer
	// We only compute fractional digits up to the input float64's precision.
	delta := 0.5 * (internal.WrapF64(n).Next() - n)
	delta = math.Max(internal.WrapF64(0).Next(), delta)
	if fraction >= delta {
		// Insert decimal point.
		buffer[fractionCursor] = '.'
		fractionCursor++

		for {
			// Shift up by one digit.
			fraction *= float64(radix)
			delta *= float64(radix)

			// Write digit.
			digit := int(fraction)
			buffer[fractionCursor] = rune(chars[digit])
			fractionCursor++

			// Calculate remainder.
			fraction -= float64(digit)

			// Round to even.
			if fraction > 0.5 || (fraction == 0.5 && (digit&1) == 1) {
				if fraction+delta > 1 {
					// We need to back trace already written digits in case of carry-over.
					for {
						fractionCursor--
						if fractionCursor == bufferSize/2 {
							// Carry over the integer part.
							integer += 1
							break
						}

						c := buffer[fractionCursor]
						// Reconstruct digit.
						var digit rune
						if c > '9' {
							digit = c - 'a' + 10
						} else {
							digit = c - '0'
						}
						if int(digit)+1 < radix {
							buffer[fractionCursor] = rune(chars[digit+1])
							fractionCursor++
							break
						}
					}
					break
				}
			}

			// If (fraction >= delta) isn't true then break.
			if fraction < delta {
				break
			}
		}
	}

	// Compute integer digits. Fill unrepresented digits with zero.
	for internal.WrapF64(integer/float64(radix)).Exponent() > 0 {
		integer /= float64(radix)
		intCursor--
		buffer[intCursor] = '0'
	}
	for {
		remainder := math.Mod(integer, float64(radix))
		intCursor--
		buffer[intCursor] = rune(chars[int(remainder)])
		integer = (integer - remainder) / float64(radix)

		// If (integer > 0) isn't true then break.
		if integer <= 0 {
			break
		}
	}

	// Add sign if negative.
	if negative {
		intCursor--
		buffer[intCursor] = '-'
	}

	// The result is the range between intCursor and fractionCursor.
	return string(buffer[intCursor:fractionCursor])
}
