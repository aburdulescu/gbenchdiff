CFLAGS = -g -O3 -Wall -fno-omit-frame-pointer
LFLAGS = -lbenchmark -lpthread

BINARY = bench

build:
	clang++ $(CFLAGS) -o $(BINARY) bench.cpp $(LFLAGS)

run: build
	./$(BINARY) --benchmark_repetitions=10 --benchmark_out=bench.json

clean:
	rm -f $(BINARY)
