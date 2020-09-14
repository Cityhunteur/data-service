package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewFromEnv(ctx context.Context, config *Config) (*DB, error) {
	connStr := dbConnectionString(config)

	var err error
	var pool *pgxpool.Pool
	err = backoff.Retry(func() error {
		pool, err = pgxpool.Connect(ctx, connStr)
		if err != nil {
			return err
		}
		return nil
	}, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %v", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close(_ context.Context) {
	db.Pool.Close()
}

func dbConnectionString(config *Config) string {
	vals := dbValues(config)
	var p []string
	for k, v := range vals {
		p = append(p, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(p, " ")
}

func setIfNotEmpty(m map[string]string, key, val string) {
	if val != "" {
		m[key] = val
	}
}

func dbValues(config *Config) map[string]string {
	p := map[string]string{}
	setIfNotEmpty(p, "dbname", config.Name)
	setIfNotEmpty(p, "user", config.User)
	setIfNotEmpty(p, "host", config.Host)
	setIfNotEmpty(p, "port", config.Port)
	setIfNotEmpty(p, "password", config.Password)
	return p
}
