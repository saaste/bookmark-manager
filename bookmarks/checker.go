package bookmarks

import (
	"fmt"
	"log"
	"net/http"
)

type BookmarkError struct {
	Title   string
	URL     string
	Message string
}

type BookmarkChecker struct {
	repo Repository
	get  func(url string) (resp *http.Response, err error)
}

func NewBookmarkChecker(repo Repository) *BookmarkChecker {
	return &BookmarkChecker{
		repo: repo,
		get:  http.Get,
	}
}

func (bc *BookmarkChecker) CheckBookbarks() ([]BookmarkError, error) {
	errors := make([]BookmarkError, 0)
	bms, err := bc.repo.GetAllWithoutPagination()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookmarks for checking: %v", err)
	}

	log.Printf("Checking %d bookmarks...\n", len(bms))
	for _, bookmark := range bms {
		resp, err := bc.get(bookmark.URL)
		if err != nil {
			errors = append(errors, BookmarkError{
				Title:   bookmark.Title,
				URL:     bookmark.URL,
				Message: err.Error(),
			})
		} else if resp.StatusCode >= 300 {
			errors = append(errors, BookmarkError{
				Title:   bookmark.Title,
				URL:     bookmark.URL,
				Message: fmt.Sprintf("Returned %s", resp.Status),
			})
		}
	}
	log.Printf("Bookmarks check done! Found %d errors\n", len(errors))
	return errors, nil
}
