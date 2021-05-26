package service

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/animals"
	"github.com/wojbog/praktyki_backend/repository/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

//TestAddNewUser test AddNewUser
func TestAddNewUser(t *testing.T) {
	p := []models.NewUser{{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6q", Email: "sss4tefan@elo.pl"},
		{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6q", Email: "sss4tefan@elo.pl"},
		{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "0009000", Pass: "Wojtek6q", Email: "@sss4tefan@elo.pl"}}
	s, c, _ := config()
	_, err := s.AddNewUser(context.Background(), p[0])
	if err != nil {
		t.Error()
	}
	_, err = s.AddNewUser(context.Background(), p[1])
	if err.Error() != "user exists" {
		t.Error()
	}
	_, err = s.AddNewUser(context.Background(), p[2])
	if err.Error() != "validation error" {
		t.Error()
	}
	c.DeleteOne(context.Background(), bson.M{"email": "sss4tefan@elo.pl"})
}

//TestLoginUser test LoginUser
func TestLoginUser(t *testing.T) {
	p := []models.User{{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6q", Email: "sss4tefan@elo.pl"},
		{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6qq", Email: "sss4tefan@elo.pl"},
		{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6q", Email: "ssssss4tefan@elo.pl"}}
	us := p[0]
	s, c, _ := config()
	str, _ := bcrypt.GenerateFromPassword([]byte(us.Pass), 14)
	us.Pass = string(str)
	c.InsertOne(context.Background(), us)

	_, err := s.LoginUser(context.Background(), p[0])
	if err != nil {
		t.Error()
	}
	_, err = s.LoginUser(context.Background(), p[1])
	if err.Error() != "incorrect password" {
		t.Error()
	}
	_, err = s.LoginUser(context.Background(), p[2])
	if err.Error() != "incorrect data" {
		t.Error()
	}
	c.DeleteOne(context.Background(), us)

}

func TestGetAnimals(t *testing.T) {
	s, _, c := config()

	testUserId, _ := primitive.ObjectIDFromHex("0")

	exampleAnimal := models.Animal{
		OwnerId:      testUserId,
		Series:       "abcdef",
		BirthDate:    time.Date(2020, 10, 3, 0, 0, 0, 0, time.UTC),
		Species:      "1",
		UtilityType:  "1",
		Sex:          "1",
		Status:       "1",
		MotherSeries: "1",
		Breed:        "1",
	}

	exampleAnimal2 := models.Animal{
		OwnerId:      testUserId,
		Series:       "123456",
		BirthDate:    time.Date(2010, 12, 3, 0, 0, 0, 0, time.UTC),
		Species:      "2",
		UtilityType:  "2",
		Sex:          "2",
		Status:       "2",
		MotherSeries: "2",
		Breed:        "2",
	}

	exampleAnimal3 := models.Animal{
		OwnerId:      testUserId,
		Series:       "zyxwut",
		BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
		Species:      "3",
		UtilityType:  "3",
		Sex:          "3",
		Status:       "3",
		MotherSeries: "3",
		Breed:        "3",
	}

	c.InsertMany(context.Background(), []interface{}{exampleAnimal, exampleAnimal2, exampleAnimal3})
	defer c.DeleteMany(context.Background(), bson.M{"ownerId": testUserId})

	testCases := []models.AnimalFilters{
		{Series: "123456"},
		{
			OwnerId:      testUserId,
			MinBirthDate: "2006-10-10",
			MaxBirthDate: "2030-10-10",
		},
		{
			OwnerId:      testUserId,
			MinBirthDate: "sie",
			MaxBirthDate: "ma co tam",
		},
		{
			OwnerId:      testUserId,
			MinBirthDate: "ugh hackers everywhere",
			MaxBirthDate: "2005-12-19",
		},
		{
			OwnerId:      testUserId,
			MinBirthDate: "2007-01-01",
			MaxBirthDate: "im like an analphabetic",
		},
		{OwnerId: testUserId, MotherSeries: "3"},
		{OwnerId: testUserId, Sex: "2"},
		{OwnerId: testUserId, Breed: "chleb"},
		{OwnerId: testUserId, UtilityType: "coco jambo"},
	}

	exp := []struct {
		animals []models.Animal
		err     error
	}{
		{[]models.Animal{exampleAnimal2}, nil},
		{[]models.Animal{exampleAnimal, exampleAnimal2}, nil},
		{[]models.Animal{exampleAnimal, exampleAnimal2, exampleAnimal3}, nil},
		{[]models.Animal{exampleAnimal3}, nil},
		{[]models.Animal{exampleAnimal, exampleAnimal2}, nil},
		{[]models.Animal{exampleAnimal3}, nil},
		{[]models.Animal{exampleAnimal2}, nil},
		{nil, nil},
		{nil, nil},
	}

	if len(testCases) != len(exp) {
		t.Fatal("Numbers of test cases and expectations are not equal!")
	}

	for i := 0; i < len(testCases); i++ {
		if res, err := s.GetAnimals(context.Background(), testCases[i]); err != nil {
			t.Errorf("Unespected error occured: %s", err)
		} else {
			if !reflect.DeepEqual(res, exp[i].animals) {
				t.Errorf("Wrong query result!\nFilter: %+v\nExpected: %v\nReceived: %v", testCases[i], exp[i].animals, res)
			}
			if err != exp[i].err {
				t.Errorf("Wrong error response!\nFilter: %+v\nExpected: %v\nReceived: %v", testCases[i], exp[i].err, err)
			}
		}
	}
}

func config() (*Service, *mongo.Collection, *mongo.Collection) {
	str1 := os.Getenv("MONGO_URL")
	if str1 == "" {
		log.Fatal("NO MONGO URL")
	}
	str2 := os.Getenv("MONGO_DB")
	if str2 == "" {
		log.Fatal("NO MONGO DB")
	}
	mongo_URL := str1
	datebasename := str2

	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_URL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("connect to DB")
	}
	DatabseName := datebasename
	db := client.Database(DatabseName)

	ucol := *db.Collection("users")
	userCol := user.NewCollection(&ucol)

	acol := *db.Collection("animals")
	animalCol := animals.NewCollection(&acol)

	s := NewService(userCol, animalCol)

	return s, &ucol, &acol
}
