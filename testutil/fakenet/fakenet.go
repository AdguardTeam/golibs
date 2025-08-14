// Package fakenet contains fake implementations of interfaces from package net
// from the standard library.
//
// It is recommended to fill all methods that shouldn't be called with:
//
//	panic(testutil.UnexpectedCall(arg1, arg2))
//
// See the package example.
package fakenet

import (
	"net"
	"time"
)

// Conn is the [net.Conn] for tests.
type Conn struct {
	OnClose            func() (err error)
	OnLocalAddr        func() (laddr net.Addr)
	OnRead             func(b []byte) (n int, err error)
	OnRemoteAddr       func() (raddr net.Addr)
	OnSetDeadline      func(t time.Time) (err error)
	OnSetReadDeadline  func(t time.Time) (err error)
	OnSetWriteDeadline func(t time.Time) (err error)
	OnWrite            func(b []byte) (n int, err error)
}

// type check
var _ net.Conn = (*Conn)(nil)

// Close implements the [net.Conn] interface for *Conn.
func (c *Conn) Close() (err error) {
	return c.OnClose()
}

// LocalAddr implements the [net.Conn] interface for *Conn.
func (c *Conn) LocalAddr() (laddr net.Addr) {
	return c.OnLocalAddr()
}

// Read implements the [net.Conn] interface for *Conn.
func (c *Conn) Read(b []byte) (n int, err error) {
	return c.OnRead(b)
}

// RemoteAddr implements the [net.Conn] interface for *Conn.
func (c *Conn) RemoteAddr() (raddr net.Addr) {
	return c.OnRemoteAddr()
}

// SetDeadline implements the [net.Conn] interface for *Conn.
func (c *Conn) SetDeadline(t time.Time) (err error) {
	return c.OnSetDeadline(t)
}

// SetReadDeadline implements the [net.Conn] interface for *Conn.
func (c *Conn) SetReadDeadline(t time.Time) (err error) {
	return c.OnSetReadDeadline(t)
}

// SetWriteDeadline implements the [net.Conn] interface for *Conn.
func (c *Conn) SetWriteDeadline(t time.Time) (err error) {
	return c.OnSetWriteDeadline(t)
}

// Write implements the [net.Conn] interface for *Conn.
func (c *Conn) Write(b []byte) (n int, err error) {
	return c.OnWrite(b)
}

// Listener is a [net.Listener] for tests.
type Listener struct {
	OnAccept func() (c net.Conn, err error)
	OnAddr   func() (addr net.Addr)
	OnClose  func() (err error)
}

// type check
var _ net.Listener = (*Listener)(nil)

// Accept implements the [net.Listener] interface for *Listener.
func (l *Listener) Accept() (c net.Conn, err error) {
	return l.OnAccept()
}

// Addr implements the [net.Listener] interface for *Listener.
func (l *Listener) Addr() (addr net.Addr) {
	return l.OnAddr()
}

// Close implements the [net.Listener] interface for *Listener.
func (l *Listener) Close() (err error) {
	return l.OnClose()
}

// PacketConn is the [net.PacketConn] for tests.
type PacketConn struct {
	OnClose            func() (err error)
	OnLocalAddr        func() (laddr net.Addr)
	OnReadFrom         func(b []byte) (n int, addr net.Addr, err error)
	OnSetDeadline      func(t time.Time) (err error)
	OnSetReadDeadline  func(t time.Time) (err error)
	OnSetWriteDeadline func(t time.Time) (err error)
	OnWriteTo          func(b []byte, addr net.Addr) (n int, err error)
}

// type check
var _ net.PacketConn = (*PacketConn)(nil)

// Close implements the [net.PacketConn] interface for *PacketConn.
func (c *PacketConn) Close() (err error) {
	return c.OnClose()
}

// LocalAddr implements the [net.PacketConn] interface for *PacketConn.
func (c *PacketConn) LocalAddr() (laddr net.Addr) {
	return c.OnLocalAddr()
}

// ReadFrom implements the [net.PacketConn] interface for *PacketConn.
func (c *PacketConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	return c.OnReadFrom(b)
}

// SetDeadline implements the [net.PacketConn] interface for *PacketConn.
func (c *PacketConn) SetDeadline(t time.Time) (err error) {
	return c.OnSetDeadline(t)
}

// SetReadDeadline implements the [net.PacketConn] interface for *PacketConn.
func (c *PacketConn) SetReadDeadline(t time.Time) (err error) {
	return c.OnSetReadDeadline(t)
}

// SetWriteDeadline implements the [net.PacketConn] interface for *PacketConn.
func (c *PacketConn) SetWriteDeadline(t time.Time) (err error) {
	return c.OnSetWriteDeadline(t)
}

// WriteTo implements the [net.PacketConn] interface for *PacketConn.
func (c *PacketConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	return c.OnWriteTo(b, addr)
}
