package database

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
)

// newTestDB returns a freshly-initialized sqlite database backed by a temp file.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := NewClient(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

type numRow struct {
	n     int
	label string
}

func scanNumRow(r *sql.Rows) (numRow, error) {
	var x numRow
	return x, r.Scan(&x.n, &x.label)
}

func TestQueryRows_MapsRows(t *testing.T) {
	db := newTestDB(t)
	if _, err := db.Exec(`CREATE TABLE nums (n INTEGER, label TEXT)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO nums VALUES (2,'b'),(1,'a')`); err != nil {
		t.Fatalf("insert: %v", err)
	}

	got, err := QueryRows(context.Background(), db, `SELECT n, label FROM nums ORDER BY n`, scanNumRow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != (numRow{1, "a"}) || got[1] != (numRow{2, "b"}) {
		t.Errorf("unexpected rows: %+v", got)
	}
}

func TestQueryRows_PassesArgs(t *testing.T) {
	db := newTestDB(t)
	if _, err := db.Exec(`CREATE TABLE nums (n INTEGER, label TEXT)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO nums VALUES (1,'a'),(2,'b'),(3,'c')`); err != nil {
		t.Fatalf("insert: %v", err)
	}

	got, err := QueryRows(context.Background(), db, `SELECT n, label FROM nums WHERE n >= ? ORDER BY n`, scanNumRow, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].n != 2 {
		t.Errorf("args not applied: %+v", got)
	}
}

func TestQueryRows_EmptyIsNonNil(t *testing.T) {
	db := newTestDB(t)
	if _, err := db.Exec(`CREATE TABLE nums (n INTEGER, label TEXT)`); err != nil {
		t.Fatalf("create table: %v", err)
	}

	got, err := QueryRows(context.Background(), db, `SELECT n, label FROM nums`, scanNumRow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil slice so it marshals as [] not null")
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d", len(got))
	}
}

func TestQueryRows_QueryError(t *testing.T) {
	db := newTestDB(t)

	_, err := QueryRows(context.Background(), db, `SELECT * FROM does_not_exist`, scanNumRow)
	if err == nil {
		t.Fatal("expected error for bad query")
	}
}

func TestQueryRows_ScanError(t *testing.T) {
	db := newTestDB(t)
	if _, err := db.Exec(`CREATE TABLE nums (n INTEGER, label TEXT)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO nums VALUES (1,'a')`); err != nil {
		t.Fatalf("insert: %v", err)
	}

	// Select two columns but bind only one in scan -> Scan returns an error.
	_, err := QueryRows(context.Background(), db, `SELECT n, label FROM nums`,
		func(r *sql.Rows) (int, error) {
			var n int
			return n, r.Scan(&n)
		})
	if err == nil {
		t.Fatal("expected scan error for column/dest mismatch")
	}
}
