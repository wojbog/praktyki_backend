package handlers
/*
package użytkjownika:
-dodawanie do bazy
*/
import (
	
	
	log "github.com/sirupsen/logrus"
	"github.com/gofiber/fiber/v2"
	"github.com/wojbog/praktyki_backend/repository/user"
	"github.com/wojbog/praktyki_backend/repository/person"
)

func CreateUser(c *fiber.Ctx) error {//dodawanie użytkownika do bazy
	

	p := new(person.Person)

	if err := c.BodyParser(p); err != nil {//zamiana na strukturę
		log.Info(err.Error())
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   err.Error()})
	}

	if tab, errv := p.Validation(); errv != nil {//walidacja
		log.Info(errv.Error())
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   errv.Error(),
			"list":    tab})
	}

	if str,err:=user.InsertUser(c.Context(),p);err!=nil {//dodawanie do bazy
		log.Info(err.Error())
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error": err.Error()})
	} else {
		log.Info("success add new user, id: " + str)
		return c.Status(200).JSON(&fiber.Map{"success": true, "id": str})
	}	
	

}
