package database

import (
	"context"
	"database/sql"
	"fmt"
)

// QueryRows runs query with args and maps each result row into a T via scan.
// It owns the open/close/iterate/error ritual so callers supply only the SQL and
// the per-row scan. The returned slice is never nil so empty results marshal as
// [] rather than null.
func QueryRows[T any](
	ctx context.Context,
	db *sql.DB,
	query string,
	scan func(*sql.Rows) (T, error),
	args ...any,
) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	out := []T{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}
	return out, nil
}
