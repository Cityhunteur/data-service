package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataDB_InsertData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testDB, _ := NewTestDatabase(t)

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			"valid",
			"name1",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DataDB{
				db: testDB,
			}
			d, err := db.InsertData(ctx, tt.title)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, d.ID)
				assert.NotEmpty(t, d.Timestamp)
			}
		})
	}
}

func TestDataDB_GetData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testDB, _ := NewTestDatabase(t)

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			"valid",
			"name1",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DataDB{
				db: testDB,
			}
			_, _ = db.InsertData(ctx, tt.title)
			d, err := db.GetData(ctx, tt.title)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.title, d.Title)
			}
		})
	}
}

func TestDataDB_GetData_NotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testDB, _ := NewTestDatabase(t)

	db := &DataDB{
		db: testDB,
	}
	_, err := db.GetData(ctx, "invalid")
	assert.ErrorIs(t, err, ErrNotFound)
}
