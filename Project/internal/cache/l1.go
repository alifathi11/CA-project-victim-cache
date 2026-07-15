package cache

import (
	"victimcacheproject/internal/config"
	"victimcacheproject/internal/model"
)

type L1 struct {
	cfg config.Config
	// TODO(stage-2): Add sets/lines and replacement metadata.
}

func NewL1(cfg config.Config) *L1 { return &L1{cfg: cfg} }
func (c *L1) Name() string        { return "L1" }

func (c *L1) Lookup(req model.Request) LookupResult {
	// TODO(stage-2): Decode block address, index and tag.
	return LookupResult{}
}

func (c *L1) Insert(block model.Block) *model.Block {
	// TODO(stage-2): Insert into direct-mapped line and return eviction.
	return nil
}

func (c *L1) Invalidate(blockAddress uint64) bool {
	// TODO(stage-2): Locate and invalidate matching block.
	return false
}
