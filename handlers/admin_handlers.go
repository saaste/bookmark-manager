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

	data := templateData{
		SiteName:        h.appConf.SiteName,
		BaseURL:         h.appConf.BaseURL,
		IsAuthenticated: isAuthenticated,
		PrivateOnly:     true,
		Bookmarks:       bookmarkResult.Bookmarks,
		Tags:            allTags,
		Pages:           h.getPages(page, bookmarkResult.PageCount),
	}

	h.parseTemplateWithFunc("index.html", r, w, data)
}

func (h *Handler) HandleBookmarkAdd(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, fmt.Sprintf("%s/login", h.appConf.BaseURL), http.StatusFound)
		return
	}

	data := adminTemplateData{
		templateData: templateData{
			SiteName:        h.appConf.SiteName,
			BaseURL:         h.appConf.BaseURL,
			IsAuthenticated: isAuthenticated,
		},
		Errors:   make(map[string]string),
		Bookmark: bookmarks.Bookmark{},
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

	data := adminTemplateData{
		templateData: templateData{
			SiteName:        h.appConf.SiteName,
			BaseURL:         h.appConf.BaseURL,
			IsAuthenticated: isAuthenticated,
		},
		Errors:   make(map[string]string),
		Bookmark: *bookmark,
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

	data := adminTemplateData{
		templateData: templateData{
			SiteName:        h.appConf.SiteName,
			BaseURL:         h.appConf.BaseURL,
			IsAuthenticated: isAuthenticated,
		},
		Bookmark: *bookmark,
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.internalServerError(w, "Failed to parse form", err)
			return
		}

		action := r.Form.Get("action")
		if action == "delete" {
			err = h.bookmarkRepo.Delete(bookmark.ID)
			if err != nil {
				h.internalServerError(w, "Failed to delete a bookmark", err)
				return
			}
		}

		http.Redirect(w, r, h.appConf.BaseURL, http.StatusFound)
		return
	}

	h.parseTemplateWithFunc("admin_bookmark_delete.html", r, w, data)
}
