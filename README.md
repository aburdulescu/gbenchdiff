# gbenchdiff
Diff results of two C++ google benchmarks

## Install

```
go get -u github.com/aburdulescu/gbenchdiff/cmd/gbenchdiff
```

## Usage

```
gbenchdiff

Usage: gbenchdiff old.json new.json

For each benchmark in both files, the tool will:
- remove outliers with interquartile range rule
- perform significance test(Man-Whitney U-test)
- print % change in mean from the first to the second file
- print the p-value and sample sizes from a test of the two distributions of benchmark times

Small p-values indicate that the two distributions are significantly different.
If the test indicates that there was no significant change between the two
benchmarks (defined as p > 0.05), a single ~ will be displayed instead of
the percent change.

Notes:
- run the benchmark with --benchmark_repetitions(=10 should be enough)
- run the benchmark with --benchmark_out=file.json
```
