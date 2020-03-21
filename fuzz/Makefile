.DEFAULT_GOAL = run

.PHONY: run
run: build fuzz

.PHONY: run-lf
run-lf: build-lf fuzz-lf

.PHONY: clean
clean:
	rm -rf corpus
	rm -rf crashers
	rm -rf suppressions
	rm -f fuzz-fuzz.zip
	rm -f *-fuzz.a
	rm -f fuzzer
	rm -f crash-*

.PHONY: build
build:
	go-fuzz-build

.PHONY: fuzz
fuzz:
	go-fuzz

.PHONY: build-lf
build-lf:
	go-fuzz-build -libfuzzer
	clang -fsanitize=fuzzer *-fuzz.a -o fuzzer

.PHONY: fuzz-lf
fuzz-lf:
	./fuzzer
