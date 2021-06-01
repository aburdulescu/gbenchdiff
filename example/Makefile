build:
	clang++ -O3 -o bench bench.cpp -lbenchmark -lpthread

run: build
	./bench --benchmark_out=bench.json
