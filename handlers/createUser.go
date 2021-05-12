//handlers in pplication
package handlers


import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/service"
	"github.com/wojbog/praktyki_backend/models"
)

//CreateUser add new user to datebase
func CreateUser(s *service.Service) func (c *fiber.Ctx) error {

	return func (c *fiber.Ctx) error {
	
	//p instance of user
	var p models.NewUser

	//convert to service.Person type 
	if err := c.BodyParser(&p); err != nil {
		log.Info(err.Error())
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   "Cannot read request"})
	}
	//add to datebase
	if status,user,tab, err:=s.AddNewUser(c.Context(),p); err!=nil {
		if status==400 {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   err.Error(),
			"list":    tab})
		}else {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   err.Error()})
		}
	} else {
		
		return c.Status(200).JSON(&fiber.Map{"success": true, "user" : user})
	}

}}
