package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		return nil, err
	}
	return pass, nil
}
