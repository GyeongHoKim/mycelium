package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gyeonghokim/mycelium/internal/db"
	"github.com/gyeonghokim/mycelium/internal/db/models"
)

func TestOpen_Migration(t *testing.T) {
	t.Parallel()
	db, err := db.Open(":memory:")
	require.NoError(t, err)
	defer db.Close()

	note := &models.Note{
		Path:        "README.md",
		ContentHash: 1234567890,
		VectorID:    "1234567890",
		UpdatedAt:   time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC),
	}
	_, err = db.Upsert(context.Background(), note)
	require.NoError(t, err)

	got, err := db.Get(context.Background(), "README.md")
	require.NoError(t, err)
	require.Equal(t, note, got)
}
