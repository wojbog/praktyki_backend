package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//NewUser type of new user
type NewUser struct {
	Name      string `json:"name" validate:"required,alpha"`
	Surname   string `json:"surname" validate:"required,alpha"`
	Email     string `json:"email" validate:"required,email"`
	Street    string `json:"street" validate:"required,alpha"`
	Number    string `json:"number" validate:"required,alphanum"`
	City      string `json:"city" validate:"required,alpha"`
	Post_code string `json:"post_code" validate:"required,postCode"`
	Pass      string `json:"pass" validate:"required,password"`
}

//UserResponse inform client about user
type UserResponse struct {
	Name      string `json:"name" `
	Surname   string `json:"surname" `
	Email     string `json:"email" `
	Street    string `json:"street" `
	Number    string `json:"number" `
	City      string `json:"city"`
	Post_code string `json:"post_code" `
}

//User whole user struct
type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Email     string             `json:"email"`
	Pass      string             `json:"pass"`
	Name      string             `json:"name" `
	Surname   string             `json:"surname" `
	Street    string             `json:"street" `
	Number    string             `json:"number" `
	City      string             `json:"city"`
	Post_code string             `json:"post_code" `
}
