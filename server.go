package main

import (
	"os"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/server"

)

func main() {
	
	app := fiber.New()
	
	server.Routing(app)
	
	
	PORT:=os.Getenv("PORT")
	if PORT != "" {
		app.Listen(":"+PORT)
	} else {
		log.Panic("NO PORT")
	}
	
	
}
