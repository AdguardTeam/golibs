package promutil

import "github.com/prometheus/client_golang/prometheus"

// EmptyRegisterer is an implementation of the [prometheus.Registerer] that does
// nothing.
type EmptyRegisterer struct{}

// type check
var _ prometheus.Registerer = EmptyRegisterer{}

// Register implements the [prometheus.Registerer] for EmptyRegistrer.  err is
// always nil.
func (EmptyRegisterer) Register(_ prometheus.Collector) (err error) {
	return nil
}

// MustRegister implements the [prometheus.Registerer] for EmptyRegistrer.
func (EmptyRegisterer) MustRegister(_ ...prometheus.Collector) {}

// Unregister implements the [prometheus.Registerer] for EmptyRegistrer.  ok is
// always true.
func (EmptyRegisterer) Unregister(_ prometheus.Collector) (ok bool) {
	return true
}
