package bookmarks

import (
	"fmt"
	"log"
	"net/http"

	"github.com/saaste/bookmark-manager/config"
)

type BookmarkError struct {
	Title   string
	URL     string
	Message string
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type BookmarkChecker struct {
	appConfig *config.AppConfig
	repo      Repository
	client    HttpClient
}

func NewBookmarkChecker(appConfig *config.AppConfig, repo Repository, client HttpClient) *BookmarkChecker {
	return &BookmarkChecker{
		appConfig: appConfig,
		repo:      repo,
		client:    client,
	}
}

func (bc *BookmarkChecker) CheckBookbarks() ([]BookmarkError, error) {
	errors := make([]BookmarkError, 0)
	bms, err := bc.repo.GetCheckable()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookmarks for checking: %v", err)
	}

	log.Printf("Checking %d bookmarks...\n", len(bms))
	for _, bookmark := range bms {

		req, err := http.NewRequest(http.MethodGet, bookmark.URL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create a request: %v", err)
		}
		req.Header.Add("User-Agent", bc.userAgentString())
		resp, err := bc.client.Do(req)
		working := true
		if err != nil {
			errors = append(errors, BookmarkError{
				Title:   bookmark.Title,
				URL:     bookmark.URL,
				Message: err.Error(),
			})
			working = false
			fmt.Printf("Bookmark #%d failed: %v\n", bookmark.ID, err)
		} else if resp.StatusCode >= 300 {
			errors = append(errors, BookmarkError{
				Title:   bookmark.Title,
				URL:     bookmark.URL,
				Message: fmt.Sprintf("Returned %s", resp.Status),
			})
			working = false
			fmt.Printf("Bookmark #%d failed with status: %s\n", bookmark.ID, resp.Status)
		}

		bookmark.IsWorking = working
		_, err = bc.repo.Update(bookmark)
		if err != nil {
			return errors, err
		}
	}
	log.Printf("Bookmarks check done! Found %d errors\n", len(errors))
	return errors, nil
}

func (bc *BookmarkChecker) userAgentString() string {
	if bc.appConfig.CheckerUserAgent != "" {
		return bc.appConfig.CheckerUserAgent
	}
	return fmt.Sprintf("Mozilla/5.0 +https://github.com/saaste/bookmark-manager BookmarkManager/%s", bc.appConfig.AppVersion)
}
