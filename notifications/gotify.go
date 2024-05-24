package notifications

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/saaste/bookmark-manager/config"
)

type MessageRequestBody struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type GotifyErrorBody struct {
	Error       string `json:"error"`
	Code        int    `json:"errorCode"`
	Description string `json:"errorDescription"`
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
type Notifier struct {
	appConfig *config.AppConfig
	client    HttpClient
}

func NewNotifier(appConfig *config.AppConfig) *Notifier {
	return &Notifier{
		appConfig: appConfig,
		client:    http.DefaultClient,
	}
}

func (n *Notifier) SendGotifyMessage(title string, message string) error {
	if n.appConfig.GotifyURL == "" || n.appConfig.GotifyToken == "" {
		return errors.New("notifications are disabled because Gotify is not configured")
	}

	body := &MessageRequestBody{
		Title:   title,
		Message: message,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshaling request body failed: %w", err)
	}

	gotifyUrl := fmt.Sprintf("%smessage", n.appConfig.GotifyURL)
	req, err := http.NewRequest(http.MethodPost, gotifyUrl, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("creating a request failed: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Gotify-Key", n.appConfig.GotifyToken)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending a request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("sending gotify message failed with status %s", resp.Status)
		}
		var errorBody GotifyErrorBody
		if err := json.Unmarshal(body, &errorBody); err != nil {
			return fmt.Errorf("sending gotify message failed with status %s", resp.Status)
		}
		return fmt.Errorf("sending gotify message failed: error %s, code %d, description %s", errorBody.Error, errorBody.Code, errorBody.Description)
	}

	return nil
}
