package controllers

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/middleware"
	"GoAndNextProject/src/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

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

	user := models.User{
		FirstName:    data["FirstName"],
		LastName:     data["LastName"],
		Email:        data["Email"],
		IsAmbassador: false,
	}

	user.SetPassword(data["Password"])

	var userExists bool
	err = database.DB.Model(&user).Where("first_name = ?", user.FirstName).Find(&userExists).Error
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "User already exists",
		})

	}

	var token string

	err = database.DB.Transaction(func(tx *gorm.DB) error {

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(&user).Error; err != nil {
			// return any error will rollback
			return err
		}

		token, err = user.GenerateJwtForUser()
		if err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		return err
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success!!",
	})
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
			"message": "Invalid Credentials",
		})

	}

	if err := user.ComparePassword(data["Password"]); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}

	token, err := user.GenerateJwtForUser()
	if err != nil {
		return err
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "Welcome!",
	})
}

func User(c *fiber.Ctx) error {

	userId, _ := middleware.GetUserId(c)

	var user models.User
	database.DB.Where("id = ?", userId).First(&user)
	if user.Id == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "user does not exist",
		})
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "Success",
	})
}
