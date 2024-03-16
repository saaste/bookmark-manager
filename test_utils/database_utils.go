package test_utils

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/saaste/bookmark-manager/migrations"
	"github.com/stretchr/testify/assert"
)

func InitTestDatabase(t *testing.T, dbFileName string) *sql.DB {
	_, err := os.Stat("../test_data")
	if os.IsNotExist(err) {
		os.Mkdir("../test_data", os.ModePerm)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("../test_data/%s", dbFileName))
	assert.Nil(t, err)

	err = migrations.RunMigrations(db)
	assert.Nil(t, err)

	_, err = db.Exec("DELETE FROM bookmark_tags")
	assert.Nil(t, err)

	_, err = db.Exec("DELETE FROM bookmarks")
	assert.Nil(t, err)

	return db
}

func DestroyTestDatabase(t *testing.T, db *sql.DB, dbFileName string) {
	defer db.Close()
	err := os.Remove(fmt.Sprintf("../test_data/%s", dbFileName))
	assert.Nil(t, err)
}
