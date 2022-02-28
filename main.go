package main

import (
	"GoAndNextProject/src/database"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	database.Connect()
	database.AutoMigrate()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	log.Fatal(app.Listen(":3000"))
}
