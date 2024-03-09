package migrations

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunMigrations(t *testing.T) {
	db, err := sql.Open("sqlite3", "test_data/test.db")
	assert.Nil(t, err)

	err = RunMigrations(db)
	assert.Nil(t, err)

	_, err = db.Query("SELECT * FROM bookmarks")
	assert.Nil(t, err)

	_, err = db.Query("SELECT * FROM bookmark_tags")
	assert.Nil(t, err)
}
