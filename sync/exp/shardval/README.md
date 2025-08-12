# shardval

shardval provides an implementation of a sharded value backed by [sync.Pool](https://pkg.go.dev/sync#Pool).

This is useful for performance-sensitive code that is running across many cores
where contention can cause significant slowdowns.

## Benchmarks

Run on a Macbook Pro (M2 Pro).

| CPUs                          | 1     | 2     | 4      | 8      | 12     |
| ----------------------------- | ----- | ----- | ------ | ------ | ------ |
| BenchmarkCounterLowBound      | 1.981 | 1.023 | 0.5134 | 0.2616 | 0.1906 |
| BenchmarkCounterAtomic        | 3.77  | 13.39 | 22.49  | 35.95  | 70.24  |
| BenchmarkCounterMutex         | 7.748 | 7.538 | 7.535  | 7.529  | 7.549  |
| BenchmarkCounterSharded       | 10.32 | 5.263 | 3.056  | 1.362  | 1.652  |
| BenchmarkCounterShardedAtomic | 12.87 | 6.655 | 4.828  | 3.57   | 2.296  |
| BenchmarkCounterShardedMutex  | 18.42 | 9.366 | 5.209  | 2.449  | 2.722  |

`BenchmarkCounterLowBound` sets a lower bound for the benchmark by incrementing a local integer.

Atomics scale very poorly -- as the number of cores increase, we generally to see increased # of operations (and hence lower ns/op), yet atomics show an increase. So for example, going from 1 core to 2 core, the ns/op increases by 3.5x, but cores increased 2x, so we're spending 7x the total CPU time to do the same work.
