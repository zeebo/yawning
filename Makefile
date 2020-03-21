avx2.s: ./avo/*.go
	( cd avo; go run . ) > avx2.s

.PHONY: fuzz
fuzz:
	( cd fuzz; make )
