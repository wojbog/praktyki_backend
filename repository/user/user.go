package user

import (
	"context"
	"errors"
	"time"
	"os"
	
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/repository/person"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Baza struct {
	db mongo.Database
	mongo_URL string 
	datebasename string 
}

func InsertUser(ctx context.Context,user *person.Person ) (string,error) {//dodawanie użytkownika
	var baza Baza
	baza.conf(ctx)
	per:=&person.Person{}
	if errv:=baza.db.Collection("users").FindOne(ctx, bson.M{"email":user.Email}).Decode(per); errv ==nil {
		return "",errors.New("account exists")
	} else {	
	result ,err:=baza.db.Collection("users").InsertOne(ctx,user)
	if err !=nil{
		log.Fatal(err)
		return "",err
	} 
	id:=result.InsertedID
		
	return  id.(primitive.ObjectID).Hex(),nil
}}

func (baza *Baza) conf(ctx context.Context) {
	
	str1:=os.Getenv("MONGO_URL")
	if str1 ==""{
		log.Fatal("NO MONGO URL")
	}
	str2:=os.Getenv("MONGO_DB")
	if str2 ==""{
		log.Fatal("NO MONGO DB")
	}
	baza.mongo_URL=str1
	baza.datebasename=str2

	client, err := mongo.NewClient(options.Client().ApplyURI(baza.mongo_URL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("connect to DB")
	}
	DatabseName:=baza.datebasename
	if DatabseName == ""{
		baza.db = *client.Database("praktyki")
	}else {
		baza.db = *client.Database(DatabseName)
	}
	
}