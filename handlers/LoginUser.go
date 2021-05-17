package handlers

import (
	"errors"
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/user"
	"github.com/wojbog/praktyki_backend/service"
)

//LoginUser endpoint login
func LoginUser(s *service.Service) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {

		var p models.User

		//convert to models.UserLogin type
		if err := c.BodyParser(&p); err != nil {
			log.Info(err.Error())
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"error":   "cannot read request"})
		}
		//check in datebase and create token
		if us, err := s.LoginUser(c.Context(), p); err != nil {
			if err == service.IncorrectPasswordError || err == user.UserNotFoundError {
				return c.Status(400).JSON(&fiber.Map{
					"success": false,
					"error":   "incorrect password or email"})
			} else {
				return c.Status(500).JSON(&fiber.Map{
					"success": false,
					"error":   "internal Server Error"})
			}
		} else {
			if token, err := CreateToken(us.Id.Hex()); err != nil {
				return c.Status(500).JSON(&fiber.Map{
					"success": false,
					"error":   "internal Server Error"})
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
		log.Error("No ACCESS_SECRET")
		return "", errors.New("NO ACCESS_SECRET")
	} else {
		token, err := claims.SignedString([]byte(secret))
		if err != nil {
			return "", err
		}
		return token, nil
	}

}
