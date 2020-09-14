package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// NewTestDatabase is a helper to create a test instance of the database.
func NewTestDatabase(tb testing.TB) (*DB, *Config) {
	tb.Helper()

	if testing.Short() {
		tb.Skipf("skipping in short mode")
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("Could not connect to docker: %s", err)
	}

	dbname, dbuser, dbpassword := "postgres", "postgres", "postgres"

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.4-alpine",
		Env: []string{
			"LANG=C",
			"POSTGRES_DB=" + dbname,
			"POSTGRES_USER=" + dbuser,
			"POSTGRES_PASSWORD=" + dbpassword,
		},
	})
	if err != nil {
		tb.Fatalf("Could not start container: %s", err)
	}
	tb.Cleanup(func() {
		if err := pool.Purge(container); err != nil {
			log.Fatalf("Could not purge container: %s", err)
		}
	})

	host := container.Container.NetworkSettings.IPAddress
	if runtime.GOOS == "darwin" {
		host = net.JoinHostPort(container.GetBoundIP("5432/tcp"), container.GetPort("5432/tcp"))
	}

	connURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbuser, dbpassword),
		Host:   host,
		Path:   dbname,
	}
	q := connURL.Query()
	q.Add("sslmode", "disable")
	connURL.RawQuery = q.Encode()

	var dbPool *pgxpool.Pool
	ctx := context.Background()
	if err := pool.Retry(func() error {
		var err error
		dbPool, err = pgxpool.Connect(ctx, connURL.String())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err := runMigrations(connURL.String()); err != nil {
		tb.Fatalf("Could not run database migrations: %s", err)
	}

	db := &DB{Pool: dbPool}
	tb.Cleanup(func() {
		db.Close(context.Background())
	})

	return db, &Config{
		Name:     dbname,
		Host:     container.GetBoundIP("5432/tcp"),
		Port:     container.GetPort("5432/tcp"),
		User:     dbuser,
		Password: dbpassword,
	}
}

func runMigrations(u string) error {
	dir := fmt.Sprintf("file://%s", migrationsDir())
	m, err := migrate.New(dir, u)
	if err != nil {
		return fmt.Errorf("could not create migrate: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrate: %w", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("closing migrate source error: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("closing migrate database error: %w", dbErr)
	}
	return nil
}

func migrationsDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Join(filepath.Dir(filename), "../../migrations")
}
