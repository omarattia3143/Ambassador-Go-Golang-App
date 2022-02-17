package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	_, err := gorm.Open(mysql.Open("root:root@tcp(db:3306)/ambassador"), &gorm.Config{})
	if err != nil {
		panic("couldn't connect to database!")
	}

	log.Println("connected to database successfully!")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	log.Fatal(app.Listen(":3000"))
}
