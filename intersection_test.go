package yawning

import (
	"testing"

	"github.com/zeebo/assert"
)

func TestIntersection(t *testing.T) {
	var it Intersection
	for i := 1; i < 5; i++ {
		bm := New()
		for j := 0; j < 60; j += i {
			bm.Add(uint32(j))
		}
		bm.Add(1000000)
		it.Include(bm)
	}

	var got []uint32
	it.Iterate(func(n uint32) bool {
		got = append(got, n)
		return true
	})
	assert.DeepEqual(t, got, []uint32{0, 12, 24, 36, 48, 1000000})
}

func BenchmarkIntersection(b *testing.B) {
	var bms []*Bitmap
	for i := 1; i < 5; i++ {
		bm := New()
		for j := 0; j < 60; j += i {
			bm.Add(uint32(j))
		}
		bms = append(bms, bm)
	}
	bms = append(bms, New())

	var it Intersection
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it.Reset()
		for _, bm := range bms {
			it.Include(bm)
		}
		it.Iterate(func(uint32) bool { return true })
	}
}
