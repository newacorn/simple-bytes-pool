package bpool

import "math/bits"

func size2class(size int) (sizeclass uint8) {
	if size <= smallSizeMax-8 {
		sizeclass = size_to_class8[divRoundUp(uintptr(size), smallSizeDiv)]
	} else if size <= _MaxSmallSize {
		sizeclass = size_to_class128[divRoundUp(uintptr(size)-smallSizeMax, largeSizeDiv)]
	} else {
		sizeclass = size_to_class_big[bsr(size)-_MaxSamllSizePower-1]
	}
	return
}

// divRoundUp returns ceil(n / a).
func divRoundUp(n, a uintptr) uintptr {
	// a is generally a power of two. This will get inlined and
	// the compiler will optimize the division.
	return (n + a - 1) / a
}

func bsr(x int) (bitsSize int) {
	bitsSize = bits.Len(uint(x))
	if isPowerOfTwo(x) {
		bitsSize -= 1
	}
	return
}

func isPowerOfTwo(x int) bool {
	return (x & (-x)) == x
}
