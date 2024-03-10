package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/saaste/bookmark-manager/bookmarks"
)

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	q := r.URL.Query().Get("q")
	page := h.getPageParam(r)

	bookmarkResult, err := h.getBookmarksWithPagination(isAuthenticated, q, "", page, h.appConf.PageSize)
	if err != nil {
		h.internalServerError(w, "Failed to fetch bookmarks", err)
		return
	}

	allTags, err := h.bookmarkRepo.GetTags(isAuthenticated)
	if err != nil {
		h.internalServerError(w, "Failed to fetch tags", err)
		return
	}

	brokenBookmarks := make([]*bookmarks.Bookmark, 0)
	if isAuthenticated {
		broken, err := h.bookmarkRepo.GetBrokenBookmarks()
		if err != nil {
			h.internalServerError(w, "Failed to fetch broken bookmarks", err)
		}
		brokenBookmarks = broken
	}

	title := "Recent Bookmarks"
	if q != "" {
		title = "Search Results"
	}

	data := templateData{
		SiteName:        h.appConf.SiteName,
		Description:     h.appConf.Description,
		Title:           title,
		BaseURL:         h.appConf.BaseURL,
		CurrentURL:      h.getCurrentURL(r, h.appConf),
		IsAuthenticated: isAuthenticated,
		Bookmarks:       bookmarkResult.Bookmarks,
		Tags:            allTags,
		TextFilter:      q,
		Pages:           h.getPages(page, bookmarkResult.PageCount),
		BrokenBookmarks: brokenBookmarks,
	}

	h.parseTemplateWithFunc("index.html", r, w, data)
}

func (h *Handler) HandleTags(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	tagsParam := chi.URLParam(r, "tags")
	page := h.getPageParam(r)

	bookmarkResult, err := h.getBookmarksWithPagination(isAuthenticated, "", tagsParam, page, h.appConf.PageSize)
	if err != nil {
		h.internalServerError(w, "Failed to fetch bookmarks", err)
		return
	}

	allTags, err := h.bookmarkRepo.GetTags(isAuthenticated)
	if err != nil {
		h.internalServerError(w, "Failed to fetch tags", err)
		return
	}

	brokenBookmarks := make([]*bookmarks.Bookmark, 0)
	if isAuthenticated {
		broken, err := h.bookmarkRepo.GetBrokenBookmarks()
		if err != nil {
			h.internalServerError(w, "Failed to fetch broken bookmarks", err)
		}
		brokenBookmarks = broken
	}

	data := templateData{
		SiteName:        h.appConf.SiteName,
		Description:     h.appConf.Description,
		Title:           fmt.Sprintf("Bookmarks With Tag: %s", tagsParam),
		BaseURL:         h.appConf.BaseURL,
		CurrentURL:      h.getCurrentURL(r, h.appConf),
		IsAuthenticated: isAuthenticated,
		Bookmarks:       bookmarkResult.Bookmarks,
		Tags:            allTags,
		TagFilter:       tagsParam,
		Pages:           h.getPages(page, bookmarkResult.PageCount),
		BrokenBookmarks: brokenBookmarks,
	}

	h.parseTemplateWithFunc("index.html", r, w, data)
}
