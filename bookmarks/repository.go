package bookmarks

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("bookmark not found")
)

type Repository interface {
	Create(bookmark *Bookmark) (*Bookmark, error)
	Update(bookmark *Bookmark) (*Bookmark, error)
	Get(id int64) (*Bookmark, error)
	Delete(id int64) error
	GetAll(showPrivate bool, page, pageSize int) (*BookmarkResult, error)
	GetPrivate(page, pageSize int) (*BookmarkResult, error)
	GetByTags(showPrivate bool, containsTags []string, page, pageSize int) (*BookmarkResult, error)
	GetByKeyword(showPrivate bool, q string, page, pageSize int) (*BookmarkResult, error)
	GetTags(showPrivate bool) ([]string, error)
	GetCheckable() ([]*Bookmark, error)
	GetBrokenBookmarks() ([]*Bookmark, error)
	BrokenBookmarksExist() (bool, error)
}

type SqliteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(db *sql.DB) *SqliteRepository {
	return &SqliteRepository{
		db: db,
	}
}

func (r *SqliteRepository) Create(bookmark *Bookmark) (*Bookmark, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}

	defer tx.Rollback()

	res, err := tx.Exec(
		"INSERT INTO bookmarks (url, title, description, is_private, created, is_working, ignore_check) VALUES (?, ?, ?, ?, ?, ?, ?)",
		bookmark.URL, bookmark.Title, bookmark.Description, bookmark.IsPrivate, bookmark.Created.UTC().Format(time.RFC3339), bookmark.IsWorking, bookmark.IgnoreCheck)
	if err != nil {
		return nil, fmt.Errorf("sql exec failed: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("fetching last inserted id failed: %w", err)
	}

	bookmark.ID = id
	err = r.setBookmarkTags(id, bookmark.Tags, tx)
	if err != nil {
		return nil, fmt.Errorf("setting bookmark tags failed: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}
	return bookmark, nil
}

func (r *SqliteRepository) Update(bookmark *Bookmark) (*Bookmark, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction begin failed: %w", err)
	}

	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE bookmarks SET url = ?, title = ?, description = ?, is_private = ?, is_working = ?, ignore_check = ?, last_status_code = ?, error_message = ? WHERE id = ?",
		bookmark.URL, bookmark.Title, bookmark.Description, bookmark.IsPrivate, bookmark.IsWorking, bookmark.IgnoreCheck, bookmark.LastStatusCode, bookmark.ErrorMessage, bookmark.ID)
	if err != nil {
		return nil, fmt.Errorf("sql exec failed: %w", err)
	}

	err = r.setBookmarkTags(bookmark.ID, bookmark.Tags, tx)
	if err != nil {
		return nil, fmt.Errorf("setting bookmark tags failed: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return bookmark, nil
}

func (r *SqliteRepository) Get(id int64) (*Bookmark, error) {
	row := r.db.QueryRow("SELECT * FROM bookmarks WHERE id = ? LIMIT 1", id)

	var bm Bookmark
	var created string
	if err := row.Scan(&bm.ID, &bm.URL, &bm.Title, &bm.Description, &bm.IsPrivate, &created, &bm.IsWorking, &bm.IgnoreCheck, &bm.LastStatusCode, &bm.ErrorMessage); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("db row scan failed: %w", err)
	}

	parsedDate, err := time.Parse(time.RFC3339, created)
	if err != nil {
		return nil, fmt.Errorf("time parsing failed: %w", err)
	}

	tagMap, err := r.getTagsByBookmarkIDs([]int64{bm.ID})
	if err != nil {
		return nil, fmt.Errorf("fetching bookmark tags failed: %w", err)
	}

	bm.Created = parsedDate
	if tags, found := tagMap[bm.ID]; found {
		bm.Tags = tags
	}

	return &bm, nil
}

func (r *SqliteRepository) Delete(id int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin failed: %w", err)
	}

	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM bookmark_tags WHERE bookmark_id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting bookmark tags failed: %w", err)
	}

	_, err = tx.Exec("DELETE FROM bookmarks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting bookmark failed: %w", err)
	}

	return tx.Commit()
}

func (r *SqliteRepository) GetAll(showPrivate bool, page int, pageSize int) (*BookmarkResult, error) {
	condition := "WHERE is_private = false"
	if showPrivate {
		condition = ""
	}

	query := fmt.Sprintf("SELECT * FROM bookmarks %s ORDER BY created DESC, id DESC LIMIT ?, ?", condition)
	offset := r.calculateOffset(page, pageSize)
	rows, err := r.db.Query(query, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("fetching bookmarks failed: %w", err)
	}
	defer rows.Close()

	bookmarks := make([]*Bookmark, 0)
	for rows.Next() {
		bm, err := r.scanBookmarkRow(rows)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bm)
	}

	err = r.fetchAndSetBookmarkTags(bookmarks)
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM bookmarks %s", condition))
	pageCount, err := r.getRowCount(row, pageSize)
	if err != nil {
		return nil, err
	}

	return &BookmarkResult{
		Bookmarks: bookmarks,
		PageCount: pageCount,
	}, nil
}

func (r *SqliteRepository) GetPrivate(page, pageSize int) (*BookmarkResult, error) {
	offset := r.calculateOffset(page, pageSize)

	rows, err := r.db.Query("SELECT * FROM bookmarks WHERE is_private = true ORDER BY created DESC, id DESC LIMIT ?, ?", offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	bookmarks := make([]*Bookmark, 0)
	for rows.Next() {
		bm, err := r.scanBookmarkRow(rows)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bm)
	}

	err = r.fetchAndSetBookmarkTags(bookmarks)
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow("SELECT COUNT(*) FROM bookmarks WHERE is_private = true")
	pageCount, err := r.getRowCount(row, pageSize)
	if err != nil {
		return nil, fmt.Errorf("fetching row count failed: %w", err)
	}

	return &BookmarkResult{
		Bookmarks: bookmarks,
		PageCount: int(pageCount),
	}, nil
}

func (r *SqliteRepository) GetByTags(showPrivate bool, containsTags []string, page, pageSize int) (*BookmarkResult, error) {
	var qMarks string

	// Question marks for tags in IN query
	if len(containsTags) > 0 {
		qMarks = strings.Repeat("?,", len(containsTags))
		qMarks = qMarks[:len(qMarks)-1]
	}

	var params []any
	for _, tag := range containsTags {
		params = append(params, tag)
	}

	builder := NewConditionBuilder()

	if !showPrivate {
		builder.Add("is_private = false")
	}

	if len(containsTags) > 0 {
		builder.Add(fmt.Sprintf("id IN (SELECT bookmark_id FROM bookmark_tags WHERE tag IN (%s))", qMarks))
	}

	condition := builder.String()
	query := fmt.Sprintf("SELECT * FROM bookmarks %s ORDER BY created DESC, id DESC LIMIT ?, ?", condition)

	offset := r.calculateOffset(page, pageSize)
	params = append(params, offset, pageSize)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	bookmarks := make([]*Bookmark, 0)

	for rows.Next() {
		bm, err := r.scanBookmarkRow(rows)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bm)
	}

	err = r.fetchAndSetBookmarkTags(bookmarks)
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM bookmarks %s", condition), params...)
	pageCount, err := r.getRowCount(row, pageSize)
	if err != nil {
		return nil, fmt.Errorf("fetching row count failed: %w", err)
	}

	return &BookmarkResult{
		Bookmarks: bookmarks,
		PageCount: int(pageCount),
	}, nil
}

func (r *SqliteRepository) GetByKeyword(showPrivate bool, q string, page, pageSize int) (*BookmarkResult, error) {

	var params = make([]any, 0)
	builder := NewConditionBuilder()

	if !showPrivate {
		builder.Add("is_private = false")
	}

	if q != "" {
		builder.Add("(title LIKE ? or description LIKE ?)")
		params = append(params, fmt.Sprintf("%%%s%%", q), fmt.Sprintf("%%%s%%", q))
	}

	condition := builder.String()

	offset := r.calculateOffset(page, pageSize)
	params = append(params, offset, pageSize)

	query := fmt.Sprintf("SELECT * FROM bookmarks %s ORDER BY created DESC, id DESC LIMIT ?, ?", condition)
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	bookmarks := make([]*Bookmark, 0)

	for rows.Next() {
		bm, err := r.scanBookmarkRow(rows)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bm)
	}

	r.fetchAndSetBookmarkTags(bookmarks)

	query = fmt.Sprintf("SELECT COUNT(*) FROM bookmarks %s", condition)
	row := r.db.QueryRow(query, params...)
	pageCount, err := r.getRowCount(row, pageSize)
	if err != nil {
		return nil, fmt.Errorf("fetching row count failed: %w", err)
	}

	return &BookmarkResult{
		Bookmarks: bookmarks,
		PageCount: int(pageCount),
	}, nil
}

func (r *SqliteRepository) GetTags(showPrivate bool) ([]string, error) {
	query := "SELECT DISTINCT bt.tag FROM bookmark_tags AS bt "
	if !showPrivate {
		query += "LEFT JOIN bookmarks AS bm ON bt.bookmark_id = bm.id WHERE bm.is_private = false "
	}
	query += "ORDER BY bt.tag ASC"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	results := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		results = append(results, tag)
	}

	return results, nil
}

func (r *SqliteRepository) GetCheckable() ([]*Bookmark, error) {
	query := "SELECT * FROM bookmarks WHERE ignore_check = false ORDER BY created DESC, id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("fetching bookmarks failed: %w", err)
	}
	defer rows.Close()

	bookmarks := make([]*Bookmark, 0)
	for rows.Next() {
		bm, err := r.scanBookmarkRow(rows)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bm)
	}

	err = r.fetchAndSetBookmarkTags(bookmarks)
	if err != nil {
		return nil, err
	}

	return bookmarks, nil
}

func (r *SqliteRepository) GetBrokenBookmarks() ([]*Bookmark, error) {
	query := "SELECT * FROM bookmarks WHERE is_working = false ORDER BY created DESC, id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("fetching bookmarks failed: %w", err)
	}
	defer rows.Close()

	bookmarks := make([]*Bookmark, 0)
	for rows.Next() {
		bm, err := r.scanBookmarkRow(rows)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, bm)
	}

	err = r.fetchAndSetBookmarkTags(bookmarks)
	if err != nil {
		return nil, err
	}

	return bookmarks, nil
}

func (r *SqliteRepository) BrokenBookmarksExist() (bool, error) {
	var count int64
	row := r.db.QueryRow("SELECT COUNT(*) FROM bookmarks WHERE is_working = false")
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("scanning query row failed: %w", err)
	}

	return count > 0, nil

}

func (r *SqliteRepository) getTagsByBookmarkIDs(bookmarkIDs []int64) (map[int64][]string, error) {
	if len(bookmarkIDs) == 0 {
		return make(map[int64][]string), nil
	}

	qMarks := strings.Repeat("?,", len(bookmarkIDs))
	qMarks = qMarks[0 : len(qMarks)-1]

	var params []any
	for _, id := range bookmarkIDs {
		params = append(params, id)
	}

	rows, err := r.db.Query(fmt.Sprintf("SELECT bookmark_id, tag FROM bookmark_tags WHERE bookmark_id IN (%s) ORDER BY bookmark_id ASC, tag ASC", qMarks), params...)
	if err != nil {
		return nil, fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	results := make(map[int64][]string)
	for rows.Next() {
		var bookmarkID int64
		var tag string

		if err := rows.Scan(&bookmarkID, &tag); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, fmt.Errorf("row scan failed: %w", err)
		}

		if _, exists := results[bookmarkID]; !exists {
			results[bookmarkID] = make([]string, 0)
		}

		results[bookmarkID] = append(results[bookmarkID], tag)
	}

	return results, nil
}

func (r *SqliteRepository) setBookmarkTags(bookmarkID int64, tags []string, tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM bookmark_tags WHERE bookmark_id = ?", bookmarkID)
	if err != nil {
		return fmt.Errorf("deleting bookmark tags failed: %w", err)
	}

	for _, tag := range tags {
		_, err := tx.Exec("INSERT INTO bookmark_tags (bookmark_id, tag) VALUES (?, ?)", bookmarkID, tag)
		if err != nil {
			return fmt.Errorf("inserting bookmark tags failed: %w", err)
		}
	}

	return nil
}

func (r *SqliteRepository) calculateOffset(page, pageSize int) int {
	return page*pageSize - pageSize
}

func (r *SqliteRepository) getRowCount(row *sql.Row, pageSize int) (int, error) {
	var bookmarkCount int64
	err := row.Scan(&bookmarkCount)
	if err != nil {
		return 0, fmt.Errorf("fetching row count failed: %w", err)
	}

	return int(math.Ceil(float64(bookmarkCount) / float64(pageSize))), nil
}

func (r *SqliteRepository) scanBookmarkRow(rows *sql.Rows) (*Bookmark, error) {
	var bm Bookmark
	var created string
	if err := rows.Scan(&bm.ID, &bm.URL, &bm.Title, &bm.Description, &bm.IsPrivate, &created, &bm.IsWorking, &bm.IgnoreCheck, &bm.LastStatusCode, &bm.ErrorMessage); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scanning bookmark row failed: %w", err)
	}

	parsedDate, err := time.Parse(time.RFC3339, created)
	if err != nil {
		return nil, fmt.Errorf("parsing time failed: %w", err)
	}
	bm.Created = parsedDate
	return &bm, nil
}

func (r *SqliteRepository) fetchAndSetBookmarkTags(bookmarks []*Bookmark) error {
	bookmarkIDs := make([]int64, 0)
	for _, bm := range bookmarks {
		bookmarkIDs = append(bookmarkIDs, bm.ID)
	}

	tags, err := r.getTagsByBookmarkIDs(bookmarkIDs)
	if err != nil {
		return fmt.Errorf("fetching bookmark tags failed: %w", err)
	}

	for i, bm := range bookmarks {
		if tags, found := tags[bm.ID]; found {
			bookmarks[i].Tags = tags
		}
	}

	return nil
}
