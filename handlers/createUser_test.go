package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/user"
	userService "github.com/wojbog/praktyki_backend/service/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//func TestCreateUser test createUser Handler
func TestCreateUser(t *testing.T) {
	s, c := Config()
	p := []models.NewUser{{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6q", Email: "sss4tefan@elo.pl"},
		{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6q", Email: "sss4tefan@elo.pl"},
		{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "0009000", Pass: "Wojtek6q", Email: "@sss4tefan@elo.pl"}}
	exp := []struct {
		status int
		err    string
	}{
		{200, ""},
		{400, "user exists"}}
	for i := 0; i < 2; i++ {
		requestByte, _ := json.Marshal(p[i])
		requestReader := bytes.NewReader(requestByte)
		req := httptest.NewRequest("POST", "http://127.0.0.1:3001/createUser", requestReader)
		req.Header.Set("Content-Type", "application/json")
		app := fiber.New()
		app.Post("/createUser", CreateUser(s))
		res, _ := app.Test(req, 3000)

		if res.StatusCode == exp[i].status {
			type resultt struct {
				Success bool   `json:"success"`
				Errors  string `json:"error"`
			}
			result := new(resultt)
			if err := json.NewDecoder(res.Body).Decode(result); err != nil {
				log.Fatal(err)
			}
			if result.Errors != exp[i].err {
				t.Errorf("bad response")
			}
		} else {
			t.Errorf("bad code")
		}

	}
	requestByte, _ := json.Marshal(p[2])
	requestReader := bytes.NewReader(requestByte)
	req := httptest.NewRequest("POST", "http://127.0.0.1:3001/createUser", requestReader)
	req.Header.Set("Content-Type", "application/json")
	app := fiber.New()
	app.Post("/createUser", CreateUser(s))
	res, _ := app.Test(req, 3000)

	if res.StatusCode == 400 {
		type taberr struct {
			Tab []string `json:"invalidFields"`
		}
		type resultt struct {
			Success bool   `json:"success"`
			Errors  taberr `json:"error"`
		}
		result := new(resultt)
		if err := json.NewDecoder(res.Body).Decode(result); err != nil {
			log.Fatal(err)
		}
		if result.Errors.Tab[0] != "Email" {
			t.Errorf("bad response")
		}
	} else {
		t.Errorf("bad code")
	}

	req = httptest.NewRequest("POST", "http://127.0.0.1:3001/createUser", nil)
	req.Header.Set("Content-Type", "application/json")
	app = fiber.New()
	app.Post("/createUser", CreateUser(s))
	res, _ = app.Test(req, 3000)

	if res.StatusCode == 400 {
		log.Info("PASS")
	} else {
		t.Errorf("bad code")
	}
	c.DeleteOne(context.Background(), bson.M{"email": "sss4tefan@elo.pl"})

}

//Config configuration function
func Config() (*userService.Service, *mongo.Collection) {
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

	col := *db.Collection("users")

	userCol := user.NewCollection(&col)

	s := userService.NewService(userCol)
	return s, &col

}
