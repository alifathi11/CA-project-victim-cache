package config

// Config contains parameters that should be varied during experiments.
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
		VictimPolicy:        "FIFO",
		L2SizeBytes:         64 * 1024,
		L2Associativity:     8,
		L1HitLatencyCycles:  1,
		VictimLatencyCycles: 2,
		L2LatencyCycles:     12,
		MemoryLatencyCycles: 100,
	}
}
