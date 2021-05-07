package baza

/*
packege incjalizujący połączenie z bazą danych
*/

import (
	"context"
	"time"
	"os"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB mongo.Database

func ConnectToMongo() { //funkcja łączy się z bazą

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("URL_BAZA")))
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

	DB = *client.Database("praktyki")

}
