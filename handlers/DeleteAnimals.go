package handlers

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/animals"
	animalsService "github.com/wojbog/praktyki_backend/service/animals"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteAnimal(s *animalsService.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		//extracting user id
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		ownerId := claims["iss"].(string)

		filter := new(models.AnimalFilters)

		if primitiveOwnerId, err := primitive.ObjectIDFromHex(ownerId); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		} else {
			filter.OwnerId = primitiveOwnerId
		}

		if err := c.BodyParser(filter); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{})
		}

		if filter.Series == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "animal must be specified by series number"})
		}

		if err := s.DeleteAnimal(c.Context(), *filter); err != nil {
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
