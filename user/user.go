package user
/*
package użytkjownika:
-dodawanie do bazy
*/
import (
	
	"fmt"
	
	"github.com/gofiber/fiber/v2"
	"github.com/wojbog/praktyki_backend/serviceDB"
	"github.com/wojbog/praktyki_backend/person"
)

func SetUser(c *fiber.Ctx) error {//dodawanie użytkownika do bazy
	fmt.Println("jestem")

	p := new(person.Person)

	if err := c.BodyParser(p); err != nil {//zamiana na strukturę
		fmt.Println(err)
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   err})
	}

	if tab, errv := p.Validation(); errv != nil {//walidacja
		fmt.Println(errv)
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   "validation error",
			"list":    tab})
	}

	if err:=serviceDB.InsertUser(p);err!=nil {//dodawanie do bazy

		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error": err.Error()})
	}

	fmt.Println(p.Name)
	fmt.Println(p.Post_code)
	fmt.Println(p)

	return c.Status(200).JSON(&fiber.Map{"success": true})

}
