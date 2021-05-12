//server route in app
package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wojbog/praktyki_backend/handlers"
	"github.com/wojbog/praktyki_backend/service"
)


//Routing routes in app
func Routing(app *fiber.App, ser *service.Service){

	
	//CreateUser endpoint add new user
	app.Post("/createUser", handlers.CreateUser(ser))


}
