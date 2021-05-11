package service

import (
	"errors"
	"regexp"
	"unicode"
	
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
)
//Person type of new user
type Person struct {//struktura użytkownika
	Name      string `json:"name" validate:"required,alpha"`
	Surname   string `json:"surname" validate:"required,alpha"`
	Email     string `json:"email" validate:"required,email"`
	Street    string `json:"street" validate:"required,alpha"`
	Number    string `json:"number" validate:"required,alphanum"`
	City      string `json:"city" validate:"required,alpha"`
	Post_code string `json:"post_code" validate:"required,postCode"`
	Pass      string `json:"pass" validate:"required,password"`
}
//validatePostCode validate PosteCode, correct format: 00-000
func validatePostCode(fl validator.FieldLevel) bool {
	
	re := regexp.MustCompile(`^\d{2}-\d{3}$`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}
//validatePassword validate Password, correct format: min 8 chars, min. 1 Capital letter,min. 1 number
func validatePassword(fl validator.FieldLevel) bool {
	
	var number bool =false
	var upper bool =false
	size:=0
	for _, char := range fl.Field().String() {
		switch {
		case unicode.IsNumber(char): 
			number=true
			size++
		case unicode.IsUpper(char) :
			upper=true
			size++
		default:
			size++
		}
	}
	if (number && upper && size>=8) {
		return true
	} else {
		return false
	}

}
//Validation validate Person struct
func (p *Person) Validation() ([]string, error) {//walidacja struktury użytkownika
	validate := validator.New()
	validate.RegisterValidation("postCode", validatePostCode)
	validate.RegisterValidation("password", validatePassword)
	if err := validate.Struct(p); err != nil {
		log.Info("validation error")
		var TabErrors []string
		for _, err := range err.(validator.ValidationErrors) {

			TabErrors = append(TabErrors, err.Field())
		}
		return TabErrors, errors.New("validation error")
	}
	return make([]string, 0), nil
}