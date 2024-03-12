package bookmarks

import (
	"database/sql"
	"testing"
	"time"

	"github.com/saaste/bookmark-manager/migrations"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	expected := createBookmark(false)

	actual, err := repo.Create(expected)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
	assert.True(t, actual.ID >= 1)

	bookmarks := getBookmarks(t, db, repo)
	assert.Len(t, bookmarks, 1)

	tags := getTags(t, db)
	assert.Len(t, tags, 2)
	assert.Equal(t, expected.Tags, tags)
}

func TestUpdate(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	expected := createBookmark(false)

	bookmark, err := repo.Create(expected)
	assert.Nil(t, err)

	bookmark.URL = "https://example.org/updated"
	bookmark.Title = "Updated title"
	bookmark.Description = "Updated description"
	bookmark.IsPrivate = false
	bookmark.Tags = []string{"updated1", "updated2"}

	actual, err := repo.Update(bookmark)
	assert.Nil(t, err)
	assert.Equal(t, bookmark, actual)

	bookmarks := getBookmarks(t, db, repo)
	assert.Len(t, bookmarks, 1)
	assert.Equal(t, bookmarks[0].Title, bookmark.Title)

	tags := getTags(t, db)
	assert.Len(t, tags, 2)
	assert.Equal(t, bookmark.Tags, tags)
}

func TestGet(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	expected := createBookmark(false)

	created, err := repo.Create(expected)
	assert.Nil(t, err)

	actual, err := repo.Get(created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, actual)

	_, err = repo.Get(1345623232)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestDelete(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	bookmark := createBookmark(false)

	created, err := repo.Create(bookmark)
	assert.Nil(t, err)

	err = repo.Delete(created.ID)
	assert.Nil(t, err)

	bookmarks := getBookmarks(t, db, repo)
	assert.Len(t, bookmarks, 0)

}

func TestGetAll(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	bookmark1 := createBookmark(false)
	_, err := repo.Create(bookmark1)
	assert.Nil(t, err)

	bookmark2 := createBookmark(false)
	_, err = repo.Create(bookmark2)
	assert.Nil(t, err)

	bookmark3 := createBookmark(true)
	_, err = repo.Create(bookmark3)
	assert.Nil(t, err)

	bookmark4 := createBookmark(true)
	_, err = repo.Create(bookmark4)
	assert.Nil(t, err)

	page1, err := repo.GetAll(true, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, page1.PageCount)
	assert.Len(t, page1.Bookmarks, 2)
	assert.Equal(t, bookmark4.ID, page1.Bookmarks[0].ID)
	assert.Equal(t, bookmark3.ID, page1.Bookmarks[1].ID)

	page2, err := repo.GetAll(true, 2, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, page2.PageCount)
	assert.Len(t, page2.Bookmarks, 2)
	assert.Equal(t, bookmark2.ID, page2.Bookmarks[0].ID)
	assert.Equal(t, bookmark1.ID, page2.Bookmarks[1].ID)

	private, err := repo.GetAll(false, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 1, private.PageCount)
	assert.Len(t, private.Bookmarks, 2)
	assert.Equal(t, bookmark2.ID, private.Bookmarks[0].ID)
	assert.Equal(t, bookmark1.ID, private.Bookmarks[1].ID)
}

func TestGetPrivate(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	bookmark1 := createBookmark(true)
	_, err := repo.Create(bookmark1)
	assert.Nil(t, err)

	bookmark2 := createBookmark(true)
	_, err = repo.Create(bookmark2)
	assert.Nil(t, err)

	bookmark3 := createBookmark(true)
	_, err = repo.Create(bookmark3)
	assert.Nil(t, err)

	bookmark4 := createBookmark(false)
	_, err = repo.Create(bookmark4)
	assert.Nil(t, err)

	page1, err := repo.GetPrivate(1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, page1.PageCount)
	assert.Len(t, page1.Bookmarks, 2)
	assert.Equal(t, bookmark3.ID, page1.Bookmarks[0].ID)
	assert.Equal(t, bookmark2.ID, page1.Bookmarks[1].ID)

	page2, err := repo.GetPrivate(2, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, page2.PageCount)
	assert.Len(t, page2.Bookmarks, 1)
	assert.Equal(t, bookmark1.ID, page2.Bookmarks[0].ID)
}

func TestGetByTags(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	bookmark1 := createBookmark(false)
	_, err := repo.Create(bookmark1)
	assert.Nil(t, err)

	bookmark2 := createBookmark(false)
	_, err = repo.Create(bookmark2)
	assert.Nil(t, err)

	bookmark3 := createBookmark(true)
	_, err = repo.Create(bookmark3)
	assert.Nil(t, err)

	bookmark4 := createBookmark(false)
	bookmark4.Tags = []string{"tag3", "tag4"}
	_, err = repo.Create(bookmark4)
	assert.Nil(t, err)

	public, err := repo.GetByTags(false, []string{"tag1"}, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 1, public.PageCount)
	assert.Len(t, public.Bookmarks, 2)
	assert.Equal(t, bookmark2.ID, public.Bookmarks[0].ID)
	assert.Equal(t, bookmark1.ID, public.Bookmarks[1].ID)

	allPage1, err := repo.GetByTags(true, []string{"tag1"}, 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, allPage1.PageCount)
	assert.Len(t, allPage1.Bookmarks, 2)
	assert.Equal(t, bookmark3.ID, allPage1.Bookmarks[0].ID)
	assert.Equal(t, bookmark2.ID, allPage1.Bookmarks[1].ID)

	allPage2, err := repo.GetByTags(true, []string{"tag1"}, 2, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, allPage2.PageCount)
	assert.Len(t, allPage2.Bookmarks, 1)
	assert.Equal(t, bookmark1.ID, allPage2.Bookmarks[0].ID)
}

func TestGetByKeyword(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	bookmark1 := createBookmark(false)
	_, err := repo.Create(bookmark1)
	assert.Nil(t, err)

	bookmark2 := createBookmark(false)
	_, err = repo.Create(bookmark2)
	assert.Nil(t, err)

	bookmark3 := createBookmark(true)
	_, err = repo.Create(bookmark3)
	assert.Nil(t, err)

	bookmark4 := createBookmark(false)
	bookmark4.Title = "Some other title"
	bookmark4.Description = "Some other description"
	_, err = repo.Create(bookmark4)
	assert.Nil(t, err)

	public, err := repo.GetByKeyword(false, "Test", 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 1, public.PageCount)
	assert.Len(t, public.Bookmarks, 2)
	assert.Equal(t, bookmark2.ID, public.Bookmarks[0].ID)
	assert.Equal(t, bookmark1.ID, public.Bookmarks[1].ID)

	allPage1, err := repo.GetByKeyword(true, "Test", 1, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, allPage1.PageCount)
	assert.Len(t, allPage1.Bookmarks, 2)
	assert.Equal(t, bookmark3.ID, allPage1.Bookmarks[0].ID)
	assert.Equal(t, bookmark2.ID, allPage1.Bookmarks[1].ID)

	allPage2, err := repo.GetByKeyword(true, "Test", 2, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, allPage2.PageCount)
	assert.Len(t, allPage2.Bookmarks, 1)
	assert.Equal(t, bookmark1.ID, allPage2.Bookmarks[0].ID)
}

func TestGetTags(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	bookmark1 := createBookmark(false)
	_, err := repo.Create(bookmark1)
	assert.Nil(t, err)

	bookmark2 := createBookmark(true)
	bookmark2.Tags = []string{"private1", "private2"}
	_, err = repo.Create(bookmark2)
	assert.Nil(t, err)

	public, err := repo.GetTags(false)
	assert.Nil(t, err)
	assert.Len(t, public, 2)
	assert.Equal(t, bookmark1.Tags[0], public[0])
	assert.Equal(t, bookmark1.Tags[1], public[1])

	all, err := repo.GetTags(true)
	assert.Nil(t, err)
	assert.Len(t, all, 4)
	assert.Equal(t, bookmark2.Tags[0], all[0])
	assert.Equal(t, bookmark2.Tags[1], all[1])
	assert.Equal(t, bookmark1.Tags[0], all[2])
	assert.Equal(t, bookmark1.Tags[1], all[3])

}

func TestGetAllWithoutPagination(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	bookmark1 := createBookmark(false)
	_, err := repo.Create(bookmark1)
	assert.Nil(t, err)

	bookmark2 := createBookmark(false)
	_, err = repo.Create(bookmark2)
	assert.Nil(t, err)

	bookmark3 := createBookmark(true)
	_, err = repo.Create(bookmark3)
	assert.Nil(t, err)

	bookmark4 := createBookmark(true)
	_, err = repo.Create(bookmark4)
	assert.Nil(t, err)

	result, err := repo.GetAllWithoutPagination()
	assert.Nil(t, err)
	assert.Len(t, result, 4)
	assert.Equal(t, bookmark4.ID, result[0].ID)
	assert.Equal(t, bookmark3.ID, result[1].ID)
	assert.Equal(t, bookmark2.ID, result[2].ID)
	assert.Equal(t, bookmark1.ID, result[3].ID)
}

func TestGetBrokenBookmarks(t *testing.T) {
	db := initTestDatabase(t)
	defer db.Close()

	repo := NewSqliteRepository(db)

	bookmark1 := createBookmark(false)
	_, err := repo.Create(bookmark1)
	assert.Nil(t, err)

	bookmark2 := createBookmark(false)
	_, err = repo.Create(bookmark2)
	assert.Nil(t, err)

	bookmark3 := createBookmark(true)
	bookmark3.IsWorking = false
	_, err = repo.Create(bookmark3)
	assert.Nil(t, err)

	bookmark4 := createBookmark(true)
	bookmark4.IsWorking = false
	_, err = repo.Create(bookmark4)
	assert.Nil(t, err)

	result, err := repo.GetBrokenBookmarks()
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, bookmark4.ID, result[0].ID)
	assert.Equal(t, bookmark3.ID, result[1].ID)
}

func initTestDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "test_data/test.db")
	assert.Nil(t, err)

	err = migrations.RunMigrations(db)
	assert.Nil(t, err)

	_, err = db.Exec("DELETE FROM bookmark_tags")
	assert.Nil(t, err)

	_, err = db.Exec("DELETE FROM bookmarks")
	assert.Nil(t, err)

	return db
}

// TODO: Destroy database after test run

func createBookmark(isPrivate bool) *Bookmark {
	return &Bookmark{
		URL:         "https://example.org",
		Title:       "Test title",
		Description: "Test description",
		IsPrivate:   isPrivate,
		Created:     removeNanoseconds(time.Now()),
		Tags:        []string{"tag1", "tag2"},
		IsWorking:   true,
	}
}

func getBookmarks(t *testing.T, db *sql.DB, repo *SqliteRepository) []*Bookmark {
	bookmarks := make([]*Bookmark, 0)

	rows, err := db.Query("SELECT * FROM bookmarks ORDER BY created ASC")
	assert.Nil(t, err)

	for rows.Next() {
		bm, err := repo.scanBookmarkRow(rows)
		assert.Nil(t, err)
		bookmarks = append(bookmarks, bm)
	}
	return bookmarks
}

func getTags(t *testing.T, db *sql.DB) []string {
	tags := make([]string, 0)

	rows, err := db.Query("SELECT tag FROM bookmark_tags ORDER BY tag ASC")
	assert.Nil(t, err)

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		assert.Nil(t, err)
		tags = append(tags, tag)
	}

	return tags
}

func removeNanoseconds(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), 0, dt.Location()).UTC()
}
