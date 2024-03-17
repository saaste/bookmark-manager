package handlers

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/saaste/bookmark-manager/feeds"
)

type FeedType int

const (
	FeedTypeAtom FeedType = iota
	FeedTypeRSS
	FeedTypeJSON
)

func (h *Handler) HandleFeed(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	tags := chi.URLParam(r, "tags")
	q := r.URL.Query().Get("q")

	bookmarkResult, err := h.getBookmarksWithPagination(isAuthenticated, q, tags, 1, 20)
	if err != nil {
		h.internalServerError(w, "Failed to fetch bookmarks", err)
		return
	}

	updated := time.Now()
	if len(bookmarkResult.Bookmarks) > 0 {
		updated = bookmarkResult.Bookmarks[0].Created
	}

	feedInfo := feeds.FeedInfo{
		SiteName:    h.appConf.SiteName,
		Description: h.appConf.Description,
		BaseURL:     h.appConf.BaseURL,
		CurrentURL:  h.getCurrentURL(r),
		AuthorName:  h.appConf.AuthorName,
		AuthorEmail: h.appConf.AuthorEmail,
	}

	var content string
	switch true {
	case strings.HasSuffix(r.URL.Path, "atom.xml"):
		w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
		content = feeds.ToAtom(feedInfo, bookmarkResult.Bookmarks)
	case strings.HasSuffix(r.URL.Path, "rss.xml"):
		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		content = feeds.ToRSS(feedInfo, bookmarkResult.Bookmarks)
	case strings.HasSuffix(r.URL.Path, "feed.json"):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		content, err = feeds.ToJSON(feedInfo, bookmarkResult.Bookmarks)
	default:
		http.Error(w, "Invalid feed type", http.StatusBadRequest)
	}

	if err != nil {
		h.internalServerError(w, "Generating feed failed", err)
		return
	}

	w.Header().Set("Last-Modified", updated.Format(time.RFC1123))
	io.WriteString(w, content)

}
