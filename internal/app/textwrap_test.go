package app

import (
	"math/rand"
	"testing"
)

func BenchmarkBuildLines(b *testing.B) {
	gen := NewGenerator(rand.NewSource(1))
	target := gen.Build(200, Options{})
	width := 60
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buildLines(target, width)
	}
}
