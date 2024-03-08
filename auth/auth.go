package auth

import (
	"log"
	"net/http"

	"github.com/saaste/bookmark-manager/config"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	appConf *config.AppConfig
}

func NewAuthenticator(appConf *config.AppConfig) *Authenticator {
	return &Authenticator{
		appConf: appConf,
	}
}

func (a *Authenticator) CreateCookieValue() string {
	return a.calculateHash()
}

func (a *Authenticator) IsValid(c *http.Cookie) bool {
	if c.Value == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(c.Value), []byte(a.appConf.Password+a.appConf.Secret))
	return err == nil
}

func (a *Authenticator) calculateHash() string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(a.appConf.Password+a.appConf.Secret), 4)
	if err != nil {
		log.Printf("failed to create the hash: %v\n", err)
		return ""
	}
	return string(bytes)
}
