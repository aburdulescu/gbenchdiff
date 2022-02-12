# Example

## Prerequisites

- Install [gbenchdiff](https://github.com/aburdulescu/gbenchdiff#install)
- Install [google benchmark library](https://github.com/google/benchmark#installation)

## Run the example

- Compile it: `make`
- Run the benchmark: `./bench --benchmark_out=old.json --benchmark_repetitions=10`
- Start some CPU consuming apps(e.g. web browser playing a video) so we have some variance between the benchmark runs
- Run the benchmark again: `./bench --benchmark_out=new.json --benchmark_repetitions=10`

Now you can run `gbenchdiff` on the two benchmark results:

`gbenchdiff old.json new.json`

It should print something like this:

```
real time  delta    note            old      new
---------  -----    ----            ---      ---
BM_foo/0   -0.17%   (p=0.03 n=4+4)  4.02ns   4.02ns
BM_foo/2   +10.31%  (p=0.02 n=4+5)  7.64ns   8.42ns
BM_foo/4   +25.65%  (p=0.02 n=5+5)  7.24ns   9.10ns
BM_foo/6   +8.91%   (p=0.01 n=5+5)  4.44ns   4.84ns
BM_foo/8   +7.27%   (p=0.01 n=5+5)  4.50ns   4.83ns
BM_foo/10  +11.73%  (p=0.01 n=5+5)  4.33ns   4.83ns
BM_foo/12  +11.74%  (p=0.01 n=5+5)  4.33ns   4.84ns
BM_foo/14  +11.21%  (p=0.01 n=5+5)  4.35ns   4.84ns
BM_foo/16  +7.73%   (p=0.01 n=5+5)  18.84ns  20.30ns
BM_foo/18  +7.91%   (p=0.01 n=5+5)  18.83ns  20.32ns
BM_foo/20  +7.76%   (p=0.01 n=5+5)  18.86ns  20.32ns
```
