package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	passwordByte := []byte(password)
	hashedPasswordByte, err := bcrypt.GenerateFromPassword(passwordByte, 10)
	if err != nil {
		return "", err
	}
	hashedPassword := string(hashedPasswordByte)
	return hashedPassword, nil
}

func CheckPasswordHash(hash, password string) error {
	passwordByte := []byte(password)
	hashByte := []byte(hash)
	err := bcrypt.CompareHashAndPassword(hashByte, passwordByte)
	return err
}
