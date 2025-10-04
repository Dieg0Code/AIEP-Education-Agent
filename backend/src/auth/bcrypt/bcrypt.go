package bcrypt

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type bcryptImpl struct{}

// CompareHashAndPassword implements Bcrypt.
func (b *bcryptImpl) CompareHashAndPassword(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println("Error comparing hash and password:", err)
		return err
	}
	return nil
}

// HashPassword implements Bcrypt.
func (b *bcryptImpl) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return "", err
	}

	return string(hashedPassword), nil
}

// NewBcrypt crea una instancia de Bcrypt.
func NewBcrypt() Bcrypt {
	return &bcryptImpl{}
}
