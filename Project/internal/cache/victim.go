package cache

import (
	"victimcacheproject/internal/config"
	"victimcacheproject/internal/model"
)

type Victim struct {
	cfg config.Config
	// TODO(stage-4): Add fully-associative entries and FIFO/LRU metadata.
}

func NewVictim(cfg config.Config) *Victim { return &Victim{cfg: cfg} }
func (c *Victim) Name() string            { return "VictimCache" }

func (c *Victim) Lookup(req model.Request) LookupResult {
	// TODO(stage-4): Search every valid victim entry.
	return LookupResult{}
}

func (c *Victim) Insert(block model.Block) *model.Block {
	// TODO(stage-4): Insert using FIFO/LRU and return overflow eviction.
	return nil
}

func (c *Victim) Invalidate(blockAddress uint64) bool {
	// TODO(stage-4): Remove matching block.
	return false
}

// Swap should exchange the requested victim block with the L1 eviction.
func (c *Victim) Swap(requestedBlockAddress uint64, l1Evicted model.Block) (model.Block, error) {
	// TODO(stage-5): Implement an atomic logical swap.
	return model.Block{}, ErrNotImplemented
}
