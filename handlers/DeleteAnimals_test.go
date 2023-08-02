package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteAnimalWithInvalidJwtReturnsUnathorized(t *testing.T) {
	salt := "soli!"
	jwt := getTestToken("qqqqqqwwwwwwr", salt)
	invJwtBearer := "Bearer " + jwt
	req, err := http.NewRequest("POST", "/deleteAnimal", nil)
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
	app.Post("/deleteAnimal", DeleteAnimal(nil))
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v", fiber.StatusUnauthorized, res.Status)
	}

}

func TestDeleteAnimal(t *testing.T) {
	s, c := configAnimals()
	id, _ := primitive.ObjectIDFromHex("1234")
	p := []models.AnimalRequest{{Status: "asdasd"}, {Series: "qwertyuiop"}, {Series: "123456789"}}

	salt := "sol"
	token := getTestToken(id.Hex(), salt)
	tokenBearer := "Bearer " + token

	type Exp struct {
		status int
		Error  string
	}

	expRes := []Exp{{400, "animal must be specified by series number"}, {400, "animal not exists"}, {200, ""}}

	for i := 0; i < 2; i++ {
		requestByte, _ := json.Marshal(p[i])
		requestReader := bytes.NewReader(requestByte)
		req, err := http.NewRequest("POST", "/deleteAnimal", requestReader)
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
		app.Post("/deleteAnimal", DeleteAnimal(s))
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
				t.Errorf("bad %v response, bug error", i)
			}
		} else {
			t.Errorf("bad %v response, code", i)
		}
	}

	an := models.Animal{OwnerId: id, Series: "123456789"}
	c.InsertOne(context.Background(), an)
	requestByte, _ := json.Marshal(p[2])
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest("POST", "/deleteAnimal", requestReader)
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
	app.Post("/deleteAnimal", DeleteAnimal(s))
	res, err := app.Test(req)
	if res.StatusCode == expRes[2].status {
	} else {
		t.Errorf("bad %v response, code", 2)
	}

	req, err = http.NewRequest("POST", "/deleteAnimal", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenBearer)
	app = fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(510).JSON(fiber.Map{
				"error": "wrong call",
			})
		},
		SigningKey: []byte(salt),
	}))
	app.Post("/deleteAnimal", DeleteAnimal(s))
	res, err = app.Test(req)
	if res.StatusCode == 400 {
	} else {
		fmt.Println(res.StatusCode)
		t.Errorf("bad %v response, code", 3)
	}

}
