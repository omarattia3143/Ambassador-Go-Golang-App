package controllers

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/models"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"sort"
	"strconv"
	"strings"
	"time"
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

func ProductsFrontend(c *fiber.Ctx) error {
	var ctx = context.Background()
	var products []models.Product

	result, err := database.Cache.Get(ctx, "products_frontend").Result()
	if err != nil {
		database.DB.Find(&products)
		bytes, err := json.Marshal(products)
		if err != nil {
			return err
		}
		database.Cache.Set(ctx, "products_frontend", bytes, time.Minute*30)
	} else {
		unmarshalErr := json.Unmarshal([]byte(result), &products)
		if unmarshalErr != nil {
			return unmarshalErr
		}
	}

	return c.JSON(products)
}

func ProductsBackend(c *fiber.Ctx) error {
	var ctx = context.Background()
	var products []models.Product

	result, err := database.Cache.Get(ctx, "products_backend").Result()
	if err != nil {
		database.DB.Find(&products)
		bytes, err := json.Marshal(products)
		if err != nil {
			return err
		}
		database.Cache.Set(ctx, "products_backend", bytes, time.Minute*30)
	} else {
		unmarshalErr := json.Unmarshal([]byte(result), &products)
		if unmarshalErr != nil {
			return unmarshalErr
		}
	}

	var searchProducts []models.Product

	//Search
	if s := c.Query("search"); s != "" {
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), strings.ToLower(s)) || strings.Contains(strings.ToLower(product.Description), strings.ToLower(s)) {
				searchProducts = append(searchProducts, product)
			}
		}
	} else {
		searchProducts = products
	}

	//Sort
	if sortParam := c.Query("sort"); sortParam != "" {
		sortParam = strings.ToLower(sortParam)

		if sortParam == "asc" {
			sort.Slice(searchProducts, func(i, j int) bool {
				return searchProducts[i].Price < searchProducts[j].Price
			})
		} else if sortParam == "desc" {
			sort.Slice(searchProducts, func(i, j int) bool {
				return searchProducts[i].Price > searchProducts[j].Price
			})
		}

	}

	//Pagination
	total := len(searchProducts)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage := 9
	data := searchProducts

	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		data = searchProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}

	return c.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": total/perPage + 1,
	})
}
