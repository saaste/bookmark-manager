package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/saaste/bookmark-manager/bookmarks"
)

func (h *Handler) HandlePrivateBookmarks(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, fmt.Sprintf("%slogin", h.appConf.BaseURL), http.StatusFound)
		return
	}

	page := h.getPageParam(r)

	bookmarkResult, err := h.bookmarkRepo.GetPrivate(page, h.appConf.PageSize)
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
	data.Title = "Private Bookmarks"
	data.Bookmarks = bookmarkResult.Bookmarks
	data.Tags = allTags
	data.Pages = h.getPages(page, bookmarkResult.PageCount)

	h.parseTemplateWithFunc("index.html", r, w, data)
}

func (h *Handler) HandleBrokenBookmarks(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, fmt.Sprintf("%slogin", h.appConf.BaseURL), http.StatusFound)
		return
	}

	bookmarks, err := h.bookmarkRepo.GetBrokenBookmarks()
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
	data.Title = "Broken Bookmarks"
	data.Bookmarks = bookmarks
	data.Tags = allTags

	h.parseTemplateWithFunc("index.html", r, w, data)
}

func (h *Handler) HandleBookmarkAdd(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, fmt.Sprintf("%slogin", h.appConf.BaseURL), http.StatusFound)
		return
	}

	baseData := h.defaultTemplateData(w, r, isAuthenticated)
	baseData.Title = "Add Bookmark"

	data := adminTemplateData{
		TemplateData: baseData,
		Errors:       make(map[string]string),
		Bookmark:     &bookmarks.Bookmark{},
		Tags:         "",
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.internalServerError(w, "Failed to parse form", err)
			return
		}

		if r.Form.Get("url") == "" {
			data.Errors["url"] = "Required"
		}

		if r.Form.Get("title") == "" {
			data.Errors["title"] = "Required"
		}

		data.Bookmark.URL = r.Form.Get("url")
		data.Bookmark.Title = r.Form.Get("title")
		data.Bookmark.Description = r.Form.Get("description")
		data.Bookmark.IsPrivate = r.Form.Get("is_private") == "1"
		data.Bookmark.Created = time.Now().UTC()
		data.Bookmark.IsWorking = true

		data.Tags = r.Form.Get("tags")
		if len(strings.TrimSpace(data.Tags)) > 0 {
			data.Bookmark.Tags = strings.Split(data.Tags, " ")
		}

		if len(data.Errors) == 0 {
			_, err := h.bookmarkRepo.Create(data.Bookmark)
			if err != nil {
				h.internalServerError(w, "Failed to create a bookmark", err)
				return
			}
			http.Redirect(w, r, h.appConf.BaseURL, http.StatusFound)
			return
		}
	}

	h.parseTemplateWithFunc("admin_bookmark_add.html", r, w, data)
}

func (h *Handler) HandleBookmarkEdit(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, fmt.Sprintf("%slogin", h.appConf.BaseURL), http.StatusFound)
		return
	}

	idParam := chi.URLParam(r, "bookmarkID")
	if idParam == "" {
		http.NotFound(w, r)
		return
	}

	bookmarkID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Invalid bookmark ID in the path: %s", idParam)
		http.NotFound(w, r)
		return
	}

	bookmark, err := h.bookmarkRepo.Get(bookmarkID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	baseData := h.defaultTemplateData(w, r, isAuthenticated)
	baseData.Title = "Edit Bookmark"

	data := adminTemplateData{
		TemplateData: baseData,
		Errors:       make(map[string]string),
		Bookmark:     bookmark,
		Tags:         strings.Join(bookmark.Tags, " "),
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.internalServerError(w, "Failed to parse form", err)
			return
		}

		if r.Form.Get("url") == "" {
			data.Errors["url"] = "Required"
		}

		if r.Form.Get("title") == "" {
			data.Errors["title"] = "Required"
		}

		data.Bookmark.URL = r.Form.Get("url")
		data.Bookmark.Title = r.Form.Get("title")
		data.Bookmark.Description = r.Form.Get("description")
		data.Bookmark.IsPrivate = r.Form.Get("is_private") == "1"
		data.Bookmark.IsWorking = r.Form.Get("is_working") == "1"
		data.Bookmark.Created = time.Now().UTC()

		data.Tags = r.Form.Get("tags")
		if len(strings.TrimSpace(data.Tags)) > 0 {
			data.Bookmark.Tags = strings.Split(data.Tags, " ")
		} else {
			data.Bookmark.Tags = make([]string, 0)
		}

		if len(data.Errors) == 0 {
			_, err := h.bookmarkRepo.Update(data.Bookmark)
			if err != nil {
				h.internalServerError(w, "Failed to create a bookmark", err)
				return
			}
			http.Redirect(w, r, h.appConf.BaseURL, http.StatusFound)
			return
		}
	}

	h.parseTemplateWithFunc("admin_bookmark_edit.html", r, w, data)
}

func (h *Handler) HandleBookmarkDelete(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, fmt.Sprintf("%slogin", h.appConf.BaseURL), http.StatusFound)
		return
	}

	idParam := chi.URLParam(r, "bookmarkID")
	if idParam == "" {
		log.Printf("No bookmark ID in the path")
		http.NotFound(w, r)
		return
	}

	bookmarkID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Invalid bookmark ID in the path: %s", idParam)
		http.NotFound(w, r)
		return
	}

	bookmark, err := h.bookmarkRepo.Get(bookmarkID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	baseData := h.defaultTemplateData(w, r, isAuthenticated)
	baseData.Title = "Delete Bookmark"

	data := adminTemplateData{
		TemplateData: baseData,
		Bookmark:     bookmark,
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.internalServerError(w, "Failed to parse form", err)
			return
		}

		err = h.bookmarkRepo.Delete(bookmark.ID)
		if err != nil {
			h.internalServerError(w, "Failed to delete a bookmark", err)
			return
		}

		http.Redirect(w, r, h.appConf.BaseURL, http.StatusFound)
		return
	}

	h.parseTemplateWithFunc("admin_bookmark_delete.html", r, w, data)
}
