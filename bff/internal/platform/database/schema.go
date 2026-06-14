package database

// schema is the SQLite DDL applied on connection.
const schema = `
	CREATE TABLE IF NOT EXISTS transactions (
		id         TEXT PRIMARY KEY,
		user_id    TEXT NOT NULL,
		date       TEXT NOT NULL,
		merchant   TEXT NOT NULL,
		category   TEXT NOT NULL,
		amount     INTEGER NOT NULL,
		direction  TEXT NOT NULL CHECK (direction IN ('INFLOW','OUTFLOW')),
		notes      TEXT NOT NULL DEFAULT '',
		split      INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_tx_user_date ON transactions (user_id, date);`
