package bpool

import (
	"github.com/gookit/goutil/testutil/assert"
	"github.com/xyproto/randomstring"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestWrite(t *testing.T) {
	pb := Get(20)
	dst := make([]byte, 0, 20*10)
	for i := 0; i < 20; i++ {
		//goland:noinspection GoUnhandledErrorResult
		pb.WriteString("1234567890")
		dst = append(dst, "1234567890"...)
	}
	assert.Equal(t, len(dst), pb.Len())
	assert.Equal(t, dst, pb.Bytes())
}
func TestBytes_ReadFrom(t *testing.T) {
	s := randomstring.EnglishFrequencyString(4093)
	r := strings.NewReader(s)
	pb := Get(123)
	n, err := pb.ReadFrom(r)
	assert.Equal(t, int64(4093), n)
	assert.Nil(t, err)
	assert.Equal(t, pb.UnsafeString(), s)
	pb.Release()
}

//goland:noinspection GoUnhandledErrorResult
func TestBytes_WriteByte(t *testing.T) {
	pb := Get(123)
	for i := 0; i < 100; i++ {
		pb.WriteByte('a')
		pb.WriteByte('b')
		pb.WriteByte('c')
	}
	assert.Equal(t, 300, pb.Len())
	assert.Equal(t, strings.Repeat("abc", 100), pb.UnsafeString())
}

//goland:noinspection GoUnhandledErrorResult
func TestBytes_WriteRune(t *testing.T) {
	pb := Get(412)
	for i := 0; i < 200; i++ {
		pb.WriteRune('a')
		pb.WriteRune('b')
		pb.WriteRune('c')
		r, _ := utf8.DecodeRune([]byte("我"))
		pb.WriteRune(r)
	}
	assert.Equal(t, strings.Repeat("abc我", 200), pb.UnsafeString())
	pb.Release()
}

func TestBytes_Grow(t *testing.T) {
	pb := Get(123)
	str := randomstring.String(323)
	pb.WriteString(str)
	for i := 0; i < 10; i++ {
		pb.Grow(i * 100)
	}
	assert.Equal(t, pb.UnsafeString(), str)
}
