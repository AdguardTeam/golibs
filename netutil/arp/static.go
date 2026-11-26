package arp

import (
	"context"
)

// StaticConfig is the configuration structure for [Static].
type StaticConfig struct {
	// Neighbors is a static slice of network neighbors, that will be used to
	// serve [Static.Table] calls.  All elements must not be nil.
	Neighbors []*Neighbor
}

// Static is the [Interface] implementation that uses a static slice as a
// neighbor table.
type Static struct {
	neighbors []*Neighbor
}

// NewStatic returns new properly initialized [*Static].
func NewStatic(c *StaticConfig) (s *Static) {
	if c == nil {
		c = &StaticConfig{}
	}

	var ns []*Neighbor
	if c.Neighbors == nil {
		ns = []*Neighbor{}
	} else {
		ns = c.Neighbors
	}

	return &Static{
		neighbors: ns,
	}
}

// type check
var _ Interface = (*Static)(nil)

// Table implements the [Interface] for Static.  It uses underlying slice deep
// copy as a neighbor table.  err is always nil.
func (s *Static) Table(_ context.Context) (t *Table, err error) {
	n := make([]*Neighbor, 0, len(s.neighbors))

	for _, neighbor := range s.neighbors {
		n = append(n, neighbor.Clone())
	}

	return &Table{Neighbors: n}, nil
}
