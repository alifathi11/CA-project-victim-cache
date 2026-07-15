package cpu

import "victimcacheproject/internal/model"

// Generator replaces a full CPU model in the first project stages.
// It emits deterministic memory requests from a benchmark trace.
type Generator struct {
	requests []model.Request
	next     int
}

func NewGenerator(requests []model.Request) *Generator {
	return &Generator{requests: requests}
}

func (g *Generator) HasNext() bool { return g.next < len(g.requests) }

func (g *Generator) Next() (model.Request, bool) {
	if !g.HasNext() {
		return model.Request{}, false
	}
	req := g.requests[g.next]
	g.next++
	return req, true
}
