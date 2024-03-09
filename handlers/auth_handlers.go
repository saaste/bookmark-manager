package handlers

import (
	"net/http"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(r)
	if isAuthenticated {
		http.Redirect(w, r, h.appConf.BaseURL, http.StatusFound)
		return
	}

	type loginTemplateData struct {
		templateData
		Error string
	}

	data := loginTemplateData{
		templateData: templateData{
			SiteName:        h.appConf.SiteName,
			Description:     h.appConf.Description,
			BaseURL:         h.appConf.BaseURL,
			CurrentURL:      h.getCurrentURL(r, h.appConf),
			IsAuthenticated: isAuthenticated,
			Title:           "Login",
		},
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			h.internalServerError(w, "Failed to parse form", err)
			return
		}
		if r.Form.Get("password") != h.appConf.Password {
			data.Error = "Invalid password"
		} else {
			hash, err := h.auth.CalculateHash()
			if err != nil {
				h.internalServerError(w, "Failed to calculate password hash", err)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "auth",
				Value:    hash,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, h.appConf.BaseURL, http.StatusFound)
			return
		}
	}

	h.parseTemplateWithFunc("login.html", r, w, data)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		MaxAge: -1,
	})
	http.Redirect(w, r, h.appConf.BaseURL, http.StatusFound)
}
