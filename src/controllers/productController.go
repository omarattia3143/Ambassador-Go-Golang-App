package controllers

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"github.com/gofiber/fiber/v2"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product

	database.DB.Find(&products)

	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {

	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	var exists models.Product
	database.DB.Where("title = ?", product.Title).First(&exists)

	if exists.Id != 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "product already exists",
		})
	}

	database.DB.Create(&product)

	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")

	var product models.Product
	product.Id = uint(id)

	database.DB.Find(&product)

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")

	product := models.Product{}
	product.Id = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product)

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")

	product := models.Product{}
	product.Id = uint(id)

	database.DB.Delete(&product)

	return c.JSON(fiber.Map{
		"message": "Deleted successfully",
	})
}
