package internal

const(
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

// Wraps a float64 into a new wrapper.
func WrapF64(n float64) F64Wrapper {
	return F64Wrapper(Float64ToUInt64(n))
}

func (wrapper F64Wrapper) Sign() int {
	if uint64(wrapper) &signMask == 0 {
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

// Returns the next greater float64. Returns positive infinity on positive infinity input.
func (wrapper F64Wrapper) Next() float64 {
	if uint64(wrapper) == infinity {
		return float64(infinity)
	}

	if wrapper.Sign() < 0 && wrapper.Significand() == 0 {
		return 0
	}

	if wrapper.Sign() < 0 {
		return UInt64ToFloat64(uint64(wrapper) - 1)
	} else {
		return UInt64ToFloat64(uint64(wrapper) + 1)
	}
}

// Returns true if the float64 is denormal.
func (wrapper F64Wrapper) IsDenormal() bool {
	return uint64(wrapper) &exponentMask == 0
}

func (wrapper F64Wrapper) Exponent() int {
	if wrapper.IsDenormal() {
		return denormalExponent
	}

	biased := int((uint64(wrapper) & exponentMask) >> physicalSignificandSize)

	return biased - exponentBias
}
