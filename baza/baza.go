package baza
/*
packege incjalizujący połączenie z bazą danych
*/



import (
	"context"
	"fmt"
	"log"
	"time"

	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB mongo.Database 
var Ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
func connectToMongo () {//funkcja łączy się z bazą

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}
	
	err = client.Connect(Ctx)
	if err != nil {
		log.Fatal(err)
	}else {
		fmt.Println("connected...")
	}
	DB=*client.Database("praktyki")
	
}
func init () {
	connectToMongo()
}