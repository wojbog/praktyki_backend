//handlers in pplication
package handlers


import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/service"
)

//CreateUser add new user to datebase
func CreateUser(s *service.Service) func (c *fiber.Ctx) error {

	return func (c *fiber.Ctx) error {
	
	//p instance of user
	var p service.Person

	//convert to service.Person type 
	if err := c.BodyParser(&p); err != nil {
		log.Info(err.Error())
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   err.Error()})
	}
	//add to datebase
	if tab, err:=s.AddNewUser(c.Context(),p); err!=nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   err.Error(),
			"list":    tab})
	} else {
		log.Info("success add new user, id: " + tab[0])
		return c.Status(200).JSON(&fiber.Map{"success": true, "id": tab[0]})
	}

}}
