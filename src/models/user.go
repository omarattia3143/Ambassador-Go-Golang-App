package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Model
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email" gorm:"unique"`
	Password     []byte   `json:"-"`
	IsAmbassador bool     `json:"-"`
	Revenue      *float64 `json:"revenue" gorm:"-"`
}

func (user *User) FullName() string {
	return user.FirstName + " " + user.LastName
}

func (user *User) ComparePassword(password string) error {

	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

func (user *User) SetPassword(password string) {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	user.Password = hashedPassword
}

type Admin User

func (admin *Admin) CalculateTotalRevenue(db *gorm.DB) float64 {

	var orders []Order

	db.Preload("OrderItems").Find(&orders, Order{
		UserId:   admin.Id,
		Complete: true,
	})

	var revenue float64 = 0
	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AdminRevenue
		}
	}
	return revenue
}

type Ambassador User

func (ambassador *Ambassador) CalculateTotalRevenue(db *gorm.DB) *float64 {
	var orders []Order

	db.Preload("OrderItems").Find(&orders, Order{
		UserId:   ambassador.Id,
		Complete: true,
	})

	var revenue float64 = 0
	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			revenue += orderItem.AmbassadorRevenue
		}
	}

	return &revenue
}
