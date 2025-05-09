package bookmarks

import (
	"time"
)

type Bookmark struct {
	ID             int64
	URL            string
	Title          string
	Description    string
	IsPrivate      bool
	Created        time.Time
	Tags           []string
	IsWorking      bool
	IgnoreCheck    bool
	LastStatusCode int
	ErrorMessage   string
	NextCheck      *time.Time
}

type BookmarkResult struct {
	Bookmarks []*Bookmark
	PageCount int
}
