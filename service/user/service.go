package service

import (
	"context"
	"errors"
	"regexp"
	"unicode"

	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/user"
	"golang.org/x/crypto/bcrypt"
)

//Service store collection
type Service struct {
	userCol *user.Collection
}

//NewService create new service
func NewService(userCol *user.Collection) *Service {
	return &Service{userCol}
}

//ValidationError contains array of invalid fields
type ValidationError struct {
	InvalidFields []string
}

func (e *ValidationError) Error() string {
	return "validation error"
}

var IncorrectPasswordError = errors.New("incorrect password")

//AddNewUser
//return userResponse,error
func (s *Service) AddNewUser(ctx context.Context, userNew models.NewUser) (models.UserResponse, error) {

	//validation
	if errv := Validate(userNew); errv != nil {
		log.Info(errv.Error())
		return models.UserResponse{}, errv
	}

	//hash
	str, _ := bcrypt.GenerateFromPassword([]byte(userNew.Pass), 14)
	userNew.Pass = string(str)

	//add to datebase
	if us, err := s.userCol.InsertUser(ctx, userNew); err != nil {
		log.Info(err.Error())
		if err == user.UserExistError {
			return models.UserResponse{}, err
		} else {
			return models.UserResponse{}, errors.New("internal Server Error")
		}
	} else {
		return us, nil
	}

}

//LoginUser
//return models.UserLogin,error
func (s *Service) LoginUser(ctx context.Context, user models.User) (models.User, error) {

	//chek in datebase
	if userDB, err := s.userCol.GetUserByEmail(ctx, user); err != nil {
		log.Info(err.Error())
		return models.User{}, err

	} else {
		if errv := bcrypt.CompareHashAndPassword([]byte(userDB.Pass), []byte(user.Pass)); errv != nil {
			log.Info("incorrect password user: " + userDB.Id.Hex())

			return models.User{}, IncorrectPasswordError
		} else {
			return userDB, nil
		}
	}
}

//validatePostCode validate PosteCode, correct format: 00-000
func validatePostCode(fl validator.FieldLevel) bool {

	re := regexp.MustCompile(`^\d{2}-\d{3}$`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1

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
		var TabErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			TabErrors = append(TabErrors, err.Field())
		}
		return &ValidationError{InvalidFields: TabErrors}
	}
	return nil
}
