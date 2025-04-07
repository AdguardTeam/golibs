// Package redisutil contains common utilities for working with Redis.
package redisutil

import "time"

// Redis-related limits.
const (
	// MinTTL is the minimum TTL that can be set when setting any TTL.
	MinTTL = 1 * time.Millisecond
)

// Redis commands, parameters, and other constants.
const (
	CmdDEL      = "DEL"
	CmdFCALL    = "FCALL"
	CmdFUNCTION = "FUNCTION"
	CmdGET      = "GET"
	CmdROLE     = "ROLE"
	CmdSET      = "SET"
)

// Parameter constants.
const (
	ParamLoad    = "LOAD"
	ParamMs      = "PX"
	ParamReplace = "REPLACE"
)

// Error string constants.
const (
	// ErrStrFunctionNotFound is an error message returned by Redis on
	// [CmdFCALL] if the requested function is not found.
	ErrStrFunctionNotFound = "ERR Function not found"
)
