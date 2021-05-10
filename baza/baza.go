package baza

/*
packege incjalizujący połączenie z bazą danych
*/

import (
	"context"
	"time"
	"os"
	log "github.com/sirupsen/logrus"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB mongo.Database

func connectToMongo() { //funkcja łączy się z bazą
	fmt.Println("elo: "+os.Getenv("MONGO_URL"))
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
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
	DatabseName:=os.Getenv("MONGO_DB")
	if DatabseName == ""{
		DB = *client.Database("praktyki")
	}else {
		DB = *client.Database(DatabseName)
	}
	

}
func init () {
	connectToMongo()
}
