// Package fakeservice contains fake implementations of interfaces from package
// service.
//
// It is recommended to fill all methods that shouldn't be called with:
//
//	panic(testutil.UnexpectedCall(arg1, arg2))
package fakeservice

import (
	"context"

	"github.com/AdguardTeam/golibs/service"
)

// Service is the [service.Interface] for tests.
type Service struct {
	OnStart    func(ctx context.Context) (err error)
	OnShutdown func(ctx context.Context) (err error)
}

// type check
var _ service.Interface = (*Service)(nil)

// Start implements the [service.Interface] interface for *Service.
func (s *Service) Start(ctx context.Context) (err error) {
	return s.OnStart(ctx)
}

// Shutdown implements the [service.Interface] interface for *Service.
func (s *Service) Shutdown(ctx context.Context) (err error) {
	return s.OnShutdown(ctx)
}

// Refresher is the [service.Refresher] for tests.
type Refresher struct {
	OnRefresh func(ctx context.Context) (err error)
}

// type check
var _ service.Refresher = (*Refresher)(nil)

// Refresh implements the [service.Refresher] interface for *Refresher.
func (r *Refresher) Refresh(ctx context.Context) (err error) {
	return r.OnRefresh(ctx)
}
