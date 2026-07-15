# Complete Test-Bench Design

The final test bench separates **correctness coverage** from the project's
**Victim Cache performance comparison**. A single conflict-only trace cannot
validate every cache level, so the suite uses five deterministic workloads.

## Workloads

### 1. `repeated`

Repeats one byte address. Expected behavior for every topology containing L1:

- first request: L1 miss
- all remaining requests: L1 hit

This is the simplest proof that L1 lookup, insertion, tag calculation, and hit
accounting work.

### 2. `sequential`

Walks a small contiguous block set repeatedly. The default working set fits in
L1, so the first pass warms the cache and later passes hit in L1.

### 3. `conflict`

Separates addresses by the total L1 capacity. In a direct-mapped L1, all those
addresses map to the same line. This is the project's main performance trace:

- `l1-l2`: repeated accesses should hit in L2 after initial memory misses
- `full`: repeated accesses should often hit in Victim Cache before L2
- the Victim Cache should reduce L2 reads and total cycles

### 4. `mixed`

Contains three phases with disjoint L1 indices:

1. one hot address, producing L1 hits
2. a small conflict set that fits in L1 plus Victim Cache, producing Victim hits
3. a larger conflict set that exceeds Victim capacity, producing L2 hits

The full hierarchy must report non-zero L1 hits, Victim hits, L2 hits, and main
memory accesses in the same run.

### 5. `writeback`

Writes more dirty blocks than one L2 set can hold. This forces dirty L2
replacement and verifies L2 writeback and memory traffic counters.

## Architectures

Each workload runs on:

- `memory`
- `l1`
- `l1-l2`
- `full-fifo`
- `full-lru`

## Automatic validation

The test bench checks:

- request and response counts
- L1 hit/miss accounting
- Victim hit/miss accounting
- L2 hit/miss/read accounting
- response-location counters
- total cycle accounting
- repeated-address L1 warm-up
- sequential L1 reuse
- conflict-trace Victim Cache benefit
- mixed-trace coverage of every level
- dirty L2 writeback behavior

The default command exits non-zero if any check fails:

```bash
go run ./cmd/testbench
```

List all individual checks:

```bash
go run ./cmd/testbench -verbose-checks
```

## CSV output

```bash
go run ./cmd/testbench -csv results.csv
```

The CSV contains per-workload, per-architecture values for cycles, average
cycles per request, all hit/miss counters and rates, Victim swaps, L2 reads and
writes, and memory accesses.
