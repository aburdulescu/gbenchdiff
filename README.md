# gbenchdiff

Diff results of two C++ google benchmarks

## Install

### From source using `go`

```
go install github.com/aburdulescu/gbenchdiff/cmd/gbenchdiff@latest
```

### Using prebuilt binary

See [here](https://github.com/aburdulescu/gbenchdiff/releases/latest)

## Usage

```
Usage: gbenchdiff [options] old.json new.json
options:
  -filter string
        select only the benchmarks with names that match the given regex
  -no-ctx
        don't compare benchmark contexts
  -with-cpu
        compare also CPU time

For each benchmark in both files, the tool will:
- remove outliers with interquartile range rule
- perform significance test(Man-Whitney U-test)
- print % change in mean from the first to the second file
- print the p-value and sample sizes from a test of the two distributions of benchmark times

Small p-values indicate that the two distributions are significantly different.
If the test indicates that there was no significant change between the two
benchmarks (defined as p > 0.05), a single ~ will be displayed instead of
the percent change.

IMPORTANT:
Run the benchmark with the following flags:
    --benchmark_out=file.json
    --benchmark_repetitions(=10 should be enough in most cases)
```

For a example, see [example](./example) directory.

## Acknowledgements

Heavily inspired by [Go benchstat](https://github.com/golang/perf) tool
(the statistics code is copied from there).
