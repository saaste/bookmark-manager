package migrations

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunMigrations(t *testing.T) {
	_, err := os.Stat("../test_data")
	if os.IsNotExist(err) {
		os.Mkdir("../test_data", os.ModePerm)
	}

	db, err := sql.Open("sqlite3", "../test_data/migration_test.db")
	assert.Nil(t, err)

	defer destroyTestDatabase(t, db)

	err = RunMigrations(db)
	assert.Nil(t, err)

	_, err = db.Query("SELECT * FROM bookmarks")
	assert.Nil(t, err)

	_, err = db.Query("SELECT * FROM bookmark_tags")
	assert.Nil(t, err)
}

func destroyTestDatabase(t *testing.T, db *sql.DB) {
	defer db.Close()
	err := os.Remove("../test_data/migration_test.db")
	assert.Nil(t, err)
}
