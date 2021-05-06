package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wojbog/praktyki_backend/user"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/newUser",user.SetUser)//endpoint dodawania uzytkownika do bazy

	app.Listen(":3000")
}
