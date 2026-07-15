# Victim Cache Project 6

A configurable functional reference simulator for four memory hierarchies:

- `memory`: CPU -> Main Memory
- `l1`: CPU -> L1 -> Main Memory
- `l1-l2`: CPU -> L1 -> L2 -> Main Memory
- `full`: CPU -> L1 -> Victim Cache -> L2 -> Main Memory

## Run one topology and one workload

The simulator supports five deterministic traces:

- `repeated`: proves L1 warm-up and L1 hits
- `sequential`: repeatedly scans a small working set that fits in L1
- `conflict`: forces direct-mapped L1 thrashing and measures Victim Cache benefit
- `mixed`: deliberately generates L1 hits, Victim hits, L2 hits, and memory accesses in one run
- `writeback`: generates dirty evictions and L2-to-memory writebacks

Examples:

```bash
go run ./cmd/sim -topology l1 -trace repeated
go run ./cmd/sim -topology l1-l2 -trace sequential
go run ./cmd/sim -topology l1-l2 -trace conflict
go run ./cmd/sim -topology full -trace mixed -victim=true -victim-policy=FIFO
go run ./cmd/sim -topology full -trace writeback -victim=true -victim-policy=LRU
```

Trace controls:

```bash
go run ./cmd/sim -topology full -trace conflict -blocks 4 -repetitions 20
```

## Complete final test bench

Run every workload on every hierarchy, test both FIFO and LRU Victim Cache policies, and execute automatic accounting and behavioral checks:

```bash
go run ./cmd/testbench
```

Run one workload only:

```bash
go run ./cmd/testbench -trace mixed
```

Write the measurements to CSV for the final report:

```bash
go run ./cmd/testbench -csv results.csv
```

Useful options:

```bash
go run ./cmd/testbench -blocks 4 -repetitions 20 -victim-policy BOTH
go run ./cmd/testbench -trace conflict -victim-policy FIFO
go run ./cmd/testbench -strict=false
go run ./cmd/testbench -verbose-checks
```

`-strict=true` is the default. The command exits with status 1 when any validation check fails, which makes it suitable for CI.

## Compare command

`cmd/compare` remains as a shorter alias for matrix reporting:

```bash
go run ./cmd/compare -trace conflict
go run ./cmd/compare -trace all -victim-policy BOTH
```

## Verification

```bash
go test ./...
go test -race ./...
go vet ./...
```

## Akita boundary

The code under `internal/system` is the functional reference model. An Akita integration should keep this model as the correctness oracle and replace the synchronous adapter with Akita components, ports, messages, and scheduled events. See `AKITA_INTEGRATION.md`.
