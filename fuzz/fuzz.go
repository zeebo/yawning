package fuzz

import (
	"encoding/binary"
	"fmt"
	"math/bits"

	"github.com/zeebo/yawning"
)

func Fuzz(data []byte) (c int) {
	var it yawning.Intersection
	var bms [8]*yawning.Bitmap
	var ss [8]map[uint32]struct{}

	reset := func() {
		for i := 0; i < 8; i++ {
			bms[i] = yawning.New()
			ss[i] = make(map[uint32]struct{})
		}
	}
	reset()

	c = -1

	var tmp []byte
	for len(data) > 0 {
		tmp, data = read(data, 1)

		switch tmp[0] % 8 {
		case 0: // check
			tmp, data = read(data, 1)
			w := tmp[0]
			if w == 0 {
				w = 255
			}

			fmt.Printf("intersect           %08b\n", w)

			var inc []map[uint32]struct{}

			it.Reset()
			for w > 0 {
				i := bits.TrailingZeros8(w) % 8
				it.Include(bms[i])
				inc = append(inc, ss[i])
				w &= w - 1
			}

			o := intersect(inc)
			if len(o) > 0 && c == -1 {
				c = 0
			} else if len(o) > 20 {
				c = 1
			}

			check(&it, o)

		case 1: // reset

			fmt.Printf("reset\n")

			reset()

		default:
			tmp, data = read(data, 1)
			w := tmp[0]
			if w == 0 {
				w = 255
			}

			tmp, data = read(data, 4)
			v := binary.BigEndian.Uint32(tmp) / 256 * 256

			fmt.Printf("add %-12d to %08b\n", v, w)

			for w > 0 {
				i := bits.TrailingZeros8(w) % 8
				bms[i].Add(v)
				ss[i][v] = struct{}{}
				w &= w - 1
			}

		}
	}

	return c
}

//
// helpers
//

func read(b []byte, n int) ([]byte, []byte) {
	out := make([]byte, n)
	return out, b[copy(out, b):]
}

func check(i *yawning.Intersection, s map[uint32]struct{}) {
	i.Iterate(func(x uint32) bool {
		if _, ex := s[x]; !ex {
			panic(fmt.Sprintf("extra: %v", x))
		}
		delete(s, x)
		return true
	})
	if len(s) != 0 {
		panic(fmt.Sprintf("remaining: %v", s))
	}
}

func intersect(ss []map[uint32]struct{}) map[uint32]struct{} {
	o := make(map[uint32]struct{})

values:
	for x := range ss[0] {
		for _, s := range ss[1:] {
			if _, ok := s[x]; !ok {
				continue values
			}
		}
		o[x] = struct{}{}
	}

	return o
}
