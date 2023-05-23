package apis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api")
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", swagger.HandlerDefault)

	router := app.Group("/api")

	// meta
	router.Get("/meta", GetMeta)

	// user
	router.Post("/register", Register)
	router.Post("/login", Login)
	router.Get("/users/me", GetUserMe)
	router.Patch("/users/me", ModifyUserMe)
	router.Delete("/users/me", DeleteUserMe)
	router.Get("/users", ListUsers)
	router.Get("/users/:id", GetUser)
	router.Patch("/users/:id", ModifyAUser)
	router.Delete("/users/:id", DeleteAUser)

	// book
	router.Get("/books", ListBooks)
	router.Post("/books", CreateABook)
	router.Patch("/books/:id", ModifyABook)

	// purchase
	router.Get("/purchases", ListPurchases)
	router.Get("/purchases/:id", GetAPurchase)
	router.Post("/purchases", CreateAPurchase)
	router.Patch("/purchases/:id", ModifyAPurchase)
	router.Post("/purchases/:id/_pay", PayAPurchase)
	router.Post("/purchases/:id/_return", ReturnAPurchase)
	router.Post("/purchases/:id/_arrive", ArriveAPurchase)

	// balance
	router.Get("/balances", ListBalances)
	router.Get("/balances/:id", GetABalance)
	router.Post("/balances", CreateABalance)

	// sale
	router.Get("/sales", ListSales)
	router.Get("/sales/:id", GetASale)
	router.Post("/sales", CreateASale)
}
