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

import "github.com/gomodule/redigo/redis"

// Conn is a [redis.Conn] for tests.
type Conn struct {
	OnClose   func() (err error)
	OnErr     func() (err error)
	OnDo      func(cmdName string, args ...any) (reply any, err error)
	OnSend    func(cmdName string, args ...any) (err error)
	OnFlush   func() (err error)
	OnReceive func() (reply any, err error)
}

// type check
var _ redis.Conn = (*Conn)(nil)

// Close implements the [redis.Conn] interface for *Conn.
func (c *Conn) Close() (err error) {
	return c.OnClose()
}

// Err implements the [redis.Conn] interface for *Conn.
func (c *Conn) Err() (err error) {
	return c.OnErr()
}

// Do implements the [redis.Conn] interface for *Conn.
func (c *Conn) Do(commandName string, args ...any) (reply any, err error) {
	return c.OnDo(commandName, args...)
}

// Send implements the [redis.Conn] interface for *Conn.
func (c *Conn) Send(commandName string, args ...any) (err error) {
	return c.OnSend(commandName, args...)
}

// Flush implements the [redis.Conn] interface for *Conn.
func (c *Conn) Flush() (err error) {
	return c.OnFlush()
}

// Receive implements the [redis.Conn] interface for *Conn.
func (c *Conn) Receive() (reply any, err error) {
	return c.OnReceive()
}
