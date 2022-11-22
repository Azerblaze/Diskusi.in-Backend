package helper

import (
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(e)
	return emailRegex
}

func IsStringContainNumber(str string) bool {
	numeric := regexp.MustCompile(`\d`).MatchString(str)
	return numeric
}
func IsStringContainUppercaseLetter(str string) bool {
	for _, v := range str {
		if unicode.IsUpper(v) {
			return true
		}
	}
	return false
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
