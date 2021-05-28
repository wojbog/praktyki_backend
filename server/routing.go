//server route in app
package server

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/wojbog/praktyki_backend/handlers"
	animalsService "github.com/wojbog/praktyki_backend/service/animals"
	userService "github.com/wojbog/praktyki_backend/service/user"
)

//Routing routes in app
func Routing(app *fiber.App, uSer *userService.Service, aSer *animalsService.Service) {

	app.Use(cors.New())

	//CreateUser endpoint add new user
	app.Post("/createUser", handlers.CreateUser(uSer))

	//login endpoint login user
	app.Post("/login", handlers.LoginUser(uSer))

	//Protected routes
	secret := os.Getenv("ACCESS_SECRET")
	if secret == "" {
		log.Fatal("No ACCESS_SECRET")
	}

	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
		SigningKey: []byte(secret),
	}))

	//getAnimals endpoint return users animal
	app.Get("/getAnimals", handlers.GetAnimals(aSer))

	//deleteAnimal endpoint delete animal
	app.Post("/deleteAnimal", handlers.DeleteAnimal(aSer))

}
