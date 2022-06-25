package controllers

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
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

type CreateOrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	link := models.Link{}

	database.DB.Preload("User").First(&link, models.Link{
		Code: request.Code,
	})

	if link.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Link",
		})
	}

	order := models.Order{
		Code:            link.Code,
		UserId:          link.UserId,
		AmbassadorEmail: link.User.Email,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Address:         request.Address,
		Country:         request.Country,
		City:            request.City,
		Zip:             request.Zip,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&order); err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, requestProduct := range request.Products {
		product := models.Product{}
		product.Id = uint(requestProduct["product_id"])
		database.DB.First(&product)

		total := product.Price * float64(requestProduct["quantity"])

		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(requestProduct["quantity"]),
			AdminRevenue:      0.9 * total,
			AmbassadorRevenue: 0.1 * total,
		}

		if err := tx.Create(&item); err != nil {
			tx.Rollback()
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": err,
			})
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(product.Title),
			Description: stripe.String(product.Description),
			Images:      []*string{stripe.String(product.Image)},
			Amount:      stripe.Int64(100 * int64(product.Price)),
			Currency:    stripe.String("usd"),
			Quantity:    stripe.Int64(int64(requestProduct["quantity"])),
		})
	}

	stripe.Key = "sk_test_51LELzyK9sY3PZhfw2JazYdwBdQlQVsOaEf4LnKXw9dRl2IE9sxhNeD8GiKsx0pVnF9EElfPXpL3meJYQyC3HZlNt00NXJg8tp9"

	params := stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String("http://localhost:5000/success?source={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("http://localhost:5000/error"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
	}

	source, err := session.New(&params)

	if err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	order.TransactionId = source.ID

	if err := tx.Save(&order); err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	tx.Commit()

	return c.JSON(source)
}

func CompleteOrder(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	order := models.Order{}

	database.DB.Preload("OrderItems").First(&order, &models.Order{
		TransactionId: data["source"],
	})

	if order.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Order not found",
		})
	}

	order.Complete = true
	database.DB.Save(&order)

	go func(order models.Order) {
		ambssadorRevenue := 0.0
		adminRevenue := 0.0

		for _, item := range order.OrderItems {
			ambssadorRevenue += item.AmbassadorRevenue
			adminRevenue += item.AdminRevenue
		}
		user := models.User{}
		user.Id = order.UserId

		database.DB.First(&user)

		database.Cache.ZIncrBy(context.Background(), "rankings", ambssadorRevenue, user.FullName())
	}(order)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
