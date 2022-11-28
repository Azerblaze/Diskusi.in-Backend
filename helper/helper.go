package helper

import (
	"regexp"
	"unicode"

	"github.com/labstack/echo/v4"
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

func CustomHeaderResponse(c echo.Context) echo.Context {

	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Response().Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	c.Response().Header().Set("Access-Control-Expose-Headers", "Authorization")

	return c
}
