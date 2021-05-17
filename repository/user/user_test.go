package user

import (
	"context"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TestInsertUser test InsertUser
func TestInsertUser(t *testing.T) {
	p := models.NewUser{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6q", Email: "sss4tefan@elo.pl"}
	c := config()
	us, _ := c.InsertUser(context.Background(), p)
	if us.Email != "sss4tefan@elo.pl" {
		t.Error()
	}
	_, err := c.InsertUser(context.Background(), p)
	if err.Error() != "user exists" {
		t.Error()
	}

	c.col.DeleteOne(context.Background(), bson.M{"email": "sss4tefan@elo.pl"})
}

//TestGetUserByEmail test GetUserByEmail
func TestGetUserByEmail(t *testing.T) {
	p := models.User{Name: "adsasd", Surname: "dasdasd", City: "asdasd", Number: "23432e", Street: "asdasd", Post_code: "00-000", Pass: "Wojtek6q", Email: "sss4tefan@elo.pl"}
	c := config()

	_, err := c.GetUserByEmail(context.Background(), p)
	if err.Error() != "incorrect data" {
		t.Error()
	}
	c.col.InsertOne(context.Background(), p)
	_, errv := c.GetUserByEmail(context.Background(), p)
	if errv != nil {
		t.Error()
	}
	c.col.DeleteOne(context.Background(), bson.M{"email": "sss4tefan@elo.pl"})
}

//config configuration function
func config() *Collection {
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

	userCol := NewCollection(&col)

	return userCol

}
