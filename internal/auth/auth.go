package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	tokenSecretBytes := []byte(tokenSecret)
	jwt, err := token.SignedString(tokenSecretBytes)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	user_idString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	user_id, err := uuid.Parse(user_idString)
	if err != nil {
		return uuid.Nil, err
	}

	return user_id, nil
}
