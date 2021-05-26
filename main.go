package main

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/repository/animals"
	"github.com/wojbog/praktyki_backend/repository/user"
	"github.com/wojbog/praktyki_backend/server"
	"github.com/wojbog/praktyki_backend/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	app := fiber.New()

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

	userCol := user.NewCollection(db.Collection("users"))

	animalsCol := animals.NewCollection(db.Collection("animals"))

	s := service.NewService(userCol, animalsCol)

	server.Routing(app, s)

	PORT := os.Getenv("PORT")
	if PORT != "" {
		app.Listen(":" + PORT)
	} else {
		log.Panic("NO PORT")
	}

}
