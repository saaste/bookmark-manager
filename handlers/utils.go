package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/saaste/bookmark-manager/bookmarks"
)

func (h *Handler) getBookmarksWithPagination(isAuthenticated bool, q, tags string, page, pageSize int) (*bookmarks.BookmarkResult, error) {
	if q != "" {
		return h.bookmarkRepo.GetByKeyword(isAuthenticated, q, page, pageSize)
	}

	if tags != "" {
		t := strings.Split(tags, " ")
		return h.bookmarkRepo.GetByTags(isAuthenticated, t, page, pageSize)
	}

	return h.bookmarkRepo.GetAll(isAuthenticated, page, pageSize)
}

func (h *Handler) parseTemplateWithFunc(templateFile string, r *http.Request, w http.ResponseWriter, data any) {
	templateFiles := []string{
		h.getTemplateFile("base.html"),
		h.getTemplateFile(templateFile),
	}

	entries, err := os.ReadDir("components")
	if err != nil {
		h.internalServerError(w, "Reading UI components failed", err)
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".html") {
			templateFiles = append(templateFiles, fmt.Sprintf("components/%s", entry.Name()))
		}
	}

	t, err := template.New("").
		Funcs(template.FuncMap{
			"paginationUrl": func(pageNumber int) string {
				return h.getCurrentURIWithParam(r, "page", pageNumber)
			},
			"feedUrl": func(feedType string) string {
				switch strings.ToLower(feedType) {
				case "atom":
					return h.getFeedURL(r, FeedTypeAtom)
				case "rss":
					return h.getFeedURL(r, FeedTypeRSS)
				case "json":
					return h.getFeedURL(r, FeedTypeJSON)
				default:
					return "invalid feed type"
				}

			},
			"anchorUrl": func(id string) string {
				return h.getAnchorURL(r, id)
			},
		}).ParseFiles(templateFiles...)
	if err != nil {
		h.internalServerError(w, fmt.Sprintf("Failed to parse template %s", templateFile), err)
		return
	}

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		h.internalServerError(w, fmt.Sprintf("Failed to execute template %s", templateFile), err)
		return
	}
}

func (h *Handler) getTemplateFile(filename string) string {
	return fmt.Sprintf("templates/%s/%s", h.appConf.Theme, filename)
}

func (h *Handler) isAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("auth")
	if err != nil {
		return false
	}

	isValid := h.auth.IsValid(cookie)

	if isValid {
		h.auth.SetCookie(w, cookie.Value)
	}

	return h.auth.IsValid(cookie)
}

func (h *Handler) authenticateAPI(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	isValid := h.auth.IsValid(cookie)
	if !isValid {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return false
	}

	return true
}

func (h *Handler) internalServerError(w http.ResponseWriter, msg string, err error) {
	log.Printf("ERROR: %s: %v", msg, err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func (h *Handler) getPageParam(r *http.Request) int {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	return page
}

func (h *Handler) getPages(currentPage, pageCount int) []Page {
	pages := make([]Page, 0)
	for i := 1; i <= pageCount; i++ {
		pages = append(pages, Page{
			Number:   i,
			IsActive: i == currentPage,
		})
	}
	return pages
}

func (h *Handler) getCurrentURL(r *http.Request) string {
	requestURI := strings.TrimPrefix(r.RequestURI, "/")
	return fmt.Sprintf("%s%s", h.appConf.BaseURL, requestURI)
}

func (h *Handler) getAnchorURL(r *http.Request, id string) string {
	return fmt.Sprintf("%s#%s", r.RequestURI, id)
}

func (h *Handler) getFeedURL(r *http.Request, feedType FeedType) string {
	q := r.URL.Query().Get("q")
	rawQuery := ""
	if q != "" {
		rawQuery = fmt.Sprintf("q=%s", q)
	}

	var feedPath string
	switch feedType {
	case FeedTypeAtom:
		feedPath = fmt.Sprintf("%s/atom.xml", strings.TrimSuffix(r.URL.Path, "/"))
	case FeedTypeRSS:
		feedPath = fmt.Sprintf("%s/rss.xml", strings.TrimSuffix(r.URL.Path, "/"))
	case FeedTypeJSON:
		feedPath = fmt.Sprintf("%s/feed.json", strings.TrimSuffix(r.URL.Path, "/"))
	}

	url := url.URL{
		Scheme:   r.URL.Scheme,
		Host:     r.URL.Host,
		Path:     feedPath,
		RawQuery: rawQuery,
	}
	return url.String()
}

func (h *Handler) getCurrentURIWithParam(r *http.Request, key string, val interface{}) string {
	queryParams := r.URL.Query()
	var value string

	switch v := val.(type) {
	case string:
		value = v
	case int:
		value = strconv.Itoa(v)
	default:
		value = "error: unknown type"
	}

	queryParams.Set(key, value)
	url := url.URL{
		Scheme:   r.URL.Scheme,
		Host:     r.URL.Host,
		Path:     r.URL.Path,
		RawQuery: queryParams.Encode(),
	}

	return url.String()
}
