package radix

import (
	"errors"
	"fmt"
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

// ErrRadixOutOfRange indicates that the specified radix is not between MinRadix and MaxRadix inclusively.
var ErrRadixOutOfRange = errors.New("radix out of range")

// ToString converts n to a radix string.
//
// The specified radix must be between 2 and 36 (inclusively.)
// If not, ErrRadixOutOfRange will be the returned error.
//
// The returned error is non-nil if an error occurs.
func ToString(n float64, radix int) (string, error) {
	// Radix should meet the requirement: n >= 2 && n <= 36
	if radix < MinRadix || radix > MaxRadix {
		return "", ErrRadixOutOfRange
	}

	// If n is NaN then the result is always NaN.
	if math.IsNaN(n) {
		return "NaN", nil
	}

	// If n is 0, then the result is always zero.
	if n == 0 {
		return "0", nil
	}

	// If n is positive or negative infinity, return infinity in string representation.
	if math.IsInf(n, 1) {
		return "Infinity", nil
	} else if math.IsInf(n, -1) {
		return "-Infinity", nil
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
	if delta <= 0 {
		return "", fmt.Errorf("delta not > 0: %f", delta)
	}
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
							if buffer[fractionCursor] != '.' {
								return "", fmt.Errorf("rune at index %d should be '.', instead: %d (%s)", fractionCursor, buffer[fractionCursor], string(buffer[fractionCursor]))
							}
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
	if fractionCursor >= bufferSize {
		return "", fmt.Errorf("fraction cursor not < buffer size: %d, %d", fractionCursor, bufferSize)
	}
	if intCursor <= 0 {
		return "", fmt.Errorf("int cursor is > 0: %d", intCursor)
	}

	// The result is the range between intCursor and fractionCursor.
	return string(buffer[intCursor:fractionCursor]), nil
}
