//user communication with datebase
package user

import (
	"context"
	"errors"
	
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	
)


//Collection type store collecion instance in Datebase
type Collection struct {
	col *mongo.Collection
}




//InsertUser add new user to Datebase
//return id of new user if correct added
func (colUser *Collection)InsertUser(ctx context.Context,user *service.Person ) (string,error) {//dodawanie użytkownika
	
	per:=&service.Person{}
	if errv:=colUser.col.FindOne(ctx, bson.M{"email":user.Email}).Decode(per); errv ==nil {
		return "",errors.New("user exists")
	} else {	
	result ,err:=colUser.col.InsertOne(ctx,user)
	if err !=nil{
		log.Fatal(err)
		return "",err
	} 
	id:=result.InsertedID
		
	return  id.(primitive.ObjectID).Hex(),nil
}}


//NewCollection creates new instance of Collection
func NewCollection(colection *mongo.Collection) *Collection {
	return &Collection{colection}
}