package service

import (
	"testing"

	"github.com/go-playground/validator"
	"github.com/wojbog/praktyki_backend/models"
)

//TestValidatePasswordReturnErrorIfPasswordIsValid
//get incorrect data
//return error if func return true
func TestValidatePasswordReturnErrorIfPasswordIsValid(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("password", validatePassword)
	arr := []string{"Wojtekir", "asdd", "asdasdaq", "12345678", "Wojtekqw+_';\"["}
	for i := 0; i < 5; i++ {
		err := validate.Var(arr[i], "password")
		if err == nil {
			t.Error(err)
		}
	}

}

//TestValidatePasswordReturnErrorIfPasswordIsInvalid
//get correct data
//return error if func return false
func TestValidatePasswordReturnErrorIfPasswordIsInvalid(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("password", validatePassword)
	arr := []string{"Wojtek9r", "asDddfgdfg887", "asdasdaqE4", "12345678Q", "Wojtekqw4"}
	for i := 0; i < 5; i++ {
		err := validate.Var(arr[i], "password")
		if err != nil {
			t.Error(err)
		}
	}

}

//TestValidatePostReturnErrorIfPostCodeIsValid
//get incorrect data
//return error if func return true
func TestValidatePostReturnErrorIfPostCodeIsValid(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("postCode", validatePostCode)
	arr := []string{"aw-asd", "123435", "234-87", "34-876h", "32y768"}
	for i := 0; i < 5; i++ {
		err := validate.Var(arr[i], "postCode")
		if err == nil {
			t.Error(err)
		}
	}

}

//TestValidatePostReturnErrorIfPostCodeIsInvalid
//get correct data
//return error if func return false
func TestValidatePostReturnErrorIfPostCodeIsInvalid(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("postCode", validatePostCode)
	arr := []string{"12-678", "87-456", "09-000", "00-000", "99-999"}
	for i := 0; i < 5; i++ {
		err := validate.Var(arr[i], "postCode")
		if err != nil {
			t.Error(err)
		}
	}

}

//TestValidateStructReturnErrorIfStructIsInvalid
//get incorrect data
//return error if func return true
func TestValidateStructReturnErrorIfStructIsInvalid(t *testing.T) {
	per := []models.NewUser{{Name: "Andasd", Surname: "asdasdP", Email: "asd.pooijk-asda@asdasd.pl", Street: "asdasd", Number: "12a", City: "ssdfsdfQWEWEQ", Post_code: "45-456", Pass: "Wertyui9"},
		{Name: "ASDSSD", Surname: "asdasd", Email: "asd.pooijk-asda@asdasd.asdasd.asd-asd.pl", Street: "ASADDSA", Number: "12243", City: "WEQQWEasdasd", Post_code: "00-000", Pass: "W+-==9tt"}}
	validate := validator.New()
	validate.RegisterValidation("postCode", validatePostCode)
	validate.RegisterValidation("password", validatePassword)
	for i := 0; i < len(per); i++ {
		err := Validate(per[i])
		if err != nil {
			t.Error(err)
		}
	}

}

//TestValidateStructReturnErrorIfStructIsValid
//get correct data
//return error if func return false
func TestValidateStructReturnErrorIfStructIsValid(t *testing.T) {
	per := []models.NewUser{{Name: "Andasd", Surname: "asdasdP", Email: "asd.pooijk-asda@asdasd.pl@", Street: "asdasd", Number: "12a", City: "ssdfsdfQWEWEQ", Post_code: "45-456", Pass: "Wertyui9"},
		{Name: "ASDSSD3", Surname: "asdasd", Email: "asd.pooijk-asda@asdasd.asdasd.asd-asd.pl", Street: "ASADDSA", Number: "12243", City: "WEQQWEasdasd", Post_code: "00-000", Pass: "W+-==9tt"},
		{Name: "ASDSSDh", Surname: "asdas9d", Email: "asd.pooijk-asda@asdasd.asdasd.asd-asd.pl", Street: "ASADDSA", Number: "12243", City: "WEQQWEasdasd", Post_code: "00-000", Pass: "W+-==9tt"},
		{Name: "ASDSSDl", Surname: "asdasd", Email: "asd.pooijk-asda@asdasd.asdasd.asd-asd.pl", Street: "ASADDSA", Number: "----", City: "WEQQWEasdasd", Post_code: "00-000", Pass: "W+-==9tt"},
		{Name: "ASDSSDl", Surname: "asdasd", Email: "asd.pooijk-asda@asdasd.asdasd.asd-asd.pl", Street: "ASADDSA", Number: "12243", City: "WEQQWEa1441sdasd", Post_code: "00-000", Pass: "W+-==9tt"},
		{Name: "ASDSSDl", Surname: "asdasd", Email: "asd.pooijk-asda@asdasd.asdasd.asd-asd.pl", Street: "ASADDSA", Number: "12243", City: "WEQQWEasdasd", Post_code: "000000", Pass: "W+-==9tt"},
		{Name: "ASDSSDl", Surname: "asdasd", Email: "asd.pooijk-asda@asdasd.asdasd.asd-asd.pl", Street: "ASADDSA", Number: "12243", City: "WEQQWEasdasd", Post_code: "00-000", Pass: "W+-=="}}
	for i := 0; i < len(per); i++ {
		err := Validate(per[i])
		if err == nil {
			t.Error(err)
		}
	}

}
