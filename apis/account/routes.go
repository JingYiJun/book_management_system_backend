package account

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router) {
	router.Post("/register", Register)
	router.Post("/login", Login)
	router.Get("/users/me", GetUserMe)
	router.Patch("/users/me", ModifyUserMe)
	router.Delete("/users/me", DeleteUserMe)
	router.Get("/users", ListUsers)
	router.Get("/users/:id", GetUser)
	router.Patch("/users/:id", ModifyUser)
	router.Delete("/users/:id", DeleteUser)
}
