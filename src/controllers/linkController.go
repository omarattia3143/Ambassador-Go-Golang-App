package controllers

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func Links(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var links []models.Link

	database.DB.Where("user_id = ?", id).Find(&links)

	for i, link := range links {
		var orders []models.Order
		database.DB.Where("Code = ? and complete = true", link.Code).Find(&orders)
		links[i].Orders = orders
	}

	return c.JSON(links)
}
