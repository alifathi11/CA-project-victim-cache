package config

import (
	"fmt"
	"strings"
)

const (
	ReplacementFIFO = "FIFO"
	ReplacementLRU  = "LRU"
)

// Config contains the parameters that may be varied during experiments.
// Sizes are expressed in bytes and latencies are expressed in cycles.
type Config struct {
	BlockSizeBytes      uint64
	L1SizeBytes         uint64
	L1Associativity     int
	VictimEntries       int
	VictimEnabled       bool
	VictimPolicy        string
	L2SizeBytes         uint64
	L2Associativity     int
	L1HitLatencyCycles  uint64
	VictimLatencyCycles uint64
	L2LatencyCycles     uint64
	MemoryLatencyCycles uint64
}

func Default() Config {
	return Config{
		BlockSizeBytes:      64,
		L1SizeBytes:         4 * 1024,
		L1Associativity:     1,
		VictimEntries:       8,
		VictimEnabled:       false,
		VictimPolicy:        ReplacementFIFO,
		L2SizeBytes:         64 * 1024,
		L2Associativity:     8,
		L1HitLatencyCycles:  1,
		VictimLatencyCycles: 2,
		L2LatencyCycles:     12,
		MemoryLatencyCycles: 100,
	}
}

// NormalizedVictimPolicy returns the configured replacement policy in the
// canonical form used by the cache implementation.
func (c Config) NormalizedVictimPolicy() string {
	return strings.ToUpper(strings.TrimSpace(c.VictimPolicy))
}

// ValidateL1 checks the invariants required by the direct-mapped L1.
func (c Config) ValidateL1() error {
	if c.BlockSizeBytes == 0 {
		return fmt.Errorf("block size must be greater than zero")
	}
	if c.L1SizeBytes == 0 || c.L1SizeBytes%c.BlockSizeBytes != 0 {
		return fmt.Errorf("L1 size must be non-zero and divisible by block size")
	}
	if c.L1Associativity != 1 {
		return fmt.Errorf("L1 must be direct-mapped (associativity = 1)")
	}
	return nil
}

// ValidateL2 checks the invariants required by the set-associative L2.
func (c Config) ValidateL2() error {
	if c.BlockSizeBytes == 0 {
		return fmt.Errorf("block size must be greater than zero")
	}
	if c.L2Associativity <= 0 {
		return fmt.Errorf("L2 associativity must be greater than zero")
	}
	l2SetBytes := c.BlockSizeBytes * uint64(c.L2Associativity)
	if c.L2SizeBytes == 0 || c.L2SizeBytes%l2SetBytes != 0 {
		return fmt.Errorf("L2 size must be non-zero and divisible by block size times associativity")
	}
	return nil
}

// ValidateVictim checks the invariants required by the Victim Cache.
func (c Config) ValidateVictim() error {
	if c.BlockSizeBytes == 0 {
		return fmt.Errorf("block size must be greater than zero")
	}
	if c.VictimEntries < 0 {
		return fmt.Errorf("victim cache entry count cannot be negative")
	}
	policy := c.NormalizedVictimPolicy()
	if policy != ReplacementFIFO && policy != ReplacementLRU {
		return fmt.Errorf("unsupported victim cache replacement policy %q", c.VictimPolicy)
	}
	return nil
}

// ValidateMemoryHierarchy checks all invariants required by the memory and
// cache components. It deliberately does not validate benchmark or CPU settings.
func (c Config) ValidateMemoryHierarchy() error {
	if err := c.ValidateL1(); err != nil {
		return err
	}
	if err := c.ValidateL2(); err != nil {
		return err
	}
	if err := c.ValidateVictim(); err != nil {
		return err
	}
	return nil
}
