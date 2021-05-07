package main

import (
	"os"
	
	"github.com/gofiber/fiber/v2"
	"github.com/wojbog/praktyki_backend/user"
	"github.com/wojbog/praktyki_backend/baza"
)

func main() {
	baza.ConnectToMongo()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/createUser",user.SetUser)//endpoint dodawania uzytkownika do bazy
	
	app.Listen(":"+os.Getenv("PORT"))
	
}
