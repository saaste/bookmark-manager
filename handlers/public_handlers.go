package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

	data := templateData{
		SiteName:        h.appConf.SiteName,
		BaseURL:         h.appConf.BaseURL,
		IsAuthenticated: isAuthenticated,
		Bookmarks:       bookmarkResult.Bookmarks,
		Tags:            allTags,
		TextFilter:      q,
		Pages:           h.getPages(page, bookmarkResult.PageCount),
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

	data := templateData{
		SiteName:        h.appConf.SiteName,
		BaseURL:         h.appConf.BaseURL,
		IsAuthenticated: isAuthenticated,
		Bookmarks:       bookmarkResult.Bookmarks,
		Tags:            allTags,
		TagFilter:       tagsParam,
		Pages:           h.getPages(page, bookmarkResult.PageCount),
	}

	h.parseTemplateWithFunc("index.html", r, w, data)
}
