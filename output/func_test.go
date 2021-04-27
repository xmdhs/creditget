package output

import (
	"testing"
	_ "time/tzdata"
)

func Benchmark_escape(B *testing.B) {
	for i := 0; i < B.N; i++ {
		escape(`1234567891011121314151617181920`)
	}
}
