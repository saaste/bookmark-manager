package handlers

import (
	"net/http"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := h.isAuthenticated(w, r)
	if isAuthenticated {
		http.Redirect(w, r, h.appConf.BaseURL, http.StatusFound)
		return
	}

	type loginTemplateData struct {
		TemplateData
		Error string
	}

	baseData := h.defaultTemplateData(w, r, isAuthenticated)
	baseData.Title = "Login"

	data := loginTemplateData{
		TemplateData: baseData,
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
			h.auth.SetCookie(w, hash)
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
