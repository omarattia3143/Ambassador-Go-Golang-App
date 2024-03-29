package controllers

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/middleware"
	"GoAndNextProject/src/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
	"time"
)

func Register(c *fiber.Ctx) error {

	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {

		c.Status(fiber.StatusBadRequest)
		err := c.JSON(fiber.Map{"message": "passwords doesn't match"})
		if err != nil {
			return err
		}
	}

	user := models.User{
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: strings.Contains(c.Path(), "/api/ambassador"),
	}

	user.SetPassword(data["password"])

	var userExists bool
	err = database.DB.Model(&user).Where("first_name = ?", user.FirstName).Find(&userExists).Error
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "User already exists",
		})

	}

	var token string
	var scope string

	if user.IsAmbassador {
		scope = "ambassador"
	} else {
		scope = "admin"
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(&user).Error; err != nil {
			// return any error will roll back
			return err
		}

		token, err = middleware.GenerateJWT(user.Id, scope)
		if err != nil {
			return err
		}

		// return nil will commit the whole transaction 1
		return nil
	})
	if err != nil {
		return err
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().AddDate(0, 0, 7),
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

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials",
		})

	}

	if err := user.ComparePassword(data["password"]); err != nil {
		return err
	}

	var scope string

	if user.IsAmbassador && strings.Contains(c.Path(), "admin") {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "invalid scope",
		})

	}

	if user.IsAmbassador {
		scope = "ambassador"
	} else {
		scope = "admin"
	}

	token, err := middleware.GenerateJWT(user.Id, scope)
	if err != nil {
		return err
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().AddDate(0, 0, 7),
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

	if strings.Contains(c.Path(), "/api/ambassador") {
		ambassador := models.Ambassador(user)
		ambassador.Revenue = ambassador.CalculateTotalRevenue(database.DB)
		return c.JSON(ambassador)
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

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	userId, _ := middleware.GetUserId(c)

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}
	user.Id = userId

	database.DB.Model(&user).Updates(&user)

	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {

		c.Status(fiber.StatusBadRequest)
		err := c.JSON(fiber.Map{"message": "passwords doesn't match"})
		if err != nil {
			return err
		}
	}

	userId, _ := middleware.GetUserId(c)

	user := models.User{}
	user.Id = userId
	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(&user)

	return c.JSON(user)
}

func Ambassadors(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Where("is_ambassador = true").Find(&users)

	return c.JSON(users)
}
