// Contains functions to convert between float64 and uint64 without changing any bits.
package internal

import "unsafe"

// Casts a 64-bit float to an unsigned 64-bit integer.
func Float64ToUInt64(n float64) uint64 {
	return *(*uint64) (unsafe.Pointer(&n))
}

// Casts an unsigned 64-bit integer to a 64-bit float.
func UInt64ToFloat64(n uint64) float64 {
	return *(*float64) (unsafe.Pointer(&n))
}
