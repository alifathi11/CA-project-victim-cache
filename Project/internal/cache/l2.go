package cache

import (
	"victimcacheproject/internal/config"
	"victimcacheproject/internal/model"
)

type L2 struct {
	cfg config.Config
	// TODO(stage-3): Add set-associative storage and replacement metadata.
}

func NewL2(cfg config.Config) *L2                   { return &L2{cfg: cfg} }
func (c *L2) Name() string                          { return "L2" }
func (c *L2) Lookup(req model.Request) LookupResult { return LookupResult{} }
func (c *L2) Insert(block model.Block) *model.Block { return nil }
func (c *L2) Invalidate(blockAddress uint64) bool   { return false }
