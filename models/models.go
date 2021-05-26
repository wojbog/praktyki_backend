package models

import (
	"time"

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

// type AnimalWithId struct {
// 	Id           primitive.ObjectID `bson:"_id" json:"id,omitempty"`
// 	OwnerId      primitive.ObjectID `json:"ownerId"`
// 	Series       string             `json:"series" `
// 	BirthDate    time.Time          `json:"birthDate" `
// 	Species      string             `json:"species"`
// 	UtilityType  string             `json:"utilityType"`
// 	Sex          string             `json:"sex"`
// 	Status       string             `json:"status"`
// 	MotherSeries string             `json:"motherSeries"`
// 	Breed        string             `json:"breed"`
// }

type Animal struct {
	OwnerId      primitive.ObjectID `bson:"ownerId" json:"ownerId"`
	Series       string             `bson:"series" json:"series" validate:"required,aplhanum"`
	BirthDate    time.Time          `bson:"birthDate" json:"birthDate" validate:"required,numeric"`
	Species      string             `bson:"species" json:"species" validate:"required,alpha"`
	UtilityType  string             `bson:"utilityType" json:"utilityType" validate:"required,utility_type"`
	Sex          string             `bson:"sex" json:"sex" validate:"required,gender"`
	Status       string             `bson:"status" json:"status" validate:"required"`
	MotherSeries string             `bson:"motherSeries" json:"motherSeries" validate:"required"`
	Breed        string             `bson:"breed" json:"breed" validate:"required,alpha"`
}

type AnimalRequest struct {
	Series       string `json:"series"`
	BirthDate    string `json:"birthDate"`
	Species      string `json:"species"`
	UtilityType  string `json:"utilityType"`
	Sex          string `json:"sex"`
	Status       string `json:"status"`
	MotherSeries string `json:"motherSeries"`
	Breed        string `json:"breed"`
}

type AnimalFilters struct {
	OwnerId      primitive.ObjectID `bson:"ownerId" json:"ownerId" mapstructure:"ownerId,omitempty"`
	Series       string             `bson:"series,omitempty" json:"series" mapstructure:"series,omitempty"`
	MinBirthDate string             `json:"minBirthDate" mapstructure:",omitempty"`
	MaxBirthDate string             `json:":" mapstructure:",omitempty"`
	Species      string             `bson:"species,omitempty" json:"species" mapstructure:"species,omitempty"`
	Breed        string             `bson:"breed,omitempty" json:"breed" mapstructure:"breed,omitempty"`
	Status       string             `bson:"status,omitempty" json:"status" mapstructure:"status,omitempty"`
	Sex          string             `bson:"sex,omitempty" json:"sex" mapstructure:"sex,omitempty"`
	UtilityType  string             `bson:"utilityType,omitempty" json:"utilityType" mapstructure:"utilityType,omitempty"`
	MotherSeries string             `bson:"motherSeries,omitempty" json:"motherId" mapstructure:"motherSeries,omitempty"`
}
