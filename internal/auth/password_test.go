package auth

import (
	"testing"

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
