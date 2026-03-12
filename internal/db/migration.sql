CREATE TABLE IF NOT EXISTS notes (
    path         TEXT PRIMARY KEY,
    content_hash INTEGER NOT NULL,
    vector_id    TEXT,
    updated_at   DATETIME NOT NULL
);
