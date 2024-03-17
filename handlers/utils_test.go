package handlers

import (
	"net/http"
	"testing"

	"github.com/saaste/bookmark-manager/config"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentURL(t *testing.T) {
	handler := Handler{
		appConf: &config.AppConfig{
			BaseURL: "https://example.com/",
		},
	}

	request := http.Request{
		RequestURI: "/feed",
	}

	actual := handler.getCurrentURL(&request)
	assert.Equal(t, actual, "https://example.com/feed")
}
