//server route in app
package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wojbog/praktyki_backend/handlers"
	"github.com/wojbog/praktyki_backend/repository/user"
)


//Routing routes in app
func Routing(app *fiber.App, col *user.Collection){

	
	//CreateUser endpoint add new user
	app.Post("/createUser", handlers.CreateUser(col))


}
