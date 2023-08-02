package handlers

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	animalsService "github.com/wojbog/praktyki_backend/service/animals"
)

func CreateAnimal(s *animalsService.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		//extracting user id
		usr := c.Locals("user").(*jwt.Token)
		claims := usr.Claims.(jwt.MapClaims)
		ownerId := claims["iss"].(string)

		//a instance of animal
		animalReq := new(models.AnimalRequest)

		if err := c.BodyParser(animalReq); err != nil {
			log.Info(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "Cannot read request"})
		}

		animal, err := models.Request2Animal(*animalReq, ownerId)
		if err != nil {
			if err.Error() == "cannot parse date" {
				log.Info(err.Error())
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"success": false,
					"error":   "cannot parse date"})
			}
			if err.Error() == "cannot convert ownerId" {
				log.Info(err.Error())
				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"error": "Unauthorized"})
			}
		}
		//pass animals to service function
		a, err := s.AddNewAnimal(c.Context(), animal)
		if err != nil {
			if err.Error() == "validation error" {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"success":   false,
					"error":     err.Error(),
					"errorBody": err})
			}

			if err.Error() == "animal exists" {
				return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
					"success": false,
					"error":   err.Error()})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"error":   err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
			"animal":  a,
		})
	}
}
