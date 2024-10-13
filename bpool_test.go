package bpool

import (
	"math/rand"
	"slices"
	"testing"
)

func TestSizeInRange(t *testing.T) {
	defer func() {
		if recover() != nil {
			t.Fatal("panic occur")
		}
	}()
	for i := 0; i <= _MaxBigSize; i++ {
		pb := Get(i)
		Put(pb)
	}
}

func TestCapSizeInList(t *testing.T) {
	for i := 0; i <= _MaxBigSize; i++ {
		pb := Get(i)
		if slices.Index(class_to_size[:], uint32(cap(pb.B))) == -1 {
			t.Fatalf("Get(%d) return bytes's cap not int list", i)
		}
		Put(pb)
	}
}

func TestGenLenEqualCapSize(t *testing.T) {
	for i := 0; i < _MaxBigSize; i++ {
		pb := Get(i)
		if len(pb.B) != 0 {
			t.Fatalf("len not equal 0 in size %d", i)
		}
		Put(pb)
	}
}

// TestGetUsePutPackGetCapEqualSize tests that the following process
func TestGetUsePutPackGetLenEqualZero(t *testing.T) {
	for i := 0; i < _MaxBigSize; i++ {
		pb := Get(i)
		n := int32(cap(pb.B))
		if n == 0 {
			continue
		}
		pb.B = pb.B[:rand.Int31n(int32(cap(pb.B)))]
		Put(pb)
		pb = Get(i)
		if len(pb.B) != 0 {
			t.Fatalf("len not 0 len in szie %d", i)
		}
		Put(pb)
	}
}

func TestPutFailedCapDiscard(t *testing.T) {
	for i := 0; i < 20; i++ {
		b := &Bytes{B: make([]byte, 20)}
		Put(b)
		pb := Get(20)
		if cap(pb.B) != int(class_to_size[size2class(20)]) {
			t.Fatalf("not crrect cap put in pool")
		}
	}
}

func TestSmallPut(t *testing.T) {
	for i := 0; i < 100; i++ {
		p := make([]byte, 0, 35)
		pb := &Bytes{B: p}
		Put(pb)
	}
	find := false
	for i := 0; i < 100; i++ {
		pb := Get(32)
		if cap(pb.B) == 35 {
			find = true
		}
	}
	if find == false {
		t.Fatal("not find 35 cap buf")
	}
}
