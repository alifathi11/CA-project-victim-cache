package memory

import "victimcacheproject/internal/model"

type MainMemory struct {
	// TODO(stage-3): Add optional backing storage and latency configuration.
}

func NewMainMemory() *MainMemory { return &MainMemory{} }

func (m *MainMemory) ReadBlock(blockAddress uint64) model.Block {
	// TODO(stage-3): Return a valid block aligned to block size.
	return model.Block{Address: blockAddress, Valid: true}
}

func (m *MainMemory) WriteBlock(block model.Block) {
	// TODO(stage-3): Persist dirty blocks only if functional data is modeled.
}
