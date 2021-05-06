package serviceDB
/*
package do zarządzania bazą danych
*/
import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/wojbog/praktyki_backend/baza"
	"github.com/wojbog/praktyki_backend/person"
)


func InsertUser(user *person.Person ) error {//dodawanie użytkownika
	if errv:=baza.DB.Collection("zadania").FindOne(context.Background(), bson.M{"email":user.Email}); errv !=nil {
		return errors.New("konto istnieje")
	} else {	
	res,err:=baza.DB.Collection("users").InsertOne(context.Background(),user)
	if err !=nil{
		fmt.Println(err)
		return err
	}
	fmt.Println(res.InsertedID)
	

	defer baza.DB.Client().Disconnect(baza.Ctx)
	return nil
}
}