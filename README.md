# radix
Convert floats to radix strings, just the way JavaScript does it.

This does **not** call JavaScript code whatsoever. This is a completely native solution to Go.

## Usage

```go
package main

import (
	"fmt"
	"github.com/CrimsonAIO/radix"
)

func main() {
	// prints: 1.41b71758e2196
	// the same result can be achieved via calling "(1.2567).toString(16)" in JavaScript.
	fmt.Println(radix.ToString(1.2567, 16))
}
```
Check it out on [The Go Playground](https://play.golang.org/p/TiKMW2LnEyH).

## Benchmarks
This benchmark was ran on an iMac with an M1 processor.
```
goos: darwin
goarch: arm64
pkg: github.com/CrimsonAIO/radix
BenchmarkConv
BenchmarkConv-8   	 3071470	       383.9 ns/op
```
If you wish to benchmark on your own system, clone this repository and run `go test -bench .`

## Credits
The logic for this has been taken from Google's V8 engine. The source can be found at: \
[DoubleToRadixCString (main function)](https://github.com/v8/v8/blob/f83601408c3207211bc8eb82a8802b01fd82c775/src/numbers/conversions.cc#L1269) \
[double.h (double-uint64 utility)](https://github.com/v8/v8/blob/f83601408c3207211bc8eb82a8802b01fd82c775/src/numbers/double.h)
