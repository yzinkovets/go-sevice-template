package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"

	"go-service-template/config"
)

type DbConnection struct {
	pool             *pgxpool.Pool
	InsertTimeoutSec int
}

func NewDBConnection(config config.DbConfig) (*DbConnection, error) {
	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Db)

	pgConfig, err := pgxpool.ParseConfig(cs)
	if err != nil {
		logrus.Errorf("Unable to parse connection string: %v\n", err)
		return nil, err
	}

	pgConfig.MaxConns = int32(config.MaxOpenConns)
	pgConfig.MinConns = int32(config.MaxIdleConns)

	// Not sure if we need these settings
	// pgConfig.MaxConnLifetime = time.Hour
	// pgConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		logrus.Errorf("Unable to create connection pool: %v\n", err)
		return nil, err
	}

	return &DbConnection{
		pool:             pool,
		InsertTimeoutSec: config.InsertTimeoutSec,
	}, nil
}

func (c *DbConnection) Close() {
	if c.pool == nil {
		return
	}
	c.pool.Close()
}

// Returns a connection from the pool
// Don't forget to release the connection using conn.Release()
func (c *DbConnection) Conn() (*pgxpool.Conn, error) {
	if c.pool == nil {
		return nil, errors.New("DB connection is not initialized")
	}

	conn, err := c.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
