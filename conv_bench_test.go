package radix

import "testing"

// Benchmarks positive conversion.
func BenchmarkConv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := ToString(1.2567, 16); err != nil {
			b.Error(err)
		}
	}
}
