package models

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type User struct {
	Id           uint   `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email" gorm:"unique"`
	Password     []byte `json:"-"`
	IsAmbassador bool   `json:"-"`
}

func (user *User) ComparePassword(password string) error {

	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

func (user *User) SetPassword(password string) {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	user.Password = hashedPassword
}

func (user *User) GenerateJwtForUser() (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	}

	//todo: remove secret from here
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
}
