package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
	"github.com/saaste/bookmark-manager/bookmarks"
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

	var feedType FeedType
	switch true {
	case strings.HasSuffix(r.URL.Path, "atom.xml"):
		feedType = FeedTypeAtom
	case strings.HasSuffix(r.URL.Path, "rss.xml"):
		feedType = FeedTypeRSS
	case strings.HasSuffix(r.URL.Path, "feed.json"):
		feedType = FeedTypeJSON
	default:
		http.Error(w, "Invalid feed type", http.StatusBadRequest)
	}

	bookmarkResult, err := h.getBookmarksWithPagination(isAuthenticated, q, tags, 1, 20)
	if err != nil {
		h.internalServerError(w, "Failed to fetch bookmarks", err)
		return
	}

	content, err := h.bookmarksToFeed(bookmarkResult.Bookmarks, w, r, feedType)
	if err != nil {
		h.internalServerError(w, "failed to generate feed", err)
		return
	}

	fmt.Println(r.URL.Path)

	io.WriteString(w, content)
}

func (h *Handler) bookmarksToFeed(bookmarks []*bookmarks.Bookmark, w http.ResponseWriter, r *http.Request, feedType FeedType) (string, error) {
	updated := time.Now()
	if len(bookmarks) > 0 {
		updated = bookmarks[0].Created
	}

	feed := &feeds.Feed{
		Id:          h.appConf.BaseURL + r.RequestURI,
		Title:       h.appConf.SiteName,
		Description: h.appConf.Description,
		Author: &feeds.Author{
			Name:  h.appConf.AuthorName,
			Email: h.appConf.AuthorEmail,
		},
		Link:    &feeds.Link{Href: h.appConf.BaseURL + r.RequestURI, Rel: "self"},
		Updated: updated,
	}

	for _, bm := range bookmarks {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       bm.Title,
			Link:        &feeds.Link{Href: bm.URL},
			Description: bm.Description,
			Updated:     bm.Created,
		})
	}

	w.Header().Set("Last-Modified", updated.Format(time.RFC1123))

	switch feedType {
	case FeedTypeAtom:
		atom, err := feed.ToAtom()
		if err != nil {
			return "", fmt.Errorf("failed to create Atom feed: %w", err)
		}
		w.Header().Set("Content-Type", "application/atom+xml; charset=utf-8")
		return atom, nil
	case FeedTypeRSS:
		rss, err := feed.ToRss()
		if err != nil {
			return "", fmt.Errorf("failed to create RSS feed: %w", err)
		}
		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		return rss, nil
	case FeedTypeJSON:
		json, err := feed.ToJSON()
		if err != nil {
			return "", fmt.Errorf("failed to create JSON feed: %w", err)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return json, nil
	default:
		return "", fmt.Errorf("unsupported feed type: %d", feedType)
	}
}
