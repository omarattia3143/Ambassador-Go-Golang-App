package routes

import (
	"GoAndNextProject/src/controllers"
	"GoAndNextProject/src/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Get("/", controllers.Home)

	api := app.Group("api")

	//admin
	admin := api.Group("admin")
	admin.Post("register", controllers.Register)
	admin.Post("login", controllers.Login)

	//admin + middleware
	adminAuthenticated := admin.Use(middleware.IsAuthenticated)
	adminAuthenticated.Get("user", controllers.User)
	adminAuthenticated.Post("logout", controllers.Logout)
	adminAuthenticated.Put("user/info", controllers.UpdateInfo)
	adminAuthenticated.Put("user/password", controllers.UpdatePassword)
	adminAuthenticated.Get("ambassadors", controllers.Ambassadors)
	//products
	adminAuthenticated.Get("products", controllers.Products)
	adminAuthenticated.Post("products", controllers.CreateProduct)
	adminAuthenticated.Get("products/:id", controllers.GetProduct)
	adminAuthenticated.Put("products/:id", controllers.UpdateProduct)
	adminAuthenticated.Delete("products/:id", controllers.DeleteProduct)
	//Links
	adminAuthenticated.Get("user/:id/links", controllers.Links)
	//orders
	adminAuthenticated.Get("orders", controllers.Orders)

	//----------------------------------------------------------------

	ambassador := api.Group("ambassador")
	ambassador.Post("register", controllers.Register)
	ambassador.Post("login", controllers.Login)
	ambassador.Get("products/frontend", controllers.ProductsFrontend)
	ambassador.Get("products/backend", controllers.ProductsBackend)

	//admin + middleware
	ambassadorAuthenticated := ambassador.Use(middleware.IsAuthenticated)
	ambassadorAuthenticated.Get("user", controllers.User)
	ambassadorAuthenticated.Post("logout", controllers.Logout)
	ambassadorAuthenticated.Put("user/info", controllers.UpdateInfo)
	ambassadorAuthenticated.Put("user/password", controllers.UpdatePassword)
	ambassadorAuthenticated.Get("ambassadors", controllers.Ambassadors)
	ambassadorAuthenticated.Post("links", controllers.CreateLink)
	ambassadorAuthenticated.Get("stats", controllers.Stats)
	ambassadorAuthenticated.Get("rankings", controllers.Rankings)

	//Checkout
	checkout := api.Group("checkout")
	checkout.Get("links/:code", controllers.GetLink)
	checkout.Post("orders", controllers.CreateOrder)
	checkout.Post("orders/confirm", controllers.CompleteOrder)
}
