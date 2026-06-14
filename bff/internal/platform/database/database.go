package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func NewClient(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?_pragma=busy_timeout(5000)&journal_mode=WAL"

	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, fmt.Errorf("opening sqlite database %q: %w", path, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging sqlite: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("applying schema: %w", err)
	}

	return db, nil
}
