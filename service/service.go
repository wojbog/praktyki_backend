package service

import (
	"context"
	"errors"
	"regexp"
	"unicode"
	"github.com/wojbog/praktyki_backend/repository/user"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	log "github.com/sirupsen/logrus"
)


//Service store collection
type Service struct {
	userCol *user.Collection
}

//AddNewUser
//return status,userResponse,table of errors,error 
func (s *Service) AddNewUser(ctx context.Context, user models.NewUser)(int,models.UserResponse,[]string,error ) {

	//validation
	if tab, errv := Validate(user); errv != nil {
		log.Info(errv.Error())
		return 400,models.UserResponse{},tab, errv
	}
	
	//hash
	str,_:=bcrypt.GenerateFromPassword([]byte(user.Pass),14)
	user.Pass=string(str)

	//add to datebase
	if user,err:=s.userCol.InsertUser(ctx,user);err!=nil {
		log.Info(err.Error())
		if err.Error()=="user exists" {
			return 400,models.UserResponse{},[]string{},err
		} else {
			return 500,models.UserResponse{},[]string{},errors.New("Internal Server Error")
		}
	}else {
		return 200,user,[]string{},nil
	}
	

}

//NewService create new service
func NewService(col *user.Collection) *Service  {
	return &Service{col}
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

	var number bool = false
	var upper bool = false
	size := 0
	for _, char := range fl.Field().String() {
		switch {
		case unicode.IsNumber(char):
			number = true
			size++
		case unicode.IsUpper(char):
			upper = true
			size++
		default:
			size++
		}
	}
	if number && upper && size >= 8 {
		return true
	} else {
		return false
	}

}

//Validation validate Person struct
func Validate(p models.NewUser) ([]string, error) {
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
