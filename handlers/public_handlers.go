package handlers

import (
	"fmt"
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

	title := "Recent Bookmarks"
	if q != "" {
		title = fmt.Sprintf("Search Results: \"%s\"", q)
	}

	data := h.defaultTemplateData(w, r, isAuthenticated)
	data.Title = title
	data.Bookmarks = bookmarkResult.Bookmarks
	data.Tags = allTags
	data.TextFilter = q
	data.Pages = h.getPages(page, bookmarkResult.PageCount)

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

	data := h.defaultTemplateData(w, r, isAuthenticated)
	data.Title = fmt.Sprintf("Bookmarks With Tag: %s", tagsParam)
	data.Bookmarks = bookmarkResult.Bookmarks
	data.Tags = allTags
	data.TagFilter = tagsParam
	data.Pages = h.getPages(page, bookmarkResult.PageCount)

	h.parseTemplateWithFunc("index.html", r, w, data)
}

func (h *Handler) defaultTemplateData(w http.ResponseWriter, r *http.Request, isAuthenticated bool) templateData {
	data := templateData{}

	brokenBookmarksExist := false
	if isAuthenticated {
		exists, err := h.bookmarkRepo.BrokenBookmarksExist()
		if err != nil {
			h.internalServerError(w, "Failed to check if broken bookmarks exist", err)
			return data
		}
		brokenBookmarksExist = exists
	}

	data.SiteName = h.appConf.SiteName
	data.Description = h.appConf.Description
	data.BaseURL = h.appConf.BaseURL
	data.AppVersion = h.appConf.AppVersion
	data.CurrentURL = h.getCurrentURL(r, h.appConf)
	data.IsAuthenticated = isAuthenticated
	data.BrokenBookmarksExist = brokenBookmarksExist

	return data
}
