package radix

import (
	"testing"
)

// Tests positive conversion.
func TestPosConv(t *testing.T) {
	if result, err := ToString(1.2567, 16); err != nil {
		t.Error(err)
	} else if result != "1.41b71758e2196" {
		t.Logf("Bad result: %s", result)
		t.Fail()
	} else {
		t.Log("Correct result")
	}
}

// Tests negative conversion.
func TestNegConv(t *testing.T) {
	if result, err := ToString(-1.2567, 16); err != nil {
		t.Error(err)
	} else if result != "-1.41b71758e2196" {
		t.Logf("Bad result: %s", result)
		t.Fail()
	} else {
		t.Logf("Correct result")
	}
}
