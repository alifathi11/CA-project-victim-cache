package model

// Block contains only metadata in the scaffold.
type Block struct {
	Address  uint64
	Tag      uint64
	Valid    bool
	Dirty    bool
	LastUsed uint64
	Inserted uint64
}
