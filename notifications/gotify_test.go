package notifications

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/saaste/bookmark-manager/config"
	"github.com/stretchr/testify/assert"
)

type MockHttpClient struct {
	response *http.Response
	error    error
}

func (c *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if c.error != nil {
		return nil, c.error
	}

	return c.response, nil
}

func TestSendGotifyMessage(t *testing.T) {
	mockGotifyErrorBody := GotifyErrorBody{
		Error:       "mock error",
		Code:        http.StatusInternalServerError,
		Description: "mocked error",
	}

	jsonBytes, err := json.Marshal(mockGotifyErrorBody)
	assert.Nil(t, err)

	tests := []struct {
		testName       string
		gotifyURL      string
		gotifyToken    string
		clientResponse *http.Response
		clientDoError  error
		expectedError  string
	}{
		{
			testName:    "Successful send",
			gotifyURL:   "https://example.org",
			gotifyToken: "test-token",
			clientResponse: &http.Response{
				StatusCode: http.StatusOK,
			},
			clientDoError: nil,
		},
		{
			testName:    "Invalid HTTP status from server",
			gotifyURL:   "https://example.org",
			gotifyToken: "test-token",
			clientResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(string(jsonBytes))),
			},
			clientDoError: nil,
			expectedError: "sending gotify message failed: error mock error, code 500, description mocked error",
		},
		{
			testName:      "Client error",
			gotifyURL:     "https://example.org",
			gotifyToken:   "test-token",
			clientDoError: fmt.Errorf("mock error"),
			expectedError: "sending a request failed: mock error",
		},
		{
			testName:      "Missing gotify URL",
			gotifyURL:     "",
			gotifyToken:   "test-token",
			expectedError: "notifications are disabled because Gotify is not configured",
		},
		{
			testName:      "Missing gotify token",
			gotifyURL:     "https://example.org",
			gotifyToken:   "",
			expectedError: "notifications are disabled because Gotify is not configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			notifier := &Notifier{
				appConfig: &config.AppConfig{
					GotifyURL:   tt.gotifyURL,
					GotifyToken: tt.gotifyToken,
				},
				client: &MockHttpClient{
					response: tt.clientResponse,
					error:    tt.clientDoError,
				},
			}

			err := notifier.SendGotifyMessage("Test title", "Test message")
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}
