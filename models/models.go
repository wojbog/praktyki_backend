package models

import (
	"errors"
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

func Animal2Request(animal Animal) AnimalRequest {
	const layoutISO = "2006-01-02"
	date := animal.BirthDate.Format(layoutISO)
	req := AnimalRequest{
		Series:       animal.Series,
		BirthDate:    date,
		Species:      animal.Species,
		UtilityType:  animal.UtilityType,
		Sex:          animal.Sex,
		Status:       animal.Status,
		MotherSeries: animal.MotherSeries,
		Breed:        animal.Breed,
	}
	return req
}

func Request2Animal(req AnimalRequest, ownerId string) (Animal, error) {

	var animal Animal
	////converting ownerId
	if primitiveOwnerId, err := primitive.ObjectIDFromHex(ownerId); err != nil && ownerId != "0" {
		return Animal{}, errors.New("cannot convert ownerId")
	} else {
		animal.OwnerId = primitiveOwnerId
	}

	////converting date
	const layoutISO = "2006-01-02"
	if date, err := time.Parse(layoutISO, req.BirthDate); err != nil {
		return Animal{}, errors.New("cannot parse date")

	} else {
		animal.Series = req.Series
		animal.BirthDate = date
		animal.Species = req.Species
		animal.UtilityType = req.UtilityType
		animal.Sex = req.Sex
		animal.Status = req.Status
		animal.MotherSeries = req.MotherSeries
		animal.Breed = req.Breed
	}

	return animal, nil
}

type Animal struct {
	OwnerId      primitive.ObjectID `bson:"ownerId" json:"ownerId"`
	Series       string             `bson:"series" json:"series" validate:"required,alphanum"`
	BirthDate    time.Time          `bson:"birthDate" json:"birthDate" validate:"required"`
	Species      string             `bson:"species" json:"species" validate:"required,oneof=cattle pig"`
	UtilityType  string             `bson:"utilityType" json:"utilityType" validate:"required,oneof=combined milk meat"`
	Sex          string             `bson:"sex" json:"sex" validate:"required,oneof=male female"`
	Status       string             `bson:"status" json:"status" validate:"required,oneof=sold current carrion"`
	MotherSeries string             `bson:"motherSeries" json:"motherSeries" validate:"required,alphanum"`
	Breed        string             `bson:"breed" json:"breed" validate:"required,oneof=MM SM HO RW PBZ WBP PULL ZB"`
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
	MaxBirthDate string             `json:"maxBirthDate" mapstructure:",omitempty"`
	Species      string             `bson:"species,omitempty" json:"species" mapstructure:"species,omitempty"`
	Breed        string             `bson:"breed,omitempty" json:"breed" mapstructure:"breed,omitempty"`
	Status       string             `bson:"status,omitempty" json:"status" mapstructure:"status,omitempty"`
	Sex          string             `bson:"sex,omitempty" json:"sex" mapstructure:"sex,omitempty"`
	UtilityType  string             `bson:"utilityType,omitempty" json:"utilityType" mapstructure:"utilityType,omitempty"`
	MotherSeries string             `bson:"motherSeries,omitempty" json:"motherId" mapstructure:"motherSeries,omitempty"`
}
