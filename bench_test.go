package bpool

import (
	"strconv"
	"sync"
	"testing"
)

var b1 []byte

func BenchmarkMakeAndGet(b *testing.B) {
	sizes := []int{16, 32, 48, 67, 96, 128, 256, 147, 480, 512, 684, 1024, 1355, 1570, 1740, 2048, 2576, 2892, 4096}
	for _, v := range sizes {
		v := v
		b.Run("Make"+strconv.Itoa(v), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b1 = make([]byte, v)
			}
		})
		b.Run("Get"+strconv.Itoa(v), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b2 := Get(v)
				Put(b2)
			}
		})
		b.Run("Empty sync.Pool"+strconv.Itoa(v), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v1 := emptyPool.Get()
				emptyPool.Put(v1)
			}
		})
	}
}

var emptyPool = sync.Pool{New: func() any {
	return &struct1{a: 99}
}}

type struct1 struct {
	a int
}

/**
goos: darwin
goarch: amd64
pkg: github.com/newacorn/simple-bytes-pool
cpu: 13th Gen Intel(R) Core(TM) i9-13900KS
BenchmarkMakeAndGet
BenchmarkMakeAndGet/Make16
BenchmarkMakeAndGet/Make16-32         	111400137	        11.41 ns/op	      16 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get16
BenchmarkMakeAndGet/Get16-32          	131064068	         7.655 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet16
BenchmarkMakeAndGet/MakeAndGet16-32   	227290814	         5.329 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make32
BenchmarkMakeAndGet/Make32-32         	96488462	        12.21 ns/op	      32 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get32
BenchmarkMakeAndGet/Get32-32          	140714461	         7.615 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet32
BenchmarkMakeAndGet/MakeAndGet32-32   	238418749	         5.753 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make48
BenchmarkMakeAndGet/Make48-32         	89292836	        13.63 ns/op	      48 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get48
BenchmarkMakeAndGet/Get48-32          	152496171	         7.625 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet48
BenchmarkMakeAndGet/MakeAndGet48-32   	217547701	         4.915 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make67
BenchmarkMakeAndGet/Make67-32         	71279605	        16.44 ns/op	      80 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get67
BenchmarkMakeAndGet/Get67-32          	162187579	         7.331 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet67
BenchmarkMakeAndGet/MakeAndGet67-32   	245798692	         4.885 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make96
BenchmarkMakeAndGet/Make96-32         	68784266	        17.75 ns/op	      96 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get96
BenchmarkMakeAndGet/Get96-32          	157987567	         7.575 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet96
BenchmarkMakeAndGet/MakeAndGet96-32   	242016192	         5.150 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make128
BenchmarkMakeAndGet/Make128-32        	56022525	        20.52 ns/op	     128 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get128
BenchmarkMakeAndGet/Get128-32         	157150724	         7.715 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet128
BenchmarkMakeAndGet/MakeAndGet128-32  	216965518	         5.692 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make256
BenchmarkMakeAndGet/Make256-32        	34112163	        36.85 ns/op	     256 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get256
BenchmarkMakeAndGet/Get256-32         	141303920	         8.302 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet256
BenchmarkMakeAndGet/MakeAndGet256-32  	249729206	         5.894 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make147
BenchmarkMakeAndGet/Make147-32        	47022843	        25.24 ns/op	     160 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get147
BenchmarkMakeAndGet/Get147-32         	147954907	         7.835 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet147
BenchmarkMakeAndGet/MakeAndGet147-32  	227494580	         5.429 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make480
BenchmarkMakeAndGet/Make480-32        	20983359	        56.58 ns/op	     480 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get480
BenchmarkMakeAndGet/Get480-32         	148042364	         7.527 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet480
BenchmarkMakeAndGet/MakeAndGet480-32  	236801770	         5.580 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make512
BenchmarkMakeAndGet/Make512-32        	19790715	        60.43 ns/op	     512 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get512
BenchmarkMakeAndGet/Get512-32         	135619582	         8.700 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet512
BenchmarkMakeAndGet/MakeAndGet512-32  	231877114	         4.906 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make684
BenchmarkMakeAndGet/Make684-32        	14137590	        84.59 ns/op	     704 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get684
BenchmarkMakeAndGet/Get684-32         	157596867	         7.841 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet684
BenchmarkMakeAndGet/MakeAndGet684-32  	197994230	         5.615 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make1024
BenchmarkMakeAndGet/Make1024-32       	10140279	       118.1 ns/op	    1024 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get1024
BenchmarkMakeAndGet/Get1024-32        	153911330	         8.085 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet1024
BenchmarkMakeAndGet/MakeAndGet1024-32 	234782688	         4.989 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make1355
BenchmarkMakeAndGet/Make1355-32       	 8040349	       153.5 ns/op	    1408 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get1355
BenchmarkMakeAndGet/Get1355-32        	149552971	         8.155 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet1355
BenchmarkMakeAndGet/MakeAndGet1355-32 	242174193	         5.113 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make1570
BenchmarkMakeAndGet/Make1570-32       	 6675844	       181.5 ns/op	    1792 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get1570
BenchmarkMakeAndGet/Get1570-32        	149395496	         8.286 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet1570
BenchmarkMakeAndGet/MakeAndGet1570-32 	245850967	         5.174 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make1740
BenchmarkMakeAndGet/Make1740-32       	 6598102	       184.1 ns/op	    1792 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get1740
BenchmarkMakeAndGet/Get1740-32        	151113766	         8.452 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet1740
BenchmarkMakeAndGet/MakeAndGet1740-32 	231839134	         5.049 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make2048
BenchmarkMakeAndGet/Make2048-32       	 5046399	       241.5 ns/op	    2048 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get2048
BenchmarkMakeAndGet/Get2048-32        	155824824	         7.996 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet2048
BenchmarkMakeAndGet/MakeAndGet2048-32 	240016550	         5.103 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make2576
BenchmarkMakeAndGet/Make2576-32       	 3749134	       329.4 ns/op	    2688 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get2576
BenchmarkMakeAndGet/Get2576-32        	149494459	         7.709 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet2576
BenchmarkMakeAndGet/MakeAndGet2576-32 	213246900	         5.137 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make2892
BenchmarkMakeAndGet/Make2892-32       	 4066686	       294.5 ns/op	    3072 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get2892
BenchmarkMakeAndGet/Get2892-32        	151844942	         8.753 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet2892
BenchmarkMakeAndGet/MakeAndGet2892-32 	227737364	         5.577 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/Make4096
BenchmarkMakeAndGet/Make4096-32       	 2523483	       519.3 ns/op	    4096 B/op	       1 allocs/op
BenchmarkMakeAndGet/Get4096
BenchmarkMakeAndGet/Get4096-32        	145016571	         8.157 ns/op	       0 B/op	       0 allocs/op
BenchmarkMakeAndGet/MakeAndGet4096
BenchmarkMakeAndGet/MakeAndGet4096-32 	209950219	         5.061 ns/op	       0 B/op	       0 allocs/op
PASS
*/
