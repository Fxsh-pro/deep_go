package main

import "testing"

var benchmarkResult uint32

func BenchmarkToLittleEndian(b *testing.B) {
	b.Run("Shifts", func(b *testing.B) {
		for i := range b.N {
			benchmarkResult = toLittleEndianShifts(uint32(i))
		}
	})

	b.Run("Pointer", func(b *testing.B) {
		for i := range b.N {
			benchmarkResult = ToLittleEndian(uint32(i)) // быстрее
		}
	})
}
