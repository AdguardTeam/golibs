package redisutil

import (
	"bytes"
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/validate"
	"github.com/gomodule/redigo/redis"
)

// Role is a valid Redis role.
type Role string

// Valid Redis roles and their string representations.
const (
	RoleStringMaster   = "master"
	RoleStringSentinel = "sentinel"
	RoleStringSlave    = "slave"

	RoleMaster   Role = RoleStringMaster
	RoleSentinel Role = RoleStringSentinel
	RoleSlave    Role = RoleStringSlave
)

// NewRole converts s into a role.
func NewRole(s string) (r Role, err error) {
	r = Role(s)
	err = r.Validate()
	if err != nil {
		return "", err
	}

	return r, nil
}

// type check
var _ validate.Interface = RoleMaster

// Validate implements the [validate.Interface] interface for Role.
func (r Role) Validate() (err error) {
	switch r {
	case RoleMaster, RoleSentinel, RoleSlave:
		return nil
	default:
		return fmt.Errorf("%w: %q", errors.ErrBadEnumValue, string(r))
	}
}

// RoleChecker is a [ConnectionTester] that simplifies checking connection
// roles in a [*Pool].
type RoleChecker struct {
	logger           *slog.Logger
	requiredRoleData []byte
}

// RoleCheckerConfig is the configuration structure for a [*RoleChecker].
type RoleCheckerConfig struct {
	// Logger is used to log the operation of the checker.  If nil,
	// [slog.Default] is used.
	Logger *slog.Logger

	// Role is the required role.  If empty, [RoleMaster] is used; otherwise it
	// must be a valid Role.
	Role Role
}

// NewRoleChecker returns a new properly initialized *RoleChecker.  If c is nil,
// the defaults are used.
func NewRoleChecker(c *RoleCheckerConfig) (rc *RoleChecker, err error) {
	if c == nil {
		c = &RoleCheckerConfig{}
	}

	c.Logger = cmp.Or(c.Logger, slog.Default())
	c.Role = cmp.Or(c.Role, RoleMaster)

	err = c.Role.Validate()
	if err != nil {
		return nil, fmt.Errorf("c.Role: %w", err)
	}

	return &RoleChecker{
		logger:           c.Logger,
		requiredRoleData: []byte(c.Role),
	}, nil
}

// type check
var _ ConnectionTester = (*RoleChecker)(nil)

// TestConnection implements the [ConnectionTester] interface for *RoleChecker.
func (c *RoleChecker) TestConnection(
	ctx context.Context,
	conn redis.Conn,
	lastUsed time.Time,
) (err error) {
	defer func() {
		if err != nil {
			c.logger.ErrorContext(ctx, "error checking conn role", slogutil.KeyError, err)
			err = fmt.Errorf("testing conn: %w", err)
		} else {
			c.logger.Log(ctx, slogutil.LevelTrace, "redis check successful")
		}
	}()

	values, err := redis.Values(conn.Do(CmdROLE))
	if err != nil {
		return fmt.Errorf("sending command %q: %w", CmdROLE, err)
	}

	err = validate.NotEmptySlice("role_values", values)
	if err != nil {
		return fmt.Errorf("parsing command %q output: %w", CmdROLE, err)
	}

	roleData, ok := values[0].([]byte)
	if !ok {
		return fmt.Errorf("want type string, got %T(%[1]v)", values[0])
	}

	if !bytes.Equal(roleData, c.requiredRoleData) {
		return fmt.Errorf("want role %q, got %q", c.requiredRoleData, roleData)
	}

	return nil
}
