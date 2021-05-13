package user

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
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
func (colUser *Collection) InsertUser(ctx context.Context, user models.NewUser) (models.UserResponse, error) { 

	per := models.UserResponse{}
	if errv := colUser.col.FindOne(ctx, bson.M{"email": user.Email}).Decode(&per); errv == nil {
		
		return models.UserResponse{}, errors.New("user exists")
	} else {
		result, err := colUser.col.InsertOne(ctx, user)
		if err != nil {
			log.Fatal(err)
			return models.UserResponse{}, err
		}

		log.Info("success add new user, id: " + result.InsertedID.(primitive.ObjectID).Hex())

		if errv := colUser.col.FindOne(ctx, bson.M{"email": user.Email}).Decode(&per); errv != nil {
			
			return models.UserResponse{}, errors.New("Cannot find in Datebase")
		} else {
			return per, nil
		}
		
	}
}
//GetUserLogin check in datebase
func (colUser *Collection) GetUserLogin(ctx context.Context, user models.UserLogin) (models.UserLogin, error) { 

	per := models.UserLogin{}
	if errv := colUser.col.FindOne(ctx, bson.M{"email": user.Email}).Decode(&per); errv != nil {
		return models.UserLogin{}, errors.New("incorrect data")
	} else {
		return per, nil
	}
}

//NewCollection creates new instance of Collection
func NewCollection(colection *mongo.Collection) *Collection {
	return &Collection{colection}
}
