package yawning

import (
	"math/bits"
)

type page struct {
	bit uint64
	buf [1024]uint64
}

func (p *page) copyFrom(q *page) {
	pb, qb := &p.buf, &q.buf
	p.bit = q.bit

	if p.bit != 0 {
		tr := bits.TrailingZeros64(p.bit) * 16
		le := (64 - bits.LeadingZeros64(p.bit)) * 16
		copy(pb[tr:le], qb[tr:le])
	}
}

func (p *page) add(x uint16) {
	bit := x % 64
	idx := x / 64
	p.buf[idx] |= 1 << bit
	p.bit |= 1 << (idx / 16 % 64)
}

func (p *page) get(x uint16) bool {
	bit := x % 64
	idx := x / 64
	return p.bit&(1<<(idx/16%64)) > 0 && p.buf[idx]&(1<<bit) > 0
}

func (p *page) and(q *page) {
	pb, qb := &p.buf, &q.buf
	p.bit &= q.bit

	for bit := p.bit; bit > 0; bit &= bit - 1 {
		idx := uint(bits.TrailingZeros64(bit)) * 16 % 1024
		if and_avx2(&pb[idx], &qb[idx]) == 0 {
			p.bit &^= 1 << (idx / 16 % 64)
		}
	}
}

func (p *page) iterate(id uint32, cb func(uint32) bool) bool {
	pb := &p.buf
	bit := p.bit

	for bit > 0 {
		idx := uint(bits.TrailingZeros64(bit)) % 64
		for i := idx * 16; i < idx*16+16 && i < 1024; i++ {
			off := 64 * uint32(i)
			for pbi := pb[i]; pbi > 0; pbi &= pbi - 1 {
				j := uint32(bits.TrailingZeros64(pbi))
				if !cb(id | j + off) {
					return false
				}
			}
		}
		bit &= bit - 1
	}

	return true
}
