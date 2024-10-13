//go:build ignore

package main

import (
	"bytes"
	"log"
	"math/bits"
	"os"
	"strconv"
)

var class_to_size = [...]uint32{0, 32, 64, 96, 128, 160, 224, 256, 320, 384, 448, 512, 640, 768, 1024, 1280, 1536, 1792, 2048, 2304, 2688, 3072, 3456, 4096, 4864, 5376, 6144, 6784, 8192, 9472, 10240, 10880, 12288, 13568, 14336, 16384, 18432, 19072, 20480, 21760, 24576, 27264, 28672, 32768,
	1 << (15 + 1), 1 << (15 + 2), 1 << (15 + 3), 1 << (15 + 4), 1 << (15 + 5), 1 << (15 + 6), 1 << (15 + 7), 1 << (15 + 8),
}

const (
	_MaxBigSizePower   = 23
	_MaxSamllSizePower = 15
	_MaxSmallSize      = 1 << _MaxSamllSizePower
	_MaxBigSize        = 1 << _MaxBigSizePower
	smallSizeDiv       = 8
	smallSizeMax       = 1024
	largeSizeDiv       = 128
)

func main() {
	output := &bytes.Buffer{}
	_ = output
	size_to_class8_map := make([]int, int(divRoundUp(smallSizeMax-8, smallSizeDiv))+1)
	size_to_class128_map := make([]int, (_MaxSmallSize-smallSizeMax)/largeSizeDiv+1)
	size_to_class_big_map := make([]int, _MaxBigSizePower-_MaxSamllSizePower)
	i := 0
	var idx int
	for j := 0; j <= _MaxBigSize; j++ {
		if j <= int(class_to_size[i]) {
			if j <= smallSizeMax-8 {
				idx := int(divRoundUp(uintptr(j), smallSizeDiv))
				size_to_class8_map[idx] = i
				continue
			}
			if j <= _MaxSmallSize {
				idx = int(divRoundUp(uintptr(j)-smallSizeMax, largeSizeDiv))
				size_to_class128_map[idx] = i
				continue
			}
			size_to_class_big_map[bsr(j)-_MaxSamllSizePower-1] = i
			continue
		}
		j--
		i++
	}
	output.WriteString("package pool\n\n")
	output.WriteString("const (\n")
	//
	constStr := `_MaxBigSizePower   = 23
	_MaxSamllSizePower = 15
	_MaxSmallSize      = 1 << _MaxSamllSizePower
	_MaxBigSize        = 1 << _MaxBigSizePower
	smallSizeDiv       = 8
	smallSizeMax       = 1024
	largeSizeDiv       = 128`
	//
	output.WriteString(constStr + "\n")
	output.WriteString("\t_NumSizeClasses = " + strconv.Itoa(len(class_to_size)) + "\n")
	output.WriteString(")\n\n")

	output.WriteString("//goland:noinspection GoSnakeCaseUsage\n")
	output.WriteString("var class_to_size = [...]uint32")
	output.WriteByte('{')
	for _, v := range class_to_size {
		output.WriteString(strconv.Itoa(int(v)) + ", ")
	}
	output.WriteString("}\n")

	output.WriteString("//goland:noinspection GoSnakeCaseUsage\n")
	output.WriteString("var size_to_class8 = [smallSizeMax/smallSizeDiv]uint8")
	output.WriteByte('{')
	for _, v := range size_to_class8_map {
		output.WriteString(strconv.Itoa(int(v)) + ", ")
	}
	output.WriteString("}\n")

	output.WriteString("//goland:noinspection GoSnakeCaseUsage\n")
	output.WriteString("var size_to_class128 = [(_MaxSmallSize-smallSizeMax)/largeSizeDiv + 1]uint8")
	output.WriteByte('{')
	for _, v := range size_to_class128_map {
		output.WriteString(strconv.Itoa(int(v)) + ", ")
	}
	output.WriteString("}\n")

	output.WriteString("//goland:noinspection GoSnakeCaseUsage\n")
	output.WriteString("var size_to_class_big = [_MaxBigSizePower-_MaxSamllSizePower]uint8")
	output.WriteByte('{')
	for _, v := range size_to_class_big_map {
		output.WriteString(strconv.Itoa(int(v)) + ", ")
	}
	output.WriteString("}\n")

	f, err := os.Create("sizeclass_table.go")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = f.Close()
	}()
	output.WriteTo(f)
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
