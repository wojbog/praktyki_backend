package handlers

import (
	"errors"
	"os"
	"time"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/service"
)
//LoginUser endpoint login
func LoginUser(s *service.Service) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {

		var p models.UserLogin

		//convert to models.UserLogin type
		if err := c.BodyParser(&p); err != nil {
			log.Info(err.Error())
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"error":   "Cannot read request"})
		}
		//check in datebase and create token
		if user, err := s.LoginUser(c.Context(), p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"success": false,
				"error":   "incorrect password or email"})
		} else {
			if token, err := CreateToken(user.Id.Hex()); err != nil {
				return c.Status(500).JSON(&fiber.Map{
					"success": false,
					"error":   "Internal Server Error"})
			} else {
				return c.Status(200).JSON(&fiber.Map{
					"success": true,
					"token":   token})
			}
		}

	}
}
//CreateToken return token if no errors
func CreateToken(id string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix()})
	if secret := os.Getenv("ACCESS_SECRET"); secret == "" {
		log.Fatal("No ACCESS_SECRET")
		return "", errors.New("")
	} else {
		token, err := claims.SignedString([]byte(secret))
		if err != nil {
			return "", err
		}
		return token, nil
	}

}
