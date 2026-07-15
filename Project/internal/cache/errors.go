package cache

import "errors"

var (
	// ErrNotImplemented remains available for unfinished layers outside the
	// memory/cache agent's ownership.
	ErrNotImplemented = errors.New("not implemented in scaffold")
	ErrBlockNotFound  = errors.New("block not found")
	ErrVictimDisabled = errors.New("victim cache is disabled")
)
