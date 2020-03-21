package yawning

import "math/bits"

type Intersection struct {
	bms  []*Bitmap
	base *Bitmap
	bmin int
	bit  page
	tmp  page
}

func (i *Intersection) Reset() {
	for n := range i.bms {
		i.bms[n] = nil
	}
	i.bms = i.bms[:0]
	i.base = nil
	i.bmin = 0
	i.bit.bit = 0
}

func (i *Intersection) IsEmpty() bool {
	return i.bit.bit == 0
}

func (i *Intersection) Include(bm *Bitmap) {
	if len(i.bms) == 0 {
		i.bit.copyFrom(&bm.bit)
	} else {
		i.bit.and(&bm.bit)
	}
	i.bms = append(i.bms, bm)

	bmin := bits.OnesCount64(bm.bit.bit)
	if i.base == nil || bmin < i.bmin {
		i.base = bm
		i.bmin = bmin
	}
}

func (i *Intersection) Iterate(cb func(uint32) bool) {
	i.bit.iterate(0, func(id uint32) bool {
		if len(i.bms) == 1 {
			return i.base.pages[uint16(id)].iterate(id<<16, cb)
		}
		i.tmp.copyFrom(i.base.pages[uint16(id)])
		for _, bm := range i.bms {
			if bm == i.base {
				continue
			}
			i.tmp.and(bm.pages[uint16(id)])
			if i.tmp.bit == 0 {
				return true
			}
		}
		return i.tmp.iterate(id<<16, cb)
	})
}
