package system

import (
	"fmt"

	"victimcacheproject/internal/cache"
	"victimcacheproject/internal/config"
	"victimcacheproject/internal/memory"
	"victimcacheproject/internal/metrics"
)

type System struct {
	Config config.Config
	L1     *cache.L1
	Victim *cache.Victim
	L2     *cache.L2
	Memory *memory.MainMemory
	Stats  metrics.Stats
}

func New(cfg config.Config) *System {
	return &System{
		Config: cfg,
		L1:     cache.NewL1(cfg),
		Victim: cache.NewVictim(cfg),
		L2:     cache.NewL2(cfg),
		Memory: memory.NewMainMemory(),
	}
}

func (s *System) Validate() error {
	if s.Config.BlockSizeBytes == 0 {
		return fmt.Errorf("block size must be non-zero")
	}
	if s.Config.L1SizeBytes%s.Config.BlockSizeBytes != 0 {
		return fmt.Errorf("L1 size must be divisible by block size")
	}
	if s.Config.L1Associativity != 1 {
		return fmt.Errorf("project baseline expects a direct-mapped L1")
	}
	if s.Config.VictimEntries < 0 {
		return fmt.Errorf("victim entries cannot be negative")
	}
	return nil
}
