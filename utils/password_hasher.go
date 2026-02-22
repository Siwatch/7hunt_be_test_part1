package utils

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

type DefaultPasswordHasher struct {
	cost int
}

func NewBcryptHasher(cost int) PasswordHasher {
	return &DefaultPasswordHasher{cost: cost}
}

func (h *DefaultPasswordHasher) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (h *DefaultPasswordHasher) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
