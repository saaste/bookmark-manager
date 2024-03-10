package handlers

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/render"
)

func (h *Handler) HandleAPIMetadata(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.authenticateAPI(w, r)
	if !isAuthenticated {
		return
	}

	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch URL %s: %v", url, err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse the document from %s: %v\n", url, err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var title string
	var description string

	title = doc.Find("title").Text()

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if description != "" {
			return
		}

		if prop, exists := s.Attr("property"); exists && (prop == "description" || prop == "og:description") {
			if content, exists := s.Attr("content"); exists {
				description = content
			}
		}
	})

	type MetadataResponse struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	meta := MetadataResponse{
		Title:       title,
		Description: description,
	}

	render.JSON(w, r, meta)
}

func (h *Handler) HandleAPITags(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.authenticateAPI(w, r)
	if !isAuthenticated {
		return
	}

	tags, err := h.bookmarkRepo.GetTags(isAuthenticated)
	if err != nil {
		h.internalServerError(w, "failed to fetch tags", err)
		return
	}

	type TagsResponse struct {
		Tags []string `json:"tags"`
	}

	response := TagsResponse{
		Tags: tags,
	}

	render.JSON(w, r, response)
}
