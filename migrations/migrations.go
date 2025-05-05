package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func RunMigrations(db *sql.DB) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	err = migration1(tx)
	if err != nil {
		return err
	}

	err = migration2(tx)
	if err != nil {
		return err
	}

	latestMigration := getLatestMigration(tx)

	if latestMigration < 3 {
		err = migration3(tx)
		if err != nil {
			return err
		}
	}

	if latestMigration < 4 {
		err = migration4(tx)
		if err != nil {
			return err
		}
	}

	if latestMigration < 5 {
		err = migration5(tx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func migration1(tx *sql.Tx) error {
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

	for _, statement := range statements {
		_, err := tx.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func migration2(tx *sql.Tx) error {
	// Hack to check if this version is already applied
	var name string
	row := tx.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'")
	err := row.Scan(&name)
	if err == nil {
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	statements := make([]string, 0)

	statements = append(statements, `
	CREATE TABLE IF NOT EXISTS migrations (
		version INTEGER PRIMARY KEY,
		created TEXT NOT NULL)`)

	statements = append(statements, `
	ALTER TABLE bookmarks
		ADD COLUMN is_working BOOLEAN NOT NULL DEFAULT true`)

	statements = append(statements, fmt.Sprintf(`
	INSERT INTO migrations (version, created) VALUES (%d, "%s")`, 2, time.Now().UTC().Format(time.RFC3339)))

	for _, statement := range statements {
		_, err := tx.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func migration3(tx *sql.Tx) error {
	statements := make([]string, 0)

	statements = append(statements, `
	ALTER TABLE bookmarks
		ADD COLUMN ignore_check BOOLEAN NOT NULL DEFAULT false`)

	statements = append(statements, fmt.Sprintf(`
	INSERT INTO migrations (version, created) VALUES (%d, "%s")`, 3, time.Now().UTC().Format(time.RFC3339)))

	for _, statement := range statements {
		_, err := tx.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func migration4(tx *sql.Tx) error {
	statements := make([]string, 0)

	statements = append(statements, `
	ALTER TABLE bookmarks
		ADD COLUMN last_status_code INTEGER NOT NULL DEFAULT 0`)

	statements = append(statements, `
	ALTER TABLE bookmarks
		ADD COLUMN error_message TEXT NOT NULL DEFAULT ''`)

	statements = append(statements, fmt.Sprintf(`
	INSERT INTO migrations (version, created) VALUES (%d, "%s")`, 4, time.Now().UTC().Format(time.RFC3339)))

	for _, statement := range statements {
		_, err := tx.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func migration5(tx *sql.Tx) error {
	statements := make([]string, 0)

	statements = append(statements, `
	ALTER TABLE bookmarks
		ADD COLUMN next_check TEXT NULL DEFAULT NULL`)

	statements = append(statements, fmt.Sprintf(`
	INSERT INTO migrations (version, created) VALUES (%d, "%s")`, 5, time.Now().UTC().Format(time.RFC3339)))

	for _, statement := range statements {
		_, err := tx.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func getLatestMigration(tx *sql.Tx) int64 {
	var version int64
	row := tx.QueryRow("SELECT version FROM migrations ORDER BY version DESC LIMIT 1")
	err := row.Scan(&version)
	if err != nil {
		return 0

	}
	return version
}
