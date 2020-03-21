package yawning

import (
	"testing"
	"time"

	"github.com/zeebo/assert"
	"github.com/zeebo/wyhash"
)

func TestPage(t *testing.T) {
	all := new(page)
	for i := 0; i < 1<<16; i++ {
		all.add(uint16(i))
	}

	t.Run("AndAll", func(t *testing.T) {
		for i := 0; i < 1<<16; i++ {
			p := new(page)
			p.add(uint16(i))
			p.and(all)
			assert.That(t, p.get(uint16(i)))
		}
	})

	t.Run("AndFuzz", func(t *testing.T) {
		ps, qs := make(map[uint16]struct{}), make(map[uint16]struct{})
		p, q := new(page), new(page)

		for i := 0; i < 2048; i++ {
			pi, qi := uint16(wyhash.Uint64()), uint16(wyhash.Uint64())
			ps[pi] = struct{}{}
			p.add(pi)
			qs[qi] = struct{}{}
			q.add(qi)
		}

		p.and(q)

		for i := range ps {
			if _, ok := qs[i]; ok {
				assert.That(t, p.get(i))
			} else {
				assert.That(t, !p.get(i))
			}
		}
	})

	t.Run("IterFuzz", func(t *testing.T) {
		s := make(map[uint16]struct{})
		p := new(page)

		for i := 0; i < 2048; i++ {
			x := uint16(wyhash.Uint64())
			s[x] = struct{}{}
			p.add(x)
		}

		p.iterate(0, func(x uint32) bool {
			_, ok := s[uint16(x)]
			assert.That(t, ok)
			delete(s, uint16(x))
			return true
		})
		assert.That(t, len(s) == 0)
	})
}

func BenchmarkPage(b *testing.B) {
	none := new(page)
	all := new(page)
	for i := 0; i < 1<<16; i++ {
		all.add(uint16(i))
	}

	b.Run("AndNone", func(b *testing.B) {
		p := new(page)
		for i := 0; i < b.N; i++ {
			p.and(none)
		}
	})

	b.Run("AndAll", func(b *testing.B) {
		p := *all
		for i := 0; i < b.N; i++ {
			p.and(all)
		}
	})

	b.Run("Iterate", func(b *testing.B) {
		p := new(page)
		for i := 0; i < 1<<16; i++ {
			p.add(uint16(i))
		}

		b.ReportAllocs()
		b.ResetTimer()

		now := time.Now()
		for i := 0; i < b.N; i++ {
			p.iterate(0, func(uint32) bool { return true })
		}

		b.ReportMetric(float64(time.Since(now))/(1<<16)/float64(b.N), "ns/val")
	})
}
