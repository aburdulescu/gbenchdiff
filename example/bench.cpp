#include <benchmark/benchmark.h>

#include <string>

static void BM_foo(benchmark::State& state) {
    auto len = state.range(0);
  for (auto _ : state) {
      std::string s(len, 'x');
  }
}

BENCHMARK(BM_foo)->DenseRange(0,20);

BENCHMARK_MAIN();
