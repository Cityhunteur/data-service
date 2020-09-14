package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

var (
	ErrNotFound = errors.New("record not found")
)

// InTx runs func f with given db isolation level.
func (db *DB) InTx(ctx context.Context, isoLevel pgx.TxIsoLevel, f func(tx pgx.Tx) error) error {
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %v", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel})
	if err != nil {
		return fmt.Errorf("starting transaction: %v", err)
	}

	if err := f(tx); err != nil {
		if errR := tx.Rollback(ctx); errR != nil {
			return fmt.Errorf("rolling back transaction: %v  for err: %v", errR, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}
	return nil
}
