package handlers

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GetAnimals handle getAnimal request
//It parses user query to models.AnimalFilters and passes that to service
func GetAnimals(s *service.Service) func(c *fiber.Ctx) error {
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

		if err := c.QueryParser(filter); err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{})
		}

		animals, err := s.GetAnimals(c.Context(), *filter)
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{})
		}

		return c.Status(200).JSON(&fiber.Map{
			"animals": animals,
		})

	}
}
