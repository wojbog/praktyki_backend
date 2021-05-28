package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/wojbog/praktyki_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"
// 	"reflect"
// 	"testing"
// 	"time"

// 	"github.com/form3tech-oss/jwt-go"
// 	"github.com/gofiber/fiber/v2"
// 	jwtware "github.com/gofiber/jwt/v2"
// 	"github.com/wojbog/praktyki_backend/models"
// 	"github.com/wojbog/praktyki_backend/repository/animals"
// 	animalsService "github.com/wojbog/praktyki_backend/service/animals"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

func TestCreateAnimalsWithInvalidJwtReturnsUnathorized(t *testing.T) {
	salt := "soli!"
	jwt := getTestToken("qqqqqqwwwwwwr", salt)
	invJwtBearer := "Bearer " + jwt

	requestByte, _ := json.Marshal(nil)
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest("POST", "/createAnimal", requestReader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", invJwtBearer)
	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			//weird status to distinguish handler error from middleware error
			return ctx.Status(510).JSON(fiber.Map{
				"error": "wrong call",
			})
		},
		SigningKey: []byte(salt),
	}))
	app.Post("/createAnimal", CreateAnimal(nil))
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	//t.Fatal(res.Body)
	if res.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v", fiber.StatusUnauthorized, res.Status)
	}
}

func TestCreateAnimalsWithInvalidBodyreturnBadRequest(t *testing.T) {
	salt := "soli!"
	testUserId, _ := primitive.ObjectIDFromHex("0")
	jwt := getTestToken(testUserId.Hex(), salt)
	jwtBearer := "Bearer " + jwt

	req, err := http.NewRequest("POST", "/createAnimal", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", jwtBearer)
	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(salt),
	}))
	app.Post("/createAnimal", CreateAnimal(nil))
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	//t.Fatal(res.Body)
	if res.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v", fiber.StatusBadRequest, res.Status)
	}
}

func TestCreateAnimalsWithWrongDateReturnsBadRequest(t *testing.T) {
	salt := "soli!"
	testUserId, _ := primitive.ObjectIDFromHex("0")
	jwt := getTestToken(testUserId.Hex(), salt)
	jwtBearer := "Bearer " + jwt

	testCase := models.AnimalRequest{
		Series:       "3",
		BirthDate:    "dropDatabase()",
		Species:      "cattle",
		UtilityType:  "commdddbbbined",
		Sex:          "female",
		Status:       "carrion",
		MotherSeries: "3",
		Breed:        "RW",
	}

	requestByte, _ := json.Marshal(testCase)
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest("POST", "/createAnimal", requestReader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", jwtBearer)

	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(salt),
	}))
	app.Post("/createAnimal", CreateAnimal(nil))
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	//t.Fatal(res.Body)
	if res.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v", fiber.StatusBadRequest, res.Status)
	}

	type result struct {
		Success bool `json:"success"`
		Error   string
	}

	exp := result{false, "cannot parse date"}
	var r result
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if r != exp {
		t.Errorf("Wrong response body!\n Expected: %+v\nReceived: %+v", exp, r)
	}

}

func TestCreateAnimalsWithValidationErrorReturnBadRequest(t *testing.T) {
	s, _ := configAnimals()
	salt := "soli!"
	testUserId, _ := primitive.ObjectIDFromHex("0")
	jwt := getTestToken(testUserId.Hex(), salt)
	jwtBearer := "Bearer " + jwt

	testCase := models.AnimalRequest{
		Series:       "3",
		BirthDate:    "2010-10-10",
		Species:      "catetle",
		UtilityType:  "",
		Sex:          "femalelelele",
		Status:       "carrion",
		MotherSeries: "3",
		Breed:        "RW",
	}

	requestByte, _ := json.Marshal(testCase)
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest("POST", "/createAnimal", requestReader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", jwtBearer)

	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(salt),
	}))
	app.Post("/createAnimal", CreateAnimal(s))
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	//t.Fatal(res.Body)
	if res.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v", fiber.StatusBadRequest, res.Status)
	}

	type errorBody struct {
		InvalidFields []string `json:"invalidFields"`
	}

	type result struct {
		Success   bool      `json:"success"`
		Error     string    `json:"error"`
		ErrorBody errorBody `json:"errorBody"`
	}

	exp := result{
		Success: false,
		Error:   "validation error",
		ErrorBody: errorBody{
			[]string{"Species", "UtilityType", "Sex"},
		},
	}
	var r result
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(r, exp) {
		t.Errorf("Wrong response body!\n Expected: %+v\nReceived: %+v", exp, r)
	}
}

func TestCreateAnimalsWithRedundantAnimalReturnsAnimalExistsError(t *testing.T) {
	s, c := configAnimals()
	ctx := context.Background()

	testUserId, _ := primitive.ObjectIDFromHex("0")
	salt := "soli!"
	jwt := getTestToken(testUserId.Hex(), salt)
	jwtBearer := "Bearer " + jwt

	test := models.Animal{OwnerId: testUserId,
		Series:       "abcdef",
		BirthDate:    time.Date(2039, 10, 3, 0, 0, 0, 0, time.UTC),
		Species:      "2",
		UtilityType:  "qwer",
		Sex:          "qwer",
		Status:       "qwe",
		MotherSeries: "qwer",
		Breed:        "qwer",
	}
	_, err := c.InsertOne(ctx, test)
	defer c.DeleteMany(ctx, bson.M{"ownerId": testUserId})
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	testCase := models.AnimalRequest{
		Series:       "abcdef",
		BirthDate:    "2010-10-10",
		Species:      "cattle",
		UtilityType:  "meat",
		Sex:          "female",
		Status:       "carrion",
		MotherSeries: "3",
		Breed:        "RW",
	}

	requestByte, _ := json.Marshal(testCase)
	requestReader := bytes.NewReader(requestByte)

	req, err := http.NewRequest("POST", "/createAnimal", requestReader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", jwtBearer)

	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(salt),
	}))
	app.Post("/createAnimal", CreateAnimal(s))
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v", fiber.StatusBadRequest, res.Status)
	}

	type result struct {
		Success bool   `json:"success"`
		Error   string `jnos:"error"`
	}

	exp := result{false, "animal exists"}
	var r result
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if r != exp {
		t.Errorf("Wrong response body!\n Expected: %+v\nReceived: %+v", exp, r)
	}
}

func TestCreateAnimalWithValidInputsReturnSucces(t *testing.T) {
	s, c := configAnimals()
	ctx := context.Background()

	salt := "soli!"
	testUserId, _ := primitive.ObjectIDFromHex("0")
	jwt := getTestToken(testUserId.Hex(), salt)
	jwtBearer := "Bearer " + jwt

	testCase := models.AnimalRequest{
		Series:       "3",
		BirthDate:    "2010-10-10",
		Species:      "pig",
		UtilityType:  "meat",
		Sex:          "male",
		Status:       "sold",
		MotherSeries: "3",
		Breed:        "RW",
	}

	requestByte, _ := json.Marshal(testCase)
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest("POST", "/createAnimal", requestReader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", jwtBearer)

	app := fiber.New()
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(salt),
	}))
	app.Post("/createAnimal", CreateAnimal(s))
	res, err := app.Test(req)
	defer c.DeleteMany(ctx, bson.M{"ownerId": testUserId})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusOK {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v\nRepsonseBody: %v", fiber.StatusOK, res.StatusCode, string(b))
	}

	type result struct {
		Success bool                 `json:"success"`
		Animal  models.AnimalRequest `json:"animal"`
	}

	exp := result{
		Success: true,
		Animal: models.AnimalRequest{
			Series:       "3",
			BirthDate:    "2010-10-10",
			Species:      "pig",
			UtilityType:  "meat",
			Sex:          "male",
			Status:       "sold",
			MotherSeries: "3",
			Breed:        "RW"},
	}

	var r result
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(r, exp) {
		t.Errorf("Bad Response!\nExpected: %+v\nReceived: %+v", exp, r)
	}

}
