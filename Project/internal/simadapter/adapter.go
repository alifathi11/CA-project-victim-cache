// Package simadapter defines the boundary where an Akita event-driven front-end
// can be attached. The default adapter is a complete synchronous reference
// runner, used for correctness tests and reproducible comparisons.
package simadapter

import (
	"fmt"
	"victimcacheproject/internal/model"
	"victimcacheproject/internal/system"
)

type Adapter struct {
	System    *system.System
	Requests  []model.Request
	Responses []model.Response
	built     bool
}

func New(sys *system.System) *Adapter { return &Adapter{System: sys} }
func (a *Adapter) SetRequests(reqs []model.Request) {
	a.Requests = append([]model.Request(nil), reqs...)
}
func (a *Adapter) Build() error {
	if a.System == nil {
		return fmt.Errorf("system is nil")
	}
	if err := a.System.Validate(); err != nil {
		return err
	}
	a.built = true
	return nil
}
func (a *Adapter) Run() error {
	if !a.built {
		return fmt.Errorf("adapter must be built before run")
	}
	a.Responses = a.System.Run(a.Requests)
	return nil
}
