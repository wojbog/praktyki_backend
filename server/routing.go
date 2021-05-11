package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wojbog/praktyki_backend/handlers"
	
)

func Routing(app *fiber.App){

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/createUser", handlers.CreateUser)//endpoint dodawania uzytkownika do bazy
}
