package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateAnimalWithInvalidJwtReturnsUnathorized(t *testing.T) {
	salt := "soli!"
	jwt := getTestToken("qqqqqqwwwwwwr", salt)
	invJwtBearer := "Bearer " + jwt
	p := models.AnimalRequest{Series: "123"}
	requestByte, _ := json.Marshal(p)
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest("POST", "/updateAnimal", requestReader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", invJwtBearer)
	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(510).JSON(fiber.Map{
				"error": "wrong call",
			})
		},
		SigningKey: []byte(salt),
	}))
	app.Post("/updateAnimal", UpdateAnimal(nil))
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v", fiber.StatusUnauthorized, res.Status)
	}

}

func TestUpdateAnimalWithNoDataRequest(t *testing.T) {
	salt := "soli!"
	jwt := getTestToken("qqqqqqwwwwwwr", salt)
	invJwtBearer := "Bearer " + jwt
	req, err := http.NewRequest("POST", "/updateAnimal", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", invJwtBearer)
	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(510).JSON(fiber.Map{
				"error": "wrong call",
			})
		},
		SigningKey: []byte(salt),
	}))
	app.Post("/updateAnimal", UpdateAnimal(nil))
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v", fiber.StatusUnauthorized, res.Status)
	}

}

func TestUpdateAnimal(t *testing.T) {
	s, c := configAnimals()
	id, _ := primitive.ObjectIDFromHex("1234")
	p := []models.AnimalRequest{{Status: "asdasd", BirthDate: "1.08.2019"},
		{BirthDate: "2010-10-10", Species: "cattle", UtilityType: "meat", Sex: "male", Status: "sold", MotherSeries: "1234", Breed: "MM"},
		{Series: "123456789", Species: "cattle", UtilityType: "meat", Sex: "male", Status: "sold", MotherSeries: "1234", Breed: "MM"},
		{Series: "123456789", BirthDate: "2006-11-12", UtilityType: "meat", Sex: "male", Status: "sold", MotherSeries: "1234", Breed: "MM"},
		{Series: "123456789", BirthDate: "2010-10-10", Species: "cattle", Sex: "male", Status: "sold", MotherSeries: "1234", Breed: "MM"},
		{Series: "123456789", BirthDate: "2010-10-10", Species: "cattle", UtilityType: "meat", Status: "sold", MotherSeries: "1234", Breed: "MM"},
		{Series: "123456789", BirthDate: "2010-10-10", Species: "cattle", UtilityType: "meat", Sex: "male", MotherSeries: "1234", Breed: "MM"},
		{Series: "123456789", BirthDate: "2010-10-10", Species: "cattle", UtilityType: "meat", Sex: "male", Status: "sold", Breed: "MM"},
		{Series: "123456789", BirthDate: "2010-10-10", Species: "cattle", UtilityType: "meat", Sex: "male", Status: "sold", MotherSeries: "1234"},
		{Series: "123456789", BirthDate: "2010-10-10", Species: "cattle", UtilityType: "meat", Sex: "male", Status: "sold", MotherSeries: "1234", Breed: "MM"}}
	salt := "sol"
	token := getTestToken(id.Hex(), salt)
	tokenBearer := "Bearer " + token

	type Exp struct {
		status int
		Error  string
	}

	expRes := []Exp{{400, "wrong date format"}, {400, "animal must be specified"},
		{400, "wrong date format"}, {400, "animal must be specified"},
		{400, "animal must be specified"}, {400, "animal must be specified"},
		{400, "animal must be specified"}, {400, "animal must be specified"},
		{400, "animal must be specified"}, {400, "animal not exists"}}

	for i := 0; i < 10; i++ {
		requestByte, _ := json.Marshal(p[i])
		requestReader := bytes.NewReader(requestByte)
		req, err := http.NewRequest("POST", "/updateAnimal", requestReader)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tokenBearer)
		app := fiber.New()
		app.Use(jwtware.New(jwtware.Config{
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				return ctx.Status(510).JSON(fiber.Map{
					"error": "wrong call",
				})
			},
			SigningKey: []byte(salt),
		}))
		app.Post("/updateAnimal", UpdateAnimal(s))
		res, err := app.Test(req)
		if res.StatusCode == expRes[i].status {
			type resultt struct {
				Errors string `json:"error"`
			}
			result := new(resultt)
			if err := json.NewDecoder(res.Body).Decode(result); err != nil {
				log.Fatal(err)
			}
			if result.Errors != expRes[i].Error {
				t.Errorf("bad %v response, bug error result: %v  exp: %v", i, result.Errors, expRes[i].Error)
			}
		} else {
			t.Errorf("bad %v response, code", i)
		}
	}
	timeMongo,err:=time.Parse("2006-01-02", p[9].BirthDate)
	if err!=nil{
		t.Fatalf("can not create models.Animal, error: %v",err)
	}
	an:=models.Animal{OwnerId:id,Series: p[9].Series,BirthDate: timeMongo,Species: p[9].Species,UtilityType: p[9].UtilityType,Sex:p[9].Sex,Status: p[9].Status,MotherSeries: p[9].MotherSeries,Breed: p[9].Breed}
	c.InsertOne(context.Background(), an)

	requestByte, _ := json.Marshal(p[9])
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest("POST", "/updateAnimal", requestReader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenBearer)
	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(510).JSON(fiber.Map{
				"error": "wrong call",
			})
		},
		SigningKey: []byte(salt),
	}))
	app.Post("/updateAnimal", UpdateAnimal(s))
	res, err := app.Test(req)
	if res.StatusCode == 200 {
	} else {
		t.Errorf("bad response result: %v exp: 200",res.StatusCode)
	}
	c.DeleteOne(req.Context(),an)

}
