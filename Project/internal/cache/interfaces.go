package cache

import "victimcacheproject/internal/model"

// Cache defines the logical behavior independently from Akita messaging.
type Cache interface {
	Name() string
	Lookup(req model.Request) LookupResult
	Insert(block model.Block) (evicted *model.Block)
	Invalidate(blockAddress uint64) bool
}

type LookupResult struct {
	Hit   bool
	Block *model.Block
}
