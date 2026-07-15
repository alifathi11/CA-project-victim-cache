package benchmark

import (
	"fmt"
	"strings"

	"victimcacheproject/internal/model"
)

// TraceKind identifies a deterministic memory-access workload.
type TraceKind string

const (
	TraceRepeated   TraceKind = "repeated"
	TraceSequential TraceKind = "sequential"
	TraceConflict   TraceKind = "conflict"
	TraceMixed      TraceKind = "mixed"
	TraceWriteback  TraceKind = "writeback"
)

// SuiteConfig contains the architectural information needed to construct
// traces that deliberately exercise specific levels of the hierarchy.
type SuiteConfig struct {
	BaseAddress     uint64
	BlockSizeBytes  uint64
	L1SizeBytes     uint64
	L2SizeBytes     uint64
	L2Associativity int
	VictimEntries   int
	NumberOfBlocks  int
	Repetitions     int
	AccessSizeBytes uint64
}

// Scenario is one named workload in the complete test bench.
type Scenario struct {
	Kind        TraceKind
	Name        string
	Description string
	Requests    []model.Request
}

func AllTraceKinds() []TraceKind {
	return []TraceKind{
		TraceRepeated,
		TraceSequential,
		TraceConflict,
		TraceMixed,
		TraceWriteback,
	}
}

func ParseTraceKind(value string) (TraceKind, error) {
	kind := TraceKind(strings.ToLower(strings.TrimSpace(value)))
	for _, supported := range AllTraceKinds() {
		if kind == supported {
			return kind, nil
		}
	}
	return "", fmt.Errorf("unsupported trace %q; expected repeated, sequential, conflict, mixed, writeback, or all", value)
}

// GenerateSuite creates all workloads used by the final project evaluation.
func GenerateSuite(cfg SuiteConfig) ([]Scenario, error) {
	scenarios := make([]Scenario, 0, len(AllTraceKinds()))
	for _, kind := range AllTraceKinds() {
		scenario, err := GenerateScenario(kind, cfg)
		if err != nil {
			return nil, err
		}
		scenarios = append(scenarios, scenario)
	}
	return scenarios, nil
}

// GenerateScenario creates one deterministic workload.
func GenerateScenario(kind TraceKind, cfg SuiteConfig) (Scenario, error) {
	if err := validateSuiteConfig(cfg); err != nil {
		return Scenario{}, err
	}

	var scenario Scenario
	scenario.Kind = kind
	switch kind {
	case TraceRepeated:
		scenario.Name = "Repeated locality"
		scenario.Description = "Repeats one address so the first access misses and all later accesses hit in L1."
		scenario.Requests = generateRepeatedTrace(cfg)
	case TraceSequential:
		scenario.Name = "Sequential working set"
		scenario.Description = "Walks a small contiguous working set repeatedly; after warm-up, the blocks should remain in L1."
		scenario.Requests = generateSequentialTrace(cfg)
	case TraceConflict:
		scenario.Name = "L1 conflict thrashing"
		scenario.Description = "Uses addresses separated by the L1 capacity so every block maps to the same direct-mapped L1 line."
		scenario.Requests = GenerateConflictTrace(ConflictConfig{
			BaseAddress:     cfg.BaseAddress,
			CacheSizeBytes:  cfg.L1SizeBytes,
			NumberOfBlocks:  cfg.NumberOfBlocks,
			Repetitions:     cfg.Repetitions,
			AccessSizeBytes: cfg.AccessSizeBytes,
		})
	case TraceMixed:
		scenario.Name = "Mixed hierarchy coverage"
		scenario.Description = "Combines an L1-hot phase, a victim-cache-sized conflict phase, and an oversized conflict phase that reaches L2 and memory."
		scenario.Requests = generateMixedTrace(cfg)
	case TraceWriteback:
		scenario.Name = "Dirty writeback pressure"
		scenario.Description = "Writes more blocks than one L2 set can hold, forcing dirty evictions and memory writebacks."
		scenario.Requests = generateWritebackTrace(cfg)
	default:
		return Scenario{}, fmt.Errorf("unsupported trace kind %q", kind)
	}

	return scenario, nil
}

func validateSuiteConfig(cfg SuiteConfig) error {
	if cfg.BlockSizeBytes == 0 {
		return fmt.Errorf("block size must be greater than zero")
	}
	if cfg.L1SizeBytes == 0 || cfg.L1SizeBytes%cfg.BlockSizeBytes != 0 {
		return fmt.Errorf("L1 size must be non-zero and divisible by block size")
	}
	if cfg.L2SizeBytes == 0 || cfg.L2Associativity <= 0 {
		return fmt.Errorf("L2 size and associativity must be greater than zero")
	}
	if cfg.L2SizeBytes%(cfg.BlockSizeBytes*uint64(cfg.L2Associativity)) != 0 {
		return fmt.Errorf("L2 size must be divisible by block size times associativity")
	}
	if cfg.NumberOfBlocks <= 0 {
		return fmt.Errorf("number of blocks must be greater than zero")
	}
	if cfg.Repetitions <= 0 {
		return fmt.Errorf("repetitions must be greater than zero")
	}
	if cfg.AccessSizeBytes == 0 || cfg.AccessSizeBytes > cfg.BlockSizeBytes {
		return fmt.Errorf("access size must be between 1 and the block size")
	}
	if cfg.VictimEntries < 0 {
		return fmt.Errorf("victim entries cannot be negative")
	}
	return nil
}

func generateRepeatedTrace(cfg SuiteConfig) []model.Request {
	total := cfg.NumberOfBlocks * cfg.Repetitions
	requests := make([]model.Request, 0, total)
	for i := 0; i < total; i++ {
		requests = appendRequest(requests, cfg.BaseAddress, model.Read, cfg.AccessSizeBytes)
	}
	return requests
}

func generateSequentialTrace(cfg SuiteConfig) []model.Request {
	requests := make([]model.Request, 0, cfg.NumberOfBlocks*cfg.Repetitions)
	for repetition := 0; repetition < cfg.Repetitions; repetition++ {
		for block := 0; block < cfg.NumberOfBlocks; block++ {
			address := cfg.BaseAddress + uint64(block)*cfg.BlockSizeBytes
			requests = appendRequest(requests, address, model.Read, cfg.AccessSizeBytes)
		}
	}
	return requests
}

func generateMixedTrace(cfg SuiteConfig) []model.Request {
	requests := make([]model.Request, 0)

	// Phase 1: keep one line hot to produce unambiguous L1 hits.
	hotAddress := cfg.BaseAddress + cfg.BlockSizeBytes
	hotAccesses := maxInt(4, cfg.Repetitions)
	for i := 0; i < hotAccesses; i++ {
		requests = appendRequest(requests, hotAddress, model.Read, cfg.AccessSizeBytes)
	}

	// Phase 2: use a conflict set that fits in L1 + Victim Cache. After the
	// first round, requests should be recovered from the Victim Cache.
	smallConflictBlocks := maxInt(2, cfg.NumberOfBlocks)
	if cfg.VictimEntries > 0 && smallConflictBlocks > cfg.VictimEntries+1 {
		smallConflictBlocks = cfg.VictimEntries + 1
	}
	smallConflictBase := cfg.BaseAddress + 2*cfg.BlockSizeBytes
	for repetition := 0; repetition < 2; repetition++ {
		for block := 0; block < smallConflictBlocks; block++ {
			address := smallConflictBase + uint64(block)*cfg.L1SizeBytes
			requests = appendRequest(requests, address, model.Read, cfg.AccessSizeBytes)
		}
	}

	// Phase 3: exceed the combined L1 + Victim capacity at one L1 index.
	// The second round must therefore reach L2 for at least some requests.
	largeConflictBlocks := maxInt(cfg.NumberOfBlocks+2, cfg.VictimEntries+2)
	largeConflictBase := cfg.BaseAddress + 3*cfg.BlockSizeBytes
	for repetition := 0; repetition < 2; repetition++ {
		for block := 0; block < largeConflictBlocks; block++ {
			address := largeConflictBase + uint64(block)*cfg.L1SizeBytes
			requests = appendRequest(requests, address, model.Read, cfg.AccessSizeBytes)
		}
	}

	return requests
}

func generateWritebackTrace(cfg SuiteConfig) []model.Request {
	// Number of sets * block size is the byte stride that maps addresses to
	// the same L2 set. Writing associativity+2 blocks guarantees L2 eviction.
	l2SetSpanBytes := cfg.L2SizeBytes / uint64(cfg.L2Associativity)
	blockCount := cfg.L2Associativity + 2
	base := cfg.BaseAddress + 4*cfg.BlockSizeBytes
	requests := make([]model.Request, 0, blockCount)
	for block := 0; block < blockCount; block++ {
		address := base + uint64(block)*l2SetSpanBytes
		requests = appendRequest(requests, address, model.Write, cfg.AccessSizeBytes)
	}
	return requests
}

func appendRequest(requests []model.Request, address uint64, op model.Operation, size uint64) []model.Request {
	return append(requests, model.Request{
		ID:      uint64(len(requests) + 1),
		Address: address,
		Op:      op,
		Size:    size,
	})
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
