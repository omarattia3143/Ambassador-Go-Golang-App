package main

import (
	"GoAndNextProject/src/database"
	"GoAndNextProject/src/routes"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	database.Connect()
	database.AutoMigrate()

	app := fiber.New()
	routes.Setup(app)

	log.Fatal(app.Listen(":8000"))
}
