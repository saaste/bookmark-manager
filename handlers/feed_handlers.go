package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
	"github.com/saaste/bookmark-manager/bookmarks"
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

	h.bookmarksToFeed(bookmarkResult.Bookmarks, w)
}

func (h *Handler) bookmarksToFeed(bookmarks []*bookmarks.Bookmark, w http.ResponseWriter) {
	updated := time.Now()
	if len(bookmarks) > 0 {
		updated = bookmarks[0].Created
	}

	feed := &feeds.Feed{
		Title:       h.appConf.SiteName,
		Description: h.appConf.Description,
		Link:        &feeds.Link{Href: h.appConf.BaseURL, Rel: "self"},
		Updated:     updated,
	}

	for _, bm := range bookmarks {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       bm.Title,
			Link:        &feeds.Link{Href: bm.URL},
			Description: bm.Description,
			Updated:     bm.Created,
		})
	}

	atom, err := feed.ToAtom()
	if err != nil {
		h.internalServerError(w, "Failed to create atom feed", err)
		return
	}

	w.Header().Set("Content-Type", "application/atom+xml")
	io.WriteString(w, atom)
}
