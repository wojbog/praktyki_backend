package serviceDB

/*
package do zarządzania bazą danych
*/
import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/baza"
	"github.com/wojbog/praktyki_backend/person"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
)


func InsertUser(ctx context.Context,user *person.Person ) (string,error) {//dodawanie użytkownika
	per:=&person.Person{}
	if errv:=baza.DB.Collection("users").FindOne(ctx, bson.M{"email":user.Email}).Decode(per); errv ==nil {
		log.Info("account exists")
		return "",errors.New("account exists")
	} else {	
	result ,err:=baza.DB.Collection("users").InsertOne(ctx,user)
	if err !=nil{
		log.Fatal(err)
		return "",err
	} 
	id:=result.InsertedID
		
	

	
	return  id.(primitive.ObjectID).Hex(),nil
}}