package handlers

import (
	"database/sql"

	"github.com/saaste/bookmark-manager/auth"
	"github.com/saaste/bookmark-manager/bookmarks"
	"github.com/saaste/bookmark-manager/config"
)

type Page struct {
	Number   int
	IsActive bool
}

type templateData struct {
	SiteName        string
	Description     string
	BaseURL         string
	Title           string
	CurrentURL      string
	IsAuthenticated bool
	PrivateOnly     bool
	Bookmarks       []*bookmarks.Bookmark
	Tags            []string
	TagFilter       string
	TextFilter      string
	Pages           []Page
	BrokenBookmarks []*bookmarks.Bookmark
}

type adminTemplateData struct {
	templateData
	Errors   map[string]string
	Bookmark *bookmarks.Bookmark
	Tags     string
}

type Handler struct {
	bookmarkRepo bookmarks.Repository
	appConf      *config.AppConfig
	auth         *auth.Authenticator
}

func NewHandler(db *sql.DB, appConf *config.AppConfig, auth *auth.Authenticator) *Handler {
	return &Handler{
		bookmarkRepo: bookmarks.NewSqliteRepository(db),
		appConf:      appConf,
		auth:         auth,
	}
}
