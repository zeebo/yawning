package yawning

import (
	"fmt"
	"testing"

	"github.com/zeebo/assert"
)

func TestIntersection_Crasher(t *testing.T) {
	x := 0
	run := func(cb func(t *testing.T)) { x++; t.Run(fmt.Sprint(x), cb) }

	run(func(t *testing.T) {
		// program:
		// add 876163072    to 00000001
		// intersect           00000001

		bm := New()
		bm.Add(876163072)
		it := new(Intersection)
		it.Include(bm)
		it.Iterate(func(x uint32) bool {
			assert.Equal(t, 876163072, x)
			return true
		})
	})

	run(func(t *testing.T) {
		// program:
		// add 808464384    to 01101111
		// intersect           01101111
		// add 808517376    to 01101111
		// intersect           01101111

		var bms [6]*Bitmap
		for i := range bms {
			bms[i] = New()
			bms[i].Add(808464384)
		}

		it := new(Intersection)
		it.Reset()
		for i := range bms {
			it.Include(bms[i])
		}
		it.Iterate(func(x uint32) bool {
			assert.Equal(t, x, 808464384)
			return true
		})

		for i := range bms {
			bms[i].Add(808517376)
		}

		it.Reset()
		for i := range bms {
			it.Include(bms[i])
		}
		it.Iterate(func(x uint32) bool {
			assert.That(t, x == 808464384 || x == 808517376)
			return true
		})
	})

	run(func(t *testing.T) {
		// program:
		// intersect           11111111

		var bms [8]*Bitmap
		for i := range bms {
			bms[i] = New()
		}

		it := new(Intersection)
		it.Reset()
		for i := range bms {
			it.Include(bms[i])
		}
		it.Iterate(func(x uint32) bool {
			t.FailNow()
			return true
		})
	})
}
