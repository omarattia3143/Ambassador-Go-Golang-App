package controllers

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func generateJwtForUser(user *models.User) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
}

func Register(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	if data["Password"] != data["PasswordConfirm"] {

		c.Status(fiber.StatusBadRequest)
		err := c.JSON(fiber.Map{"message": "passwords doesn't match"})
		if err != nil {
			return err
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["Password"]), 12)
	if err != nil {
		return err
	}

	user := models.User{
		FirstName:    data["FirstName"],
		LastName:     data["LastName"],
		Email:        data["Email"],
		Password:     hashedPassword,
		IsAmbassador: false,
	}

	var userExists bool
	err = database.DB.Model(&user).Where("first_name = ?", user.FirstName).Find(&userExists).Error
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "user already exists"})
	}

	var token string

	err = database.DB.Transaction(func(tx *gorm.DB) error {

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(&user).Error; err != nil {
			// return any error will rollback
			return err
		}

		token, err = generateJwtForUser(&user)
		if err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		return err
	}

	return c.JSON(token)
}

func Login(c *fiber.Ctx) error {

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["Email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials 1",
		})
	}

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["Password"]))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials 2",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Welcome!",
	})
}
