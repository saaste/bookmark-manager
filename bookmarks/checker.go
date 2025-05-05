package bookmarks

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/saaste/bookmark-manager/config"
)

type BookmarkError struct {
	Title   string
	URL     string
	Message string
}

type URLCheckResult struct {
	StatusCode int
	Error      error
	NextCheck  *time.Time
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

		// Try first with HEAD request
		checkResult := bc.checkWithHead(bookmark.URL)
		if checkResult.Error != nil {
			// If that fails, try with GET request
			checkResult = bc.checkWithGet(bookmark.URL)
			if checkResult.Error != nil {
				fmt.Printf("Bookmark #%d failed with status %d: %v\n", bookmark.ID, checkResult.StatusCode, checkResult.Error)
				errors = append(errors, BookmarkError{
					Title:   bookmark.Title,
					URL:     bookmark.URL,
					Message: checkResult.Error.Error(),
				})
			}
		}

		bookmark.IsWorking = checkResult.Error == nil
		bookmark.LastStatusCode = checkResult.StatusCode
		bookmark.NextCheck = checkResult.NextCheck
		if checkResult.Error != nil {
			bookmark.ErrorMessage = checkResult.Error.Error()
		} else {
			bookmark.ErrorMessage = ""
		}

		_, err = bc.repo.Update(bookmark)
		if err != nil {
			return errors, err
		}
	}
	log.Printf("Bookmarks check done! Found %d errors\n", len(errors))
	return errors, nil
}

func (bc *BookmarkChecker) checkWithHead(siteUrl string) *URLCheckResult {
	return bc.checkURL(siteUrl, http.MethodHead)
}

func (bc *BookmarkChecker) checkWithGet(siteURL string) *URLCheckResult {
	return bc.checkURL(siteURL, http.MethodGet)
}

func (bc *BookmarkChecker) checkURL(siteUrl string, method string) *URLCheckResult {
	req, err := http.NewRequest(method, siteUrl, nil)
	if err != nil {
		return &URLCheckResult{Error: fmt.Errorf("failed to create a request: %v", err)}
	}

	parsedURL, err := url.Parse(siteUrl)
	if err != nil {
		return &URLCheckResult{Error: fmt.Errorf("failed to parse URL: %v", err)}
	}

	req.Header.Set("User-Agent", bc.userAgentString())
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Host", parsedURL.Host)

	resp, err := bc.client.Do(req)
	if err != nil {
		return &URLCheckResult{Error: fmt.Errorf("failed to send a request to %s: %w", siteUrl, err)}
	}

	if resp.StatusCode >= 300 {

		if resp.StatusCode == http.StatusTooManyRequests {
			// Check for Retry-After header
			nextCheck := time.Now().Add(24 * 7 * time.Hour)
			if resp.Header.Get("Retry-After") != "" {
				retryAfterSeconds, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
				if err != nil {
					return &URLCheckResult{Error: fmt.Errorf("failed to parse Retry-After header: %v", err)}
				}
				nextCheck = time.Now().Add(time.Duration(retryAfterSeconds) * time.Second)
			}
			return &URLCheckResult{
				StatusCode: resp.StatusCode,
				Error:      errors.New("429 Too Many Requests"),
				NextCheck:  &nextCheck,
			}
		}

		errorMessage := resp.Status

		if resp.Header.Get("Cf-Mitigated") == "challenge" {
			errorMessage = "Check blocked by Cloudflare challenge"
		}

		return &URLCheckResult{
			StatusCode: resp.StatusCode,
			Error:      errors.New(errorMessage),
		}
	}

	return &URLCheckResult{
		StatusCode: resp.StatusCode,
	}
}

func (bc *BookmarkChecker) userAgentString() string {
	if bc.appConfig.CheckerUserAgent != "" {
		return bc.appConfig.CheckerUserAgent
	}
	return fmt.Sprintf("Mozilla/5.0 +https://github.com/saaste/bookmark-manager BookmarkManager/%s", bc.appConfig.AppVersion)
}
