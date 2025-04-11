// Package fakeredis contains fake implementations of interfaces from packages
// redis and redisutil.
//
// It is recommended to fill all methods that shouldn't be called with:
//
//	panic("not implemented")
//
// in the body of the test, so that if the method is called the panic backtrace
// points to the method definition in the test.
package fakeredis

import (
	"context"
	"fmt"
	"time"

	"github.com/AdguardTeam/golibs/redisutil"
	"github.com/gomodule/redigo/redis"
)

// Conn is a [redis.Conn] for tests.
type Conn struct {
	OnClose   func() (err error)
	OnDo      func(cmdName string, args ...any) (reply any, err error)
	OnErr     func() (err error)
	OnFlush   func() (err error)
	OnReceive func() (reply any, err error)
	OnSend    func(cmdName string, args ...any) (err error)
}

// type check
var _ redis.Conn = (*Conn)(nil)

// Close implements the [redis.Conn] interface for *Conn.
func (c *Conn) Close() (err error) {
	return c.OnClose()
}

// Do implements the [redis.Conn] interface for *Conn.
func (c *Conn) Do(cmdName string, args ...any) (reply any, err error) {
	return c.OnDo(cmdName, args...)
}

// Err implements the [redis.Conn] interface for *Conn.
func (c *Conn) Err() (err error) {
	return c.OnErr()
}

// Flush implements the [redis.Conn] interface for *Conn.
func (c *Conn) Flush() (err error) {
	return c.OnFlush()
}

// Receive implements the [redis.Conn] interface for *Conn.
func (c *Conn) Receive() (reply any, err error) {
	return c.OnReceive()
}

// Send implements the [redis.Conn] interface for *Conn.
func (c *Conn) Send(cmdName string, args ...any) (err error) {
	return c.OnSend(cmdName, args...)
}

// NewConn returns a new *Conn all methods of which panic.
func NewConn() (c *Conn) {
	return &Conn{
		OnClose: func() (err error) {
			panic(fmt.Errorf("unexpected call to fakeredis.(*Conn).Close()"))
		},
		OnDo: func(cmdName string, args ...any) (reply any, err error) {
			panic(fmt.Errorf("unexpected call to fakeredis.(*Conn).Do(%v, %v)", cmdName, args))
		},
		OnErr: func() (err error) {
			panic(fmt.Errorf("unexpected call to fakeredis.(*Conn).Err()"))
		},
		OnFlush: func() (err error) {
			panic(fmt.Errorf("unexpected call to fakeredis.(*Conn).Flush()"))
		},
		OnReceive: func() (reply any, err error) {
			panic(fmt.Errorf("unexpected call to fakeredis.(*Conn).Receive()"))
		},
		OnSend: func(cmdName string, args ...any) (err error) {
			panic(fmt.Errorf("unexpected call to fakeredis.(*Conn).Send(%v, %v)", cmdName, args))
		},
	}
}

// ConnectionTester is a [redisutil.ConnectionTester] for tests.
type ConnectionTester struct {
	OnTestConnection func(ctx context.Context, c redis.Conn, lastUsed time.Time) (err error)
}

// type check
var _ redisutil.ConnectionTester = (*ConnectionTester)(nil)

// TestConnection implements the [redisutil.ConnectionTester] interface for
// *ConnectionTester.
func (p *ConnectionTester) TestConnection(
	ctx context.Context,
	c redis.Conn,
	lastUsed time.Time,
) (err error) {
	return p.OnTestConnection(ctx, c, lastUsed)
}

// Dialer is a [redisutil.Dialer] for tests.
type Dialer struct {
	OnDialContext func(ctx context.Context) (c redis.Conn, err error)
}

// type check
var _ redisutil.Dialer = (*Dialer)(nil)

// DialContext implements the [redisutil.Dialer] interface for *Dialer.
func (d *Dialer) DialContext(ctx context.Context) (c redis.Conn, err error) {
	return d.OnDialContext(ctx)
}

// Pool is a [redisutil.Pool] for tests.
type Pool struct {
	OnClose func() (err error)
	OnGet   func(ctx context.Context) (c redis.Conn, err error)
}

// type check
var _ redisutil.Pool = (*Pool)(nil)

// Close implements the [redisutil.Pool] interface for *Pool.
func (p *Pool) Close() (err error) {
	return p.OnClose()
}

// Get implements the [redisutil.Pool] interface for *Pool.
func (p *Pool) Get(ctx context.Context) (c redis.Conn, err error) {
	return p.OnGet(ctx)
}
