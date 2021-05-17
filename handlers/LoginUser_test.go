package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

//TestCreateTokenReturnTrueIfFuncReturnToken
//return error if return error
func TestCreateTokenReturnErrorIfFuncReturnError(t *testing.T) {

	if _, err := CreateToken("sdfsdfsd"); err != nil {
		t.Error("incorrect")
	}
}

//TestCreateTokenReturnErrorIfFuncReturnToken
//return error if return Token
func TestCreateTokenReturnErrorIfFuncReturnToken(t *testing.T) {
	os.Setenv("ACCESS_SECRET", "")
	if _, err := CreateToken("fgsfg"); err == nil {
		t.Error("incorrect")
	}
	os.Setenv("ACCESS_SECRET", "secret")
}

//TestLoginUser test LoginUser handler
func TestLoginUser(t *testing.T) {
	p := models.NewUser{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6q", Email: "sss4tefan@elo.pl"}
	s, c := Config()
	s.AddNewUser(context.Background(), p)
	user := []struct {
		Email string `json:"email"`
		Pass  string `json:"pass"`
	}{{"sssssss4tefan@elo.pl", "Wojtek6q"},
		{"sss4tefan@elo.pl", "Dojojtek6q"},
		{"sss4tefan@elo.pl", "Wojtek6q"}}
	exp := []struct {
		succ   bool
		status int
		err    string
	}{
		{false, 400, "incorrect password or email"},
		{false, 400, "incorrect password or email"},
		{true, 200, ""}}
	for i := 0; i < 2; i++ {
		requestByte, _ := json.Marshal(user[i])
		requestReader := bytes.NewReader(requestByte)
		req := httptest.NewRequest("POST", "http://127.0.0.1:3001/login", requestReader)
		req.Header.Set("Content-Type", "application/json")
		app := fiber.New()
		app.Post("/login", LoginUser(s))
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
	requestByte, _ := json.Marshal(user[2])
	requestReader := bytes.NewReader(requestByte)
	req := httptest.NewRequest("POST", "http://127.0.0.1:3001/login", requestReader)
	req.Header.Set("Content-Type", "application/json")
	app := fiber.New()
	app.Post("/login", LoginUser(s))
	res, _ := app.Test(req, 3000)

	if res.StatusCode == exp[2].status {
		type resultt struct {
			Success bool   `json:"success"`
			Errors  string `json:"error"`
		}
		result := new(resultt)
		if err := json.NewDecoder(res.Body).Decode(result); err != nil {
			log.Fatal(err)
		}
		if result.Success != exp[2].succ {
			t.Errorf("bad response")
		}
	} else {
		t.Errorf("bad code")
	}

	c.DeleteOne(context.Background(), bson.M{"email": "sss4tefan@elo.pl"})

	req = httptest.NewRequest("POST", "http://127.0.0.1:3001/login", nil)
	req.Header.Set("Content-Type", "application/json")
	app = fiber.New()
	app.Post("/login", LoginUser(s))
	res, _ = app.Test(req, 3000)
	if res.StatusCode == 400 {
		type resultt struct {
			Success bool   `json:"success"`
			Errors  string `json:"error"`
		}
		result := new(resultt)
		if err := json.NewDecoder(res.Body).Decode(result); err != nil {
			log.Fatal(err)
		}
		if result.Errors != "cannot read request" {
			t.Errorf("bad response")
		}
	} else {
		t.Errorf("bad code")
	}
}
