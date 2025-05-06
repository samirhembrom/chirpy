package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "hello@123"
	passwordByte := []byte(password)
	_, err := bcrypt.GenerateFromPassword(passwordByte, 10)
	if err != nil {
		t.Errorf(`HashPassword("hello@123") = %v, want match for nil`, err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "hello@123"
	passwordByte := []byte(password)
	hashedPasswordByte, err := bcrypt.GenerateFromPassword(passwordByte, 10)
	hashedPassword := string(hashedPasswordByte)
	err = CheckPasswordHash(hashedPassword, password)
	if err != nil {
		t.Errorf(`CheckPasswordHash("hello@123") = %v, want match for nil`, err)
	}
}

func TestMakeJWT(t *testing.T) {
	user_id, err := uuid.NewRandom()
	_, err = MakeJWT(user_id, "jackthereaper", time.Duration(1))
	if err != nil {
		t.Errorf(`MakeJWT("random, jackthereaper, 1" = %v, want match for nil)`, err)
	}
}

func TestValidateJWT(t *testing.T) {
	user_id, err := uuid.NewRandom()
	expiresIn, err := time.ParseDuration("1m")
	jwt, err := MakeJWT(user_id, "jackthereaper", expiresIn)
	if err != nil {
		t.Errorf(`MakeJWT("random, jackthereaper, 1" = %v, want match for nil)`, err)
	}
	valid_user_id, err := ValidateJWT(jwt, "jackthereaper")
	if user_id != valid_user_id || err != nil {
		t.Errorf(
			`ValidateJWT("jwt, jackthereaper" = %v, %v, want match for %v, nil)`,
			valid_user_id,
			err,
			user_id,
		)
	}
}
