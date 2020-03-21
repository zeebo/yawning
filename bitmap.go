package yawning

type Bitmap struct {
	pages map[uint16]*page
	pid   uint16
	pp    *page
	bit   page
}

func New() *Bitmap {
	return &Bitmap{pages: make(map[uint16]*page)}
}

func (b *Bitmap) Add(x uint32) {
	id := uint16(x >> 16)
	if b.pp != nil && id == b.pid {
		b.pp.add(uint16(x))
		return
	}

	p, ok := b.pages[id]
	if !ok {
		p = new(page)
		b.pages[id] = p
		b.bit.add(id)
	}
	p.add(uint16(x))

	b.pid = id
	b.pp = p
}

func (b *Bitmap) AddMany(xs []uint32) {
	for _, x := range xs {
		b.Add(x)
	}
}
