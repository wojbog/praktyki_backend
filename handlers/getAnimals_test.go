package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/animals"
	animalsService "github.com/wojbog/praktyki_backend/service/animals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getTestToken(id string, salt string) string {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: time.Now().Add(10 * time.Second).Unix()})

	token, err := claims.SignedString([]byte(salt))
	if err != nil {
		log.Fatalf("Cannot create test JWT token: %s", err)
	}
	return token

}

func TestGetAnimalsWithInvalidJwtReturnsUnathorized(t *testing.T) {
	salt := "soli!"
	jwt := getTestToken("qqqqqqwwwwwwr", salt)
	invJwtBearer := "Bearer " + jwt
	req, err := http.NewRequest("GET", "/getAnimals", nil)
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
	app.Get("/getAnimals", GetAnimals(nil))
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Bad Status!\nExpected:%v\nReceived:%v", fiber.StatusUnauthorized, res.Status)
	}

}
func TestGetAnimals(t *testing.T) {
	s, c := configAnimals()

	salt := "sol"
	testUserId, _ := primitive.ObjectIDFromHex("0")
	token := getTestToken(testUserId.Hex(), salt)
	tokenBearer := "Bearer " + token

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
	c.InsertMany(context.Background(), []interface{}{exampleAnimal, exampleAnimal2})
	defer c.DeleteMany(context.Background(), bson.M{"ownerId": testUserId})

	testQueryParams := []map[string]string{
		{"status": "2"},
		{"breed": "1"},
		{"minBirthDate": "2015-01-12"},
		{"minBirthDate": "1999-01-01"},
		{"minBirthDate": "1999-01-01", "maxBirthDate": "2007-07-10"},
		{"maxBirthDate": "2031-07-11"},
		{"mminBirthDate": "dawno temu"},
		{"utilityType": "1"},
		{"species": "nie wiem w sumie"},
	}

	exp := []struct {
		status  int
		animals []models.Animal
	}{
		{fiber.StatusOK, []models.Animal{exampleAnimal2}},
		{fiber.StatusOK, []models.Animal{exampleAnimal}},
		{fiber.StatusOK, []models.Animal{exampleAnimal}},
		{fiber.StatusOK, []models.Animal{exampleAnimal, exampleAnimal2}},
		{fiber.StatusOK, nil},
		{fiber.StatusOK, []models.Animal{exampleAnimal, exampleAnimal2}},
		{fiber.StatusOK, []models.Animal{exampleAnimal, exampleAnimal2}},
		{fiber.StatusOK, []models.Animal{exampleAnimal}},
		{fiber.StatusOK, nil},
	}

	if len(testQueryParams) != len(exp) {
		t.Fatalf("Numbers of test cases and expectations are not equal!")
	}

	for i := 0; i < len(testQueryParams); i++ {
		req, err := http.NewRequest("GET", "/getAnimals", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", tokenBearer)
		q := req.URL.Query()

		for key, value := range testQueryParams[i] {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()

		app := fiber.New()
		app.Use(jwtware.New(jwtware.Config{
			SigningKey: []byte(salt),
		}))
		app.Get("/getAnimals", GetAnimals(s))

		res, err := app.Test(req)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != exp[i].status {
			t.Errorf("Bad status code!\n expected: %v \n received: %+v", exp[i].status, res.StatusCode)
		} else {
			r := &struct {
				Animals []models.Animal `json:"animals"`
				Error   string          `json:"error"`
			}{}

			if err := json.NewDecoder(res.Body).Decode(r); err != nil {
				t.Fatal("Could not decode json response")
			}

			if !reflect.DeepEqual(r.Animals, exp[i].animals) {
				t.Errorf("Bad animals array!\nQueryParams: %+v\nRequest params: %+v\nexpected: %+v\nreceived: %+v\n\n", testQueryParams[i], req.URL.Query(), exp[i].animals, r.Animals)
			}
		}
	}
}

func configAnimals() (*animalsService.Service, *mongo.Collection) {
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
	}

	DatabseName := datebasename
	db := client.Database(DatabseName)

	col := *db.Collection("animals")

	animalCol := animals.NewCollection(&col)

	s := animalsService.NewService(animalCol)
	return s, &col

}
