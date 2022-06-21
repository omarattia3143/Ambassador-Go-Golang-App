package controllers

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/middleware"
	"GoAndNextProject/src/models"
	"github.com/bxcodec/faker/v3"
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

type CreateLinkRequest struct {
	Products []int
}

func CreateLink(c *fiber.Ctx) error {
	var request CreateLinkRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	id, _ := middleware.GetUserId(c)

	link := models.Link{
		Code:   faker.Username(),
		UserId: id,
	}

	for productId := range request.Products {
		product := models.Product{
			Model: models.Model{Id: uint(productId)},
		}
		link.Products = append(link.Products, product)
	}

	database.DB.Create(&link)

	return c.JSON(link)
}

func Stats(c *fiber.Ctx) error {
	id, _ := middleware.GetUserId(c)
	var links []models.Link

	database.DB.Find(&links, models.Link{
		UserId: id,
	})

	var result []any
	var orders []models.Order

	for _, link := range links {
		database.DB.Preload("OrderItems").Find(&orders, models.Order{
			Code:     link.Code,
			Complete: true,
		})

		revenue := 0.0
		for _, order := range orders {
			revenue += order.GetTotal()
		}

		result = append(result, fiber.Map{
			"code":    link.Code,
			"count":   len(orders),
			"revenue": revenue,
		})
	}
	return c.JSON(result)
}

func Rankings(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Find(&users, models.User{
		IsAmbassador: true,
	})

	var result []any

	for _, user := range users {
		var ambassador = models.Ambassador(user)
		result = append(result, fiber.Map{
			user.FullName(): ambassador.CalculateTotalRevenue(database.DB),
		})
	}
	return c.JSON(result)
}
