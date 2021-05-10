package main

import (
	"os"
	
	"github.com/gofiber/fiber/v2"
	"github.com/wojbog/praktyki_backend/user"
	log "github.com/sirupsen/logrus"
)

func main() {
	
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/createUser",user.SetUser)//endpoint dodawania uzytkownika do bazy
	
	PORT:=os.Getenv("PORT")
	if PORT != "" {
		app.Listen(":"+PORT)
	} else {
		log.Panic("NO PORT")
	}
	
	
}
