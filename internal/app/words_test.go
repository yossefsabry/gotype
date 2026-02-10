package app

import (
	"math/rand"
	"testing"
)

func BenchmarkGeneratorBuild(b *testing.B) {
	gen := NewGenerator(rand.NewSource(1))
	options := Options{Punctuation: true, Numbers: true}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.Build(200, options)
	}
}
