package routes

import (
	"GoAndNextProject/src/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Get("/", controllers.Home)

	api := app.Group("api")

	admin := api.Group("admin")
	admin.Post("register", controllers.Register)
	admin.Post("login", controllers.Login)
	admin.Get("user", controllers.User)
	admin.Post("logout", controllers.Logout)

}
