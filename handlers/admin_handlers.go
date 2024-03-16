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
		http.Redirect(w, r, fmt.Sprintf("%s/login", h.appConf.BaseURL), http.StatusFound)
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

	brokenBookmarks, err := h.bookmarkRepo.GetBrokenBookmarks()
	if err != nil {
		h.internalServerError(w, "Failed to fetch broken bookmarks", err)
		return
	}

	data := templateData{
		SiteName:        h.appConf.SiteName,
		Description:     h.appConf.Description,
		Title:           "Private Bookmarks",
		BaseURL:         h.appConf.BaseURL,
		CurrentURL:      h.getCurrentURL(r, h.appConf),
		IsAuthenticated: isAuthenticated,
		PrivateOnly:     true,
		Bookmarks:       bookmarkResult.Bookmarks,
		Tags:            allTags,
		Pages:           h.getPages(page, bookmarkResult.PageCount),
		BrokenBookmarks: brokenBookmarks,
	}

	h.parseTemplateWithFunc("index.html", r, w, data)
}

func (h *Handler) HandleBrokenBookmarks(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, fmt.Sprintf("%s/login", h.appConf.BaseURL), http.StatusFound)
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

	data := templateData{
		SiteName:        h.appConf.SiteName,
		Description:     h.appConf.Description,
		Title:           "Broken Bookmarks",
		BaseURL:         h.appConf.BaseURL,
		CurrentURL:      h.getCurrentURL(r, h.appConf),
		IsAuthenticated: isAuthenticated,
		PrivateOnly:     true,
		Bookmarks:       bookmarks,
		Tags:            allTags,
		BrokenBookmarks: bookmarks,
	}

	h.parseTemplateWithFunc("index.html", r, w, data)
}

func (h *Handler) HandleBookmarkAdd(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, fmt.Sprintf("%s/login", h.appConf.BaseURL), http.StatusFound)
		return
	}

	brokenBookmarks, err := h.bookmarkRepo.GetBrokenBookmarks()
	if err != nil {
		h.internalServerError(w, "Failed to fetch broken bookmarks", err)
		return
	}

	data := adminTemplateData{
		templateData: templateData{
			SiteName:        h.appConf.SiteName,
			Description:     h.appConf.Description,
			Title:           "Add Bookmark",
			BaseURL:         h.appConf.BaseURL,
			CurrentURL:      h.getCurrentURL(r, h.appConf),
			IsAuthenticated: isAuthenticated,
			BrokenBookmarks: brokenBookmarks,
		},
		Errors:   make(map[string]string),
		Bookmark: &bookmarks.Bookmark{},
		Tags:     "",
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
		data.Bookmark.Tags = strings.Split(data.Tags, " ")

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
		http.Redirect(w, r, fmt.Sprintf("%s/login", h.appConf.BaseURL), http.StatusFound)
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

	brokenBookmarks, err := h.bookmarkRepo.GetBrokenBookmarks()
	if err != nil {
		h.internalServerError(w, "Failed to fetch broken bookmarks", err)
		return
	}

	data := adminTemplateData{
		templateData: templateData{
			SiteName:        h.appConf.SiteName,
			Description:     h.appConf.Description,
			Title:           "Edit Bookmark",
			BaseURL:         h.appConf.BaseURL,
			CurrentURL:      h.getCurrentURL(r, h.appConf),
			IsAuthenticated: isAuthenticated,
			BrokenBookmarks: brokenBookmarks,
		},
		Errors:   make(map[string]string),
		Bookmark: bookmark,
		Tags:     strings.Join(bookmark.Tags, " "),
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
		data.Bookmark.Tags = strings.Split(data.Tags, " ")

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
		http.Redirect(w, r, fmt.Sprintf("%s/login", h.appConf.BaseURL), http.StatusFound)
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

	// TODO: There is no need to fetch broken bookmarks every time.
	// Just fetch a boolean indicating if broken bookmarks exists
	brokenBookmarks, err := h.bookmarkRepo.GetBrokenBookmarks()
	if err != nil {
		h.internalServerError(w, "Failed to fetch broken bookmarks", err)
		return
	}

	data := adminTemplateData{
		templateData: templateData{
			SiteName:        h.appConf.SiteName,
			Description:     h.appConf.Description,
			Title:           "Delete Bookmark",
			BaseURL:         h.appConf.BaseURL,
			CurrentURL:      h.getCurrentURL(r, h.appConf),
			IsAuthenticated: isAuthenticated,
			BrokenBookmarks: brokenBookmarks,
		},
		Bookmark: bookmark,
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
