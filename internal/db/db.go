// Package db provides SQLite storage for note metadata (path, content_hash, updated_at).
package db

import (
	"context"
	"database/sql"

	_ "embed"

	_ "modernc.org/sqlite" // Pure Go SQLite driver (no cgo)

	"github.com/gyeonghokim/mycelium/internal/db/models"
)

//go:embed migration.sql
var migrationSQL string

type DB struct {
	conn *sql.DB
}

func Open(path string) (*DB, error) {
	// 주의: 이때 DB 연결 안함 sql.Open은 추상화를 위한 준비만 하고 실제 연결은 Ping으로 확인.
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err = db.PingContext(context.Background()); err != nil {
		return nil, err
	}
	if err = migrate(db); err != nil {
		return nil, err
	}
	return &DB{conn: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.ExecContext(context.Background(), migrationSQL)
	return err
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) Upsert(ctx context.Context, note *models.Note) (*models.Note, error) {
	_, err := d.conn.ExecContext(
		ctx,
		"INSERT INTO notes (path, content_hash, vector_id, updated_at) VALUES (?, ?, ?, ?) ON CONFLICT(path) DO UPDATE SET content_hash = ?, vector_id = ?, updated_at = ?",
		note.Path,
		note.ContentHash,
		note.VectorID,
		note.UpdatedAt.UTC(),
		note.ContentHash,
		note.VectorID,
		note.UpdatedAt.UTC(),
	)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (d *DB) Get(ctx context.Context, path string) (*models.Note, error) {
	var note models.Note
	err := d.conn.QueryRowContext(ctx, "SELECT path, content_hash, vector_id, updated_at from notes where path = ?", path).
		Scan(&note.Path, &note.ContentHash, &note.VectorID, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (d *DB) Delete(ctx context.Context, path string) error {
	_, err := d.conn.ExecContext(ctx, "DELETE FROM notes WHERE path = ?", path)
	return err
}

func (d *DB) All(ctx context.Context) ([]*models.Note, error) {
	var notes []*models.Note
	rows, err := d.conn.QueryContext(ctx, "SELECT path, content_hash, vector_id, updated_at from notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var note models.Note
		if err = rows.Scan(&note.Path, &note.ContentHash, &note.VectorID, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}
	// rows.Next() 실패해서 false 반환하고 나서 발생하는 에러 잡음.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}
