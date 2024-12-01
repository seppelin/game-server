package models

import (
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserID int64

const USER_ANON UserID = -1

type User struct {
	ID      UserID
	Name    string
	Email   string
	PwdHash string
}

func ValidPassword(password string) bool {
	pattern := `^[a-zA-Z\d@$!%*?&#]{4,32}`
	regex := regexp.MustCompile(pattern)
	return regex.Match([]byte(password))
}

func ValidUsername(name string) bool {
	pattern := `^[a-zA-Z0-9_]{3,20}`
	regex := regexp.MustCompile(pattern)
	return regex.Match([]byte(name))
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func VerifyPassword(pwd_hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(pwd_hash), []byte(password))
}

type UserDB interface {
	InsertUser(string, string) UserID
	UpdateUser(User) error
	User(UserID) User
}
