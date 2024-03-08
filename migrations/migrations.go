package migrations

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func RunMigrations(db *sql.DB) error {
	statements := make([]string, 0)

	statements = append(statements, `
	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY,
		url TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		is_private BOOL NOT NULL,
		created TEXT NOT NULL)`)

	statements = append(statements, `
	CREATE TABLE IF NOT EXISTS bookmark_tags (
		bookmark_id INTEGER,
		tag VARCHAR(30))`)

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	for _, statement := range statements {
		_, err := tx.Exec(statement)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
