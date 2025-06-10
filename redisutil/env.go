package redisutil

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/timeutil"
	"github.com/caarlos0/env/v11"
)

// Environment represents the Redis pool configuration that is kept in the
// environment.
type Environment struct {
	// Host is the value of the REDIS_HOST environment variable, which is used
	// together with Port to create the address to connect to.  The default
	// value is "localhost".
	Host string `env:"REDIS_HOST" envDefault:"localhost"`

	// Network is the value of the REDIS_NETWORK environment variable, which
	// shows what kind of IP protocol version to use:
	//   - "ip" means both;
	//   - "ip4" means IPv4 only;
	//   - "ip6" means IPv6 only.
	// All other values are invalid.  The default value is "ip4".
	Network string `env:"REDIS_NETWORK" envDefault:"ip4"`

	// IdleTimeout is the value of the REDIS_IDLE_TIMEOUT environment variable,
	// which is used to set the idle timeout for connections in a pool; see
	// [DefaultPoolConfig.IdleTimeout].  The duration should be in the
	// [time.Duration] format.  The default value is "5m".
	IdleTimeout timeutil.Duration `env:"REDIS_IDLE_TIMEOUT" envDefault:"5m"`

	// MaxConnLifetime is the value of the REDIS_MAX_CONN_LIFETIME environment
	// variable, which is used to set the maximum total duration of connections
	// in a pool; see [DefaultPoolConfig.MaxConnLifetime].  The duration should
	// be in the [time.Duration] format.  The default value is "0s", which means
	// that the lifetime is not limited.
	MaxConnLifetime timeutil.Duration `env:"REDIS_MAX_CONN_LIFETIME" envDefault:"0s"`

	// MaxActive is the value of the REDIS_MAX_ACTIVE environment variable,
	// which is used to set the maximum number of connections in a pool; see
	// [DefaultPoolConfig.MaxActive].  The default value is 100.
	MaxActive int `env:"REDIS_MAX_ACTIVE" envDefault:"100"`

	// MaxIdle is the value of the REDIS_MAX_IDLE environment variable, which is
	// used to set the maximum number of idle connections in a pool; see
	// [DefaultPoolConfig.MaxIdle].  The default value is 100.
	MaxIdle int `env:"REDIS_MAX_IDLE" envDefault:"100"`

	// Port is the value of the REDIS_PORT environment variable, which is used
	// together with HOST to create the address to connect to.  The default
	// value is 6379.
	Port uint16 `env:"REDIS_PORT" envDefault:"6379"`

	// DBIndex is the value of the REDIS_DB environment variable, denoting the
	// index of Redis database to use.  The default value is 0.
	DBIndex uint8 `env:"REDIS_DB" envDefault:"0"`

	// Wait is the value of the REDIS_WAIT environment variable, which selects
	// if the pool must wait for a connection once the MaxActive limit is
	// reached; see [DefaultPoolConfig.Wait].  The default is to wait.
	Wait bool `env:"REDIS_WAIT" envDefault:"1"`
}

// Prefix values for the loggers of entities initialized by
// [NewPoolFromEnvironment].
const (
	LogPrefixValuePool        = "redis_pool"
	LogPrefixValueRoleChecker = "redis_role_checker"
)

// NewPoolFromEnvironment creates a new pool based on the environment.  If
// baseLogger is nil, [slog.Default] is used.  If mtrc is nil,
// [EmptyPoolMetrics] is used.
//
// See [Environment] and its fields for more information about the environment
// variables that this function uses.
//
// TODO(a.garipov):  Find ways of testing.
func NewPoolFromEnvironment(
	ctx context.Context,
	baseLogger *slog.Logger,
	mtrc PoolMetrics,
) (p *DefaultPool, err error) {
	envs := &Environment{}
	err = env.Parse(envs)
	if err != nil {
		return nil, fmt.Errorf("parsing environment: %w", err)
	}

	dialer, err := NewDefaultDialer(&DefaultDialerConfig{
		Addr: &netutil.HostPort{
			Host: envs.Host,
			Port: envs.Port,
		},
		Network: envs.Network,
		DBIndex: envs.DBIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("creating dialer: %w", err)
	}

	if baseLogger == nil {
		baseLogger = slog.Default()
	}

	connTester, err := NewRoleChecker(&RoleCheckerConfig{
		Logger: baseLogger.With(slogutil.KeyPrefix, LogPrefixValueRoleChecker),
	})
	if err != nil {
		return nil, fmt.Errorf("creating role checker: %w", err)
	}

	if mtrc == nil {
		mtrc = EmptyPoolMetrics{}
	}

	p, err = NewDefaultPool(&DefaultPoolConfig{
		Logger:           baseLogger.With(slogutil.KeyPrefix, LogPrefixValuePool),
		ConnectionTester: connTester,
		Dialer:           dialer,
		Metrics:          mtrc,
		IdleTimeout:      time.Duration(envs.IdleTimeout),
		MaxConnLifetime:  time.Duration(envs.MaxConnLifetime),
		MaxActive:        envs.MaxActive,
		MaxIdle:          envs.MaxIdle,
		Wait:             envs.Wait,
	})
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	return p, nil
}
