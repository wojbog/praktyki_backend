//handlers in pplication
package handlers


import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/service"
	"github.com/wojbog/praktyki_backend/repository/user"
)
//CreateUser add new user to datebase
func CreateUser(col *user.Collection) func (c *fiber.Ctx) error {

	return func (c *fiber.Ctx) error {
	
	//p instance of user
	p := new(service.Person)

	//convert to service.Person type 
	if err := c.BodyParser(p); err != nil {
		log.Info(err.Error())
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   err.Error()})
	}
	//validation
	if tab, errv := p.Validation(); errv != nil {
		log.Info(errv.Error())
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   errv.Error(),
			"list":    tab})
	}
	//add to datebase
	if str,err:=col.InsertUser(c.Context(),p);err!=nil {//dodawanie do bazy
		log.Info(err.Error())
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error": err.Error()})
	} else {
		log.Info("success add new user, id: " + str)
		return c.Status(200).JSON(&fiber.Map{"success": true, "id": str})
	}	
	

}}
