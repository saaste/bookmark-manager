package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ServeFiles(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (h *Handler) HandleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("robots.txt")
	if err != nil {
		h.internalServerError(w, err.Error(), err)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.Write(content)
}
