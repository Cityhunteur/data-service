package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/google/uuid"

	"github.com/cityhunteur/data-service/internal/model"
)

type DataDB struct {
	db *DB
}

func NewDataDB(db *DB) *DataDB {
	return &DataDB{
		db: db,
	}
}

// InsertData creates a new data record with the given title.
func (db *DataDB) InsertData(ctx context.Context, title string) (*model.Data, error) {
	var d model.Data
	id := uuid.New().String()
	err := db.db.InTx(ctx, pgx.Serializable, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, `
			INSERT INTO
				Data
				(id, title)
			VALUES
				($1, $2)
			RETURNING id, title, timestamp
		`, id, title)

		if err := row.Scan(&d.ID, &d.Title, &d.Timestamp); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("insert data: %w", err)
	}
	return &d, nil
}

// GetData returns a data record with the matching title, if any.
func (db *DataDB) GetData(ctx context.Context, title string) (*model.Data, error) {
	var d model.Data
	err := db.db.InTx(ctx, pgx.ReadCommitted, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, `
			SELECT
				id, title, timestamp
			FROM
				Data
			WHERE 
				title = $1
			ORDER BY
				timestamp
		`, title)

		if err := row.Scan(&d.ID, &d.Title, &d.Timestamp); err != nil {
			return err
		}
		return nil
	})
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get data: %w", err)
	}
	return &d, nil
}
