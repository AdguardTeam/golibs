// Package redisutil contains common utilities for working with Redis.
//
// # Integration testing
//
// To test with a real Redis database, call the tests with the environment
// variable TEST_REDIS_PORT set to the port of your database.
//
// The tests use the database at index 15, which is chosen because it's the
// largest database index on most instances.  The instance should also have the
// master role.
//
// NOTE:  The database is flushed between tests, so make sure that that database
// is not used for storing important data.
package redisutil

import "time"

// DefaultPort is the default Redis port
const DefaultPort = 6379

// Redis-related limits.
const (
	// MinTTL is the minimum TTL that can be set when setting any TTL.
	MinTTL = 1 * time.Millisecond
)

// Redis commands, parameters, and other constants.
const (
	CmdDEL      = "DEL"
	CmdFCALL    = "FCALL"
	CmdFLUSHDB  = "FLUSHDB"
	CmdFUNCTION = "FUNCTION"
	CmdGET      = "GET"
	CmdROLE     = "ROLE"
	CmdSET      = "SET"
)

// Parameter constants.
const (
	ParamASYNC   = "ASYNC"
	ParamLOAD    = "LOAD"
	ParamNX      = "NX"
	ParamPX      = "PX"
	ParamREPLACE = "REPLACE"
	ParamSYNC    = "SYNC"
)

// Response constants.
const (
	RespOK = "OK"
)

// Error string constants.
const (
	// ErrStrFunctionNotFound is an error message returned by Redis on
	// [CmdFCALL] if the requested function is not found.
	ErrStrFunctionNotFound = "ERR Function not found"
)
