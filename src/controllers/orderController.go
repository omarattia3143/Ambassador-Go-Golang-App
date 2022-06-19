package controllers

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"github.com/gofiber/fiber/v2"
)

func Orders(c *fiber.Ctx) error {
	var orders []models.Order
	database.DB.Preload("OrderItems").Find(&orders)
	for _, order := range orders {
		order.SetFullName()
		order.GetTotal()
	}
	return c.JSON(orders)
}
