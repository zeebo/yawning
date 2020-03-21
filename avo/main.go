package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	{
		TEXT("and_avx2", 0, `func(a, b *uint64) uint32`)

		var (
			a = Mem{Base: Load(Param("a"), GP64())}
			b = Mem{Base: Load(Param("b"), GP64())}
		)

		ya, yb := YMM(), YMM()

		mr, r := GP32(), GP32()
		MOVL(U32(1), mr)
		MOVL(U32(0), r)

		for i := 0; i < 4; i++ {
			VMOVDQU(a.Offset(32*i), ya)
			VMOVDQU(b.Offset(32*i), yb)
			VPTEST(ya, yb)
			CMOVLNE(mr, r)
			VPAND(ya, yb, ya)
			VMOVDQU(ya, a.Offset(32*i))
		}

		VZEROUPPER()
		Store(r, ReturnIndex(0))
		RET()
	}

	Generate()
}
