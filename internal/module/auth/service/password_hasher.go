package service

import (
	"github.com/gin-demo/recipes-web/internal/module/auth/domain"
	"github.com/matthewhartstonge/argon2"
)

type PasswordHasher interface {
	Hash(string) (string, error)
	Compare(raw, hash string) error
}

type argonPasswordHasher struct {
	cfg *argon2.Config
}

func New() PasswordHasher {
	defaultCfg := argon2.DefaultConfig()
	return &argonPasswordHasher{
		cfg: &defaultCfg,
	}
}

func (up *argonPasswordHasher) Hash(password string) (string, error) {
	encoded, err := up.cfg.HashEncoded([]byte(password))
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func (up *argonPasswordHasher) Compare(raw, hash string) error {
	if _, err := argon2.VerifyEncoded([]byte(raw), []byte(hash)); err != nil {
		return domain.ErrInvalidCredentials
	}

	return nil
}
