package external

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBCluster manages primary and replica connections with read/write splitting.
type DBCluster struct {
	primary  *pgxpool.Pool
	replicas []*pgxpool.Pool
	counter  uint64
}

// NewDBCluster creates a database cluster with primary and optional replicas.
func NewDBCluster(ctx context.Context, primaryURL string, replicaURLs []string) (*DBCluster, error) {
	primary, err := newPool(ctx, primaryURL, 100, 150)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to primary: %w", err)
	}

	var replicas []*pgxpool.Pool
	for i, url := range replicaURLs {
		replica, err := newPool(ctx, url, 50, 100)
		if err != nil {
			// Close already created pools
			primary.Close()
			for _, r := range replicas {
				r.Close()
			}
			return nil, fmt.Errorf("failed to connect to replica %d: %w", i+1, err)
		}
		replicas = append(replicas, replica)
	}

	return &DBCluster{
		primary:  primary,
		replicas: replicas,
	}, nil
}

// Primary returns the primary pool for write operations.
func (c *DBCluster) Primary() *pgxpool.Pool {
	return c.primary
}

// Replica returns a replica pool for read operations (round-robin).
// Falls back to primary if no replicas are configured.
func (c *DBCluster) Replica() *pgxpool.Pool {
	if len(c.replicas) == 0 {
		return c.primary
	}
	idx := atomic.AddUint64(&c.counter, 1) % uint64(len(c.replicas))
	return c.replicas[idx]
}

// Close closes all connections in the cluster.
func (c *DBCluster) Close() {
	c.primary.Close()
	for _, r := range c.replicas {
		r.Close()
	}
}

func newPool(ctx context.Context, connString string, minConns, maxConns int32) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MinConns = minConns
	config.MaxConns = maxConns
	config.MaxConnIdleTime = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return pool, nil
}
