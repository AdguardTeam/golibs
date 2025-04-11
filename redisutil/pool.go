package redisutil

import (
	"cmp"
	"context"
	"log/slog"
	"time"

	"github.com/AdguardTeam/golibs/validate"
	"github.com/gomodule/redigo/redis"
)

// Pool is a pool of Redis connections.
type Pool interface {
	// Get returns a connection from the pool.  If the context expires before
	// the connection is complete, an error must be returned; any expiration on
	// the context must not affect the returned connection.  If the function
	// completes without error, then the application should close the returned
	// connection.
	Get(ctx context.Context) (c redis.Conn, err error)
}

// DefaultPool is a warpper around [redis.DefaultPool] with metrics and
// additional options.
type DefaultPool struct {
	metrics PoolMetrics
	pool    *redis.Pool
}

// DefaultPoolConfig is the configuration for the default Redis pool.
type DefaultPoolConfig struct {
	// Logger is used to log the operation of the Redis pool.  If nil,
	// [slog.Default] is used.
	Logger *slog.Logger

	// Dialer is used to create and configure connections.  It must not be nil.
	Dialer Dialer

	// ConnectionTester checks the health of an idle connection before the
	// connection is used again.  If nil, no checks are performed.
	ConnectionTester ConnectionTester

	// Metrics is used for the collection of the Redis pool statistics.  If nil,
	// [EmptyPoolMetrics] is used.
	Metrics PoolMetrics

	// MaxActive is the maximum number of connections allocated by the Redis
	// connection-pool at a given time.  When zero, there is no limit on the
	// number of connections in the pool.
	MaxActive int

	// MaxIdle is the maximum number of idle connections in the pool.  When
	// zero, there is no limit.
	MaxIdle int

	// IdleTimeout is the time after remaining, idle connection will be closed.
	// If the value is zero, then idle connections are not closed.  Applications
	// should set the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// MaxConnLifetime is the total duration of any connection's lifetime.  If
	// the value is zero then the pool does not close connections based on age.
	MaxConnLifetime time.Duration

	// Wait, if true, makes he pool wait for a connection once the
	// [PoolConfig.MaxActive] limit is reached.
	Wait bool
}

// ConnectionTester checks the health of an idle connection before the
// connection is used again by the application.
type ConnectionTester interface {
	// TestConnection returns an error if c is not healthy.  lastUsed is the
	// time when the connection was returned to the pool.  c must not be nil.
	TestConnection(ctx context.Context, c redis.Conn, lastUsed time.Time) (err error)
}

// NewPool returns a new properly initialized *DefaultPool.  c should not be nil
// and should be valid.
func NewPool(c *DefaultPoolConfig) (p *DefaultPool, err error) {
	err = validate.NotNil("c", c)
	if err != nil {
		return nil, err
	}

	err = validate.NotNilInterface("c.Dialer", c.Dialer)
	if err != nil {
		return nil, err
	}

	c.Logger = cmp.Or(c.Logger, slog.Default())
	c.Metrics = cmp.Or[PoolMetrics](c.Metrics, EmptyPoolMetrics{})

	var checkConn func(context.Context, redis.Conn, time.Time) error
	if c.ConnectionTester != nil {
		checkConn = c.ConnectionTester.TestConnection
	}

	return &DefaultPool{
		metrics: c.Metrics,
		pool: &redis.Pool{
			DialContext:         c.Dialer.DialContext,
			TestOnBorrowContext: checkConn,
			MaxIdle:             c.MaxIdle,
			MaxActive:           c.MaxActive,
			IdleTimeout:         c.IdleTimeout,
			Wait:                c.Wait,
			MaxConnLifetime:     c.MaxConnLifetime,
		},
	}, nil
}

// Get returns a connection from the pool and also updates the pool metrics.  If
// the context expires before the connection is complete, an error is returned;
// any expiration on the context will not affect the returned connection.  If
// the function completes without error, then the application should close the
// returned connection.
func (p *DefaultPool) Get(ctx context.Context) (c redis.Conn, err error) {
	c, err = p.pool.GetContext(ctx)

	stats := p.pool.Stats()
	p.metrics.Update(ctx, stats, err)

	if err != nil {
		return nil, err
	}

	return c, nil
}

// PoolMetrics is an interface that is used for the collection of the Redis pool
// statistics.
type PoolMetrics interface {
	// Update updates the total number of active connections and increments the
	// total number of errors if necessary.
	Update(ctx context.Context, s redis.PoolStats, err error)
}

// EmptyPoolMetrics is the implementation of the [PoolMetrics] interface that
// does nothing.
type EmptyPoolMetrics struct{}

// type check
var _ PoolMetrics = EmptyPoolMetrics{}

// Update implements the [PoolMetrics] interface for EmptyPoolMetrics.
func (EmptyPoolMetrics) Update(_ context.Context, _ redis.PoolStats, _ error) {}
