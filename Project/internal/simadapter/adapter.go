package simadapter

import (
	"errors"

	"victimcacheproject/internal/system"
)

var ErrNotConnected = errors.New("Akita adapter is not connected yet")

type Adapter struct {
	System *system.System
	// TODO: Add Akita engine, components, ports and connections.
}

func New(sys *system.System) *Adapter { return &Adapter{System: sys} }

func (a *Adapter) Build() error {
	// TODO(stage-4): Instantiate and connect Akita components.
	return ErrNotConnected
}

func (a *Adapter) Run() error {
	// TODO(stage-4): Start the Akita simulation engine.
	return ErrNotConnected
}
