package benchmark

import "victimcacheproject/internal/model"

type ConflictConfig struct {
	BaseAddress     uint64
	CacheSizeBytes  uint64
	NumberOfBlocks  int
	Repetitions     int
	AccessSizeBytes uint64
}

// GenerateConflictTrace will create addresses separated by the L1 capacity,
// making them map to the same direct-mapped L1 index.
func GenerateConflictTrace(cfg ConflictConfig) []model.Request {
	// TODO(stage-2): Generate a deterministic alternating access trace.
	return nil
}
