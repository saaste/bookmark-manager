package bookmarks

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/saaste/bookmark-manager/config"
	"github.com/saaste/bookmark-manager/test_utils"
	"github.com/stretchr/testify/assert"
)

type MockHttp struct {
	response *http.Response
	getError error
}

func (h MockHttp) Do(req *http.Request) (resp *http.Response, err error) {
	if h.getError != nil {
		return nil, h.getError
	}
	return h.response, nil
}

func TestCheckBookbarks(t *testing.T) {
	db := test_utils.InitTestDatabase(t, dbFileName)
	defer test_utils.DestroyTestDatabase(t, db, dbFileName)
	repo := NewSqliteRepository(db)

	tests := []struct {
		testName       string
		responseCode   int
		responseStatus string
		getError       error
		expected       []BookmarkError
	}{
		{
			testName:       "OK response",
			responseCode:   http.StatusOK,
			responseStatus: "200 OK",
			getError:       nil,
			expected:       make([]BookmarkError, 0),
		},
		{
			testName:       "Status code response",
			responseCode:   http.StatusNotFound,
			responseStatus: "404 Not Found",
			getError:       nil,
			expected: []BookmarkError{
				{
					Title:   "Test title",
					URL:     "https://example.org",
					Message: "404 Not Found",
				},
			},
		},
		{
			testName: "Error response",
			getError: fmt.Errorf("mock error"),
			expected: []BookmarkError{
				{
					Title:   "Test title",
					URL:     "https://example.org",
					Message: "failed to send a request to https://example.org: mock error",
				},
			},
		},
	}

	_, err := repo.Create(createBookmark(false))
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			client := &MockHttp{
				response: &http.Response{
					StatusCode: tt.responseCode,
					Status:     tt.responseStatus,
				},
				getError: tt.getError,
			}

			checker := &BookmarkChecker{
				repo:   repo,
				client: client,
				appConfig: &config.AppConfig{
					AppVersion: "test-version",
				},
			}

			bookmarkErrors, err := checker.CheckBookbarks()
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, bookmarkErrors)
		})
	}

}
