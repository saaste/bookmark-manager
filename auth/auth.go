package auth

import (
	"net/http"

	"github.com/saaste/bookmark-manager/config"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	appConf                *config.AppConfig
	generateFromPassword   func([]byte, int) ([]byte, error)
	compareHashAndPassword func([]byte, []byte) error
}

func NewAuthenticator(appConf *config.AppConfig) *Authenticator {
	return &Authenticator{
		appConf:                appConf,
		generateFromPassword:   bcrypt.GenerateFromPassword,
		compareHashAndPassword: bcrypt.CompareHashAndPassword,
	}
}

func (a *Authenticator) IsValid(c *http.Cookie) bool {
	if c.Value == "" {
		return false
	}

	err := a.compareHashAndPassword([]byte(c.Value), []byte(a.appConf.Password+a.appConf.Secret))
	return err == nil
}

func (a *Authenticator) CalculateHash() (string, error) {
	bytes, err := a.generateFromPassword([]byte(a.appConf.Password+a.appConf.Secret), 4)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (a *Authenticator) UpdateCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.MaxAge = 60 * 60 * 24 * 30
	http.SetCookie(w, cookie)
}
