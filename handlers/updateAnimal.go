package handlers

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/animals"
	animalsService "github.com/wojbog/praktyki_backend/service/animals"
)

//UpdateAnimal handle endpoint UpdateAnimal
func UpdateAnimal(s *animalsService.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		//extracting user id
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		ownerId := claims["iss"].(string)

		animalReq := new(models.AnimalRequest)

		if err := c.BodyParser(animalReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{})
		}

		if animal, err := models.Request2Animal(*animalReq, ownerId); err != nil {
			if err.Error() == "cannot convert ownerId" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized",
				})
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "wrong date format",
				})
			}

		} else {
			if animal.Series == "" || animal.Species == "" || animal.UtilityType == "" || animal.Sex == "" || animal.Status == "" || animal.MotherSeries == "" || animal.Breed == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "animal must be specified"})
			}

			if err := s.UpdateAnimal(c.Context(), animal); err != nil {
				if err == animals.AnimalNotexist {
					return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
						"error": err.Error()})
				} else {
					return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{})
				}
			} else {
				return c.Status(200).JSON(&fiber.Map{
					"success": true})
			}

		}
	}
}
