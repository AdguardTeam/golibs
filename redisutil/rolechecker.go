package redisutil

import (
	"bytes"
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/validate"
	"github.com/gomodule/redigo/redis"
)

// RoleChecker is a [ConnectionTester] that simplifies checking connection
// roles in a [*DefaultPool].
type RoleChecker struct {
	logger   *slog.Logger
	roleData []byte
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

	if c.Logger == nil {
		c.Logger = slog.Default()
	}

	c.Role = cmp.Or(c.Role, RoleMaster)
	err = c.Role.Validate()
	if err != nil {
		return nil, fmt.Errorf("c.Role: %w", err)
	}

	return &RoleChecker{
		logger:   c.Logger,
		roleData: []byte(c.Role),
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
		return fmt.Errorf("want []byte, got %T(%[1]v)", values[0])
	}

	if !bytes.Equal(roleData, c.roleData) {
		return fmt.Errorf("want role %q, got %q", c.roleData, roleData)
	}

	return nil
}
