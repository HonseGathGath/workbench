package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS posts (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	slug       TEXT UNIQUE NOT NULL,
	title      TEXT NOT NULL,
	body       TEXT NOT NULL,
	tags       TEXT NOT NULL DEFAULT '',
	published  INTEGER NOT NULL DEFAULT 0,
	post_date  TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
	token      TEXT PRIMARY KEY,
	expires_at TEXT NOT NULL
);
`

func openDB(path string) (*sql.DB, error) {
	// modernc.org/sqlite is a pure-Go driver (no cgo), which keeps
	// deployment to a plain Go buildpack simple.
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
