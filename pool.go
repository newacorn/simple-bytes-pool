package bpool

import (
	"bufio"
	"errors"
	"github.com/newacorn/goutils/unsafefn"
	"io"
	"sync"
	"unicode/utf8"
)

const (
	Block1k = 1 << 10
	Block2k = 1 << 11
	Block4k = 1 << 12
	Block8k = 1 << 13
)

type BytesPool struct {
	pools [_NumSizeClasses]sync.Pool
}

type Bytes struct {
	B []byte
}

func (b *Bytes) Bytes() []byte {
	return b.B
}

func (b *Bytes) Len() int {
	return len(b.B)
}
func (b *Bytes) Cap() int {
	return cap(b.B)
}
func (b *Bytes) Release() {
	defaultPool.Put(b)
}

// ReadFrom The function appends all the data read from r to b.

// WriteTo implements io.WriterTo.
func (b *Bytes) WriteTo(w io.Writer) (n int64, err error) {
	if nBytes := len(b.B); nBytes > 0 {
		m, e := w.Write(b.B)
		if m > nBytes {
			panic("bpool.Bytes.WriteTo: invalid Write count")
		}
		n = int64(m)
		if e != nil {
			err = e
			return
		}
		if m != nBytes {
			err = io.ErrShortWrite
			return
		}
	}
	return
}

func (b *Bytes) WriteRune(r rune) (n int, err error) {
	// Compare as uint32 to correctly handle negative runes.
	if uint32(r) < utf8.RuneSelf {
		//goland:noinspection GoUnhandledErrorResult
		b.WriteByte(byte(r))
		return 1, nil
	}
	b.Grow(utf8.UTFMax)
	oldL := len(b.B)
	b.B = utf8.AppendRune(b.B, r)
	n = len(b.B) - oldL
	return
}

// Write implements io.Writer - it appends p to ByteBuffer.B
func (b *Bytes) Write(p []byte) (n int, err error) {
	n = len(p)
	if cap(b.B)-len(b.B) >= len(p) {
		b.B = append(b.B, p...)
		return
	}
	b.slowWrite(p)
	return
}

func (b *Bytes) AvailableBuffer() []byte {
	return b.B[len(b.B):cap(b.B)]
}

func (b *Bytes) Available() int {
	return cap(b.B) - len(b.B)
}

func (b *Bytes) slowWrite(p []byte) {
	b2 := Get(len(p) + len(b.B) + len(p)>>1)
	b2.B = b2.B[:len(b.B)+len(p)]
	copy(b2.B, b.B)
	copy(b2.B[len(b.B):], p)
	b.B, b2.B = b2.B, b.B
	Put(b2)
	return
}

func (b *Bytes) slowWriteStr(p string) {
	b2 := Get(len(p) + len(b.B) + len(p)>>1)
	b2.B = b2.B[:len(b.B)+len(p)]
	copy(b2.B, b.B)
	copy(b2.B[len(b.B):], p)
	b.B, b2.B = b2.B, b.B
	Put(b2)
	return
}

// WriteByte appends the byte c to the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
//
// The function always returns nil.
func (b *Bytes) WriteByte(c byte) error {
	b.Grow(1)
	b.B = b.B[:len(b.B)+1]
	b.B[len(b.B)-1] = c
	return nil
}

func (b *Bytes) Swap(new []byte) (old []byte) {
	b.B, old = new, b.B
	return
}

// WriteString appends s to ByteBuffer.B.
func (b *Bytes) WriteString(s string) {
	if cap(b.B)-len(b.B) >= len(s) {
		b.B = append(b.B, s...)
		return
	}
	b.slowWriteStr(s)
}

// Set sets ByteBuffer.B to p.
func (b *Bytes) Set(p []byte) {
	b.B = b.B[:0]
	//goland:noinspection GoUnhandledErrorResult
	b.Write(p)
}

// SetString sets ByteBuffer.B to s.
func (b *Bytes) SetString(s string) {
	b.Set(unsafefn.S2B(s))
}

// String returns string representation of ByteBuffer.B.
func (b *Bytes) String() string {
	return string(b.B)
}

func (b *Bytes) UnsafeString() string {
	return unsafefn.B2S(b.B)
}

// Reset makes ByteBuffer.B empty.
func (b *Bytes) Reset() {
	b.B = b.B[:0]
}

// MinRead is the minimum slice size passed to a Read call by
// [Buffer.ReadFrom]. As long as the [Buffer] has at least MinRead bytes beyond
// what is required to hold the contents of r, ReadFrom will not grow the
// underlying buffer.
const MinRead = 512

func (b *Bytes) Grow(n int) {
	if cap(b.B)-len(b.B) < n {
		b2 := Get(len(b.B) + n)
		b2.B = append(b2.B, b.B...)
		b2.B, b.B = b.B, b2.B
		Put(b2)
	}
}

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with [ErrTooLarge].
func (b *Bytes) ReadFrom(r io.Reader) (n int64, err error) {
	for {
		b.Grow(MinRead)
		m, e := r.Read(b.B[len(b.B):cap(b.B)])
		if m < 0 {
			panic(errNegativeRead)
		}
		b.B = b.B[:len(b.B)+m]
		n += int64(m)
		if e == io.EOF {
			return n, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return n, e
		}
	}
}

//var defaultBufioReaderPool = NewBufioReaderPool(Block4k)
//var defaultBufioWriterPool = NewBufioWriterPool(Block4k)

type BufioReaderPool struct {
	*sync.Pool
	bufioSize int
}
type BufioWriterPool struct {
	*sync.Pool
	bufioSize int
}

const maxBufioSize = 1 << maxBufioSizePower
const minBufioSize = 1 << minBufioSizePower
const minBufioSizePower = 8
const maxBufioSizePower = 23

var defaultBufioReaderPools [maxBufioSizePower - minBufioSizePower + 1]sync.Pool
var defaultBufioWriterPools [maxBufioSizePower - minBufioSizePower + 1]sync.Pool

var defaultBufioReaderTypedPools = newBufioReaderTypedPools()
var defaultBufioWriterTypedPools = newBufioWriterTypedPools()

func newBufioReaderTypedPools() (brp [maxBufioSizePower - minBufioSizePower + 1]BufioReaderPool) {
	for i := minBufioSizePower; i <= maxBufioSizePower; i++ {
		brp[i-minBufioSizePower] = BufioReaderPool{
			Pool:      &defaultBufioReaderPools[i-minBufioSizePower],
			bufioSize: 1 << i,
		}
	}
	return
}

func newBufioWriterTypedPools() (bwp [maxBufioSizePower - minBufioSizePower + 1]BufioWriterPool) {
	for i := minBufioSizePower; i <= maxBufioSizePower; i++ {
		bwp[i-minBufioSizePower] = BufioWriterPool{
			Pool:      &defaultBufioWriterPools[i-minBufioSizePower],
			bufioSize: 1 << i,
		}
	}
	return
}
func GetBw(size int) (bw *bufio.Writer) {
	bit := bsr(size)
	if bit < minBufioSizePower || bit > maxBufioSizePower {
		return bufio.NewWriterSize(nil, size)
	}
	v := defaultBufioWriterPools[bit-minBufioSizePower].Get()
	if v == nil {
		bw = bufio.NewWriterSize(nil, 1<<bit)
	} else {
		bw = v.(*bufio.Writer)
	}
	return
}
func PutBw(bw *bufio.Writer) {
	if bw == nil {
		return
	}
	if !isPowerOfTwo(bw.Size()) {
		return
	}
	bit := bsr(bw.Size())
	if bit < minBufioSizePower || bit > maxBufioSizePower {
		return
	}
	bw.Reset(nil)
	defaultBufioWriterPools[bit-minBufioSizePower].Put(bw)
	return
}
func GetBr(size int) (br *bufio.Reader) {
	bit := bsr(size)
	if bit < minBufioSizePower || bit > maxBufioSizePower {
		return bufio.NewReaderSize(nil, size)
	}
	v := defaultBufioReaderPools[bit-minBufioSizePower].Get()
	if v == nil {
		br = bufio.NewReaderSize(nil, 1<<bit)
	} else {
		br = v.(*bufio.Reader)
	}
	return
}
func PutBr(br *bufio.Reader) {
	if br == nil {
		return
	}
	if !isPowerOfTwo(br.Size()) {
		return
	}
	bit := bsr(br.Size())
	if bit < minBufioSizePower || bit > maxBufioSizePower {
		return
	}
	br.Reset(nil)
	defaultBufioReaderPools[bit-minBufioSizePower].Put(br)
	return
}
func GetBrPool(size int) (brPool *sync.Pool) {
	bit := bsr(size)
	if bit < minBufioSizePower || bit > maxBufioSizePower {
		return &sync.Pool{}
	}
	return &defaultBufioReaderPools[bit-minBufioSizePower]
}

func GetBwPool(size int) (bwPool *sync.Pool) {
	bit := bsr(size)
	if bit < minBufioSizePower || bit > maxBufioSizePower {
		return &sync.Pool{}
	}
	return &defaultBufioWriterPools[bit-minBufioSizePower]
}

func GetBrTypedPool(size int) (brPool BufioReaderPool) {
	bit := bsr(size)
	if bit < minBufioSizePower || bit > maxBufioSizePower {
		return BufioReaderPool{
			Pool:      &sync.Pool{},
			bufioSize: size,
		}
	}
	return defaultBufioReaderTypedPools[bit-minBufioSizePower]
}

func GetBwTypedPool(size int) (bwPool BufioWriterPool) {
	bit := bsr(size)
	if bit < minBufioSizePower || bit > maxBufioSizePower {
		return BufioWriterPool{Pool: &sync.Pool{}, bufioSize: size}
	}
	return defaultBufioWriterTypedPools[bit-minBufioSizePower]
}

func (brP BufioReaderPool) Get() (br *bufio.Reader) {
	br, ok := brP.Pool.Get().(*bufio.Reader)
	if !ok {
		br = bufio.NewReaderSize(nil, brP.bufioSize)
	}
	return
}
func (brP BufioReaderPool) Put(br *bufio.Reader) {
	brP.Pool.Put(br)
}

func (bwP BufioWriterPool) Get() (bw *bufio.Writer) {
	bw, ok := bwP.Pool.Get().(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriterSize(nil, bwP.bufioSize)
	}
	return
}
func (bwP BufioWriterPool) Put(bw *bufio.Writer) {
	bwP.Pool.Put(bw)
}

func (b *Bytes) RecycleToPool00() {
	defaultPool.Put(b)
}

var defaultPool = New()

func Get(size int) *Bytes {
	return defaultPool.Get(size)
}
func Put(bytes *Bytes) {
	defaultPool.Put(bytes)
}

func New() (bp *BytesPool) {
	bp = &BytesPool{}
	for i := 0; i < _NumSizeClasses; i++ {
		bp.pools[i].New = func() interface{} {
			return &Bytes{B: unsafefn.Bytes(0, int(class_to_size[i]))}
		}
	}
	return
}
func (bp *BytesPool) Get(size int) *Bytes {
	if size == 0 {
		return &Bytes{}
	}
	if size <= _MaxBigSize {
		return bp.pools[size2class(size)].Get().(*Bytes)
	}
	return &Bytes{B: unsafefn.Bytes(0, size)}
}
func (bp *BytesPool) Put(bytes *Bytes) {
	if bytes == nil || cap(bytes.B) > _MaxBigSize || cap(bytes.B) < _MinByteSize {
		return
	}
	class := size2class(cap(bytes.B))
	floorSize := class_to_size[class]
	if cap(bytes.B) < int(floorSize) {
		if cap(bytes.B) <= _MaxSmallSize {
			// class cant less  zero
			// because  cap(bytes.B) >= _MinByteSize
			class = class - 1
		} else {
			return
		}
	}
	bytes.B = bytes.B[:0]
	bp.pools[class].Put(bytes)
}

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	pb := Get(Block4k)
	buf := pb.B[:Block4k]
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	Put(pb)
	return written, err
}

// errInvalidWrite means that a write returned an impossible count.
var errInvalidWrite = errors.New("invalid write result")
var errNegativeRead = errors.New("bytes.Buffer: reader returned negative count from Read")
