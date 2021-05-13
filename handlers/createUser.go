//handlers in pplication
package handlers

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/service"
)

//CreateUser add new user to datebase
func CreateUser(s *service.Service) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {

		//p instance of user
		var p models.NewUser

		//convert to models.NewUser type
		if err := c.BodyParser(&p); err != nil {
			log.Info(err.Error())
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"error":   "Cannot read request"})
		}

		//add to datebase
		if user, err := s.AddNewUser(c.Context(), p); err != nil {
			if err.Error() == "user exists" {
				return c.Status(400).JSON(&fiber.Map{
					"success": false,
					"error":   err})
			} else if err.Error() == "validation error" {
				return c.Status(400).JSON(&fiber.Map{
					"success": false,
					"error":   err})

			} else {
				return c.Status(500).JSON(&fiber.Map{
					"success": false,
					"error":   err})
			}
		} else {

			return c.Status(200).JSON(&fiber.Map{"success": true, "user": user})
		}

	}
}
