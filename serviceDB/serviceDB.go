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
	
)


func InsertUser(ctx context.Context,user *person.Person ) error {//dodawanie użytkownika
	per:=&person.Person{}
	if errv:=baza.DB.Collection("users").FindOne(ctx, bson.M{"email":user.Email}).Decode(per); errv ==nil {
		log.Info("account exists")
		return errors.New("account exists")
	} else {	
	_,err:=baza.DB.Collection("users").InsertOne(ctx,user)
	if err !=nil{
		log.Fatal(err)
		return err
	}
	

	
	return nil
}}