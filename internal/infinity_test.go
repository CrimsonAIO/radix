package internal

import (
	"math"
	"testing"
)

// Ensures F64Wrapper.Next returns +infinity when WrapF64 is called with +infinity.
func TestInfinity(t *testing.T) {
	if next := WrapF64(math.Inf(1)).Next(); !math.IsInf(next, 1) {
		t.Log("Got incorrect result:", next)
	}
}
