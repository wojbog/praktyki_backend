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
//return userResponse,error 
func (s *Service) AddNewUser(ctx context.Context, user models.NewUser)(models.UserResponse,error) {

	//validation
	if  errv := Validate(user); errv != nil {
		log.Info(errv.Error())
		return models.UserResponse{}, errv
	}
	
	//hash
	str,_:=bcrypt.GenerateFromPassword([]byte(user.Pass),14)
	user.Pass=string(str)

	//add to datebase
	if user,err:=s.userCol.InsertUser(ctx,user);err!=nil {
		log.Info(err.Error())
		if err.Error()=="user exists" {
			return models.UserResponse{},err
		} else {
			return models.UserResponse{},errors.New("Internal Server Error")
		}
	}else {
		return user,nil
	}
	

}
//LoginUser
//return models.UserLogin,error
func (s *Service) LoginUser(ctx context.Context, user models.UserLogin)(models.UserLogin,error) {

		//chek in datebase
		if userDB,err:=s.userCol.GetUserLogin(ctx,user);err!=nil {
			log.Info(err.Error())
			return models.UserLogin{},err
			
		}else {
			if errv := bcrypt.CompareHashAndPassword([]byte(userDB.Pass), []byte(user.Pass)); errv!=nil {
				log.Info("incorrect password user: "+userDB.Id.Hex())
				return models.UserLogin{},errors.New("incorrect password")
			}else {
				return user,nil
			}
			
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
func Validate(p models.NewUser) error {
	validate := validator.New()
	validate.RegisterValidation("postCode", validatePostCode)
	validate.RegisterValidation("password", validatePassword)
	if err := validate.Struct(p); err != nil {
		var TabErrors string
		for _, err := range err.(validator.ValidationErrors) {
			TabErrors+=" "+err.Field()
			
		}
		return  errors.New("validation-error:"+TabErrors)
	}
	return  nil
}
