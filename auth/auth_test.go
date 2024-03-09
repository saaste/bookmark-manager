package auth

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/saaste/bookmark-manager/config"
	"github.com/stretchr/testify/assert"
)

var appConf = &config.AppConfig{
	Password: "password",
	Secret:   "supersecret",
}

func TestCalculateHashSuccess(t *testing.T) {
	authenticator := &Authenticator{
		appConf: appConf,
		generateFromPassword: func(password []byte, cost int) ([]byte, error) {
			return []byte("test"), nil
		},
	}
	hash, err := authenticator.CalculateHash()
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}

func TestCalculateHashFailure(t *testing.T) {
	authenticator := &Authenticator{
		appConf: appConf,
		generateFromPassword: func(password []byte, cost int) ([]byte, error) {
			return nil, fmt.Errorf("mock error")
		},
	}
	hash, err := authenticator.CalculateHash()
	assert.EqualError(t, err, "mock error")
	assert.Empty(t, hash)
}

func TestIsValidSuccess(t *testing.T) {
	authenticator := &Authenticator{
		appConf: appConf,
		compareHashAndPassword: func(hash, password []byte) error {
			return nil
		},
	}

	cookie := &http.Cookie{
		Value: "test",
	}

	assert.True(t, authenticator.IsValid(cookie))
}

func TestIsInvalidInvalidHash(t *testing.T) {
	authenticator := &Authenticator{
		appConf: appConf,
		compareHashAndPassword: func(hash, password []byte) error {
			return fmt.Errorf("invalid hash")
		},
	}

	cookie := &http.Cookie{
		Value: "test",
	}

	assert.False(t, authenticator.IsValid(cookie))
}

func TestIsValidEmptyHash(t *testing.T) {
	authenticator := &Authenticator{
		appConf: appConf,
	}

	cookie := &http.Cookie{}
	assert.False(t, authenticator.IsValid(cookie))
}
