package service

import (
	"context"
	"errors"
	"regexp"
	"time"
	"unicode"

	"github.com/go-playground/validator"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/animals"
	"github.com/wojbog/praktyki_backend/repository/user"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

//Service store collection
type Service struct {
	userCol    *user.Collection
	animalsCol *animals.Collection
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

//GetAnimals returns array of filter.ownerId's animals.
//It converts models.AnimalFilters to compatible with db map with filters
//If no animals return null
func (s *Service) GetAnimals(ctx context.Context, filter models.AnimalFilters) ([]models.Animal, error) {
	const layoutISO = "2006-01-02"

	maxDate, err := time.Parse(layoutISO, filter.MaxBirthDate)
	if err != nil {
		maxDate, _ = time.Parse(layoutISO, "9999-12-31")
	}
	minDate, _ := time.Parse(layoutISO, filter.MinBirthDate)

	filter.MinBirthDate = ""
	filter.MaxBirthDate = ""

	mapFilter := make(map[string]interface{})

	if err := mapstructure.Decode(filter, &mapFilter); err != nil {
		return nil, err
	}

	mapFilter["birthDate"] = bson.M{
		"$gte": minDate,
		"$lt":  maxDate,
	}

	if animals, err := s.animalsCol.GetAnimals(ctx, mapFilter); err != nil {
		log.Info(err.Error())
		return nil, err
	} else {
		return animals, nil
	}
}

//NewService create new service
func NewService(userCol *user.Collection, animalsCol *animals.Collection) *Service {
	return &Service{userCol, animalsCol}
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
