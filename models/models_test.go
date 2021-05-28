package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAnimal2Request(t *testing.T) {
	oi, _ := primitive.ObjectIDFromHex("0")
	a := Animal{
		OwnerId:      oi,
		Series:       "1",
		BirthDate:    time.Date(2000, 10, 10, 0, 0, 0, 0, time.UTC),
		Species:      "1",
		UtilityType:  "1",
		Sex:          "1",
		Status:       "1",
		MotherSeries: "1",
		Breed:        "1",
	}
	exp := AnimalRequest{
		Series:       "1",
		BirthDate:    "2000-10-10",
		Species:      "1",
		UtilityType:  "1",
		Sex:          "1",
		Status:       "1",
		MotherSeries: "1",
		Breed:        "1",
	}

	res := Animal2Request(a)

	if res != exp {
		t.Errorf("Wrong result!\nExpected: %+v\nReceived: %+v", res, exp)
	}
}

func TestRequest2AnimalWithValidInputConvertsProperly(t *testing.T) {
	oi, _ := primitive.ObjectIDFromHex("0")

	req := AnimalRequest{
		Series:       "1",
		BirthDate:    "2000-10-10",
		Species:      "1",
		UtilityType:  "1",
		Sex:          "1",
		Status:       "1",
		MotherSeries: "1",
		Breed:        "1",
	}
	exp := Animal{
		OwnerId:      oi,
		Series:       "1",
		BirthDate:    time.Date(2000, 10, 10, 0, 0, 0, 0, time.UTC),
		Species:      "1",
		UtilityType:  "1",
		Sex:          "1",
		Status:       "1",
		MotherSeries: "1",
		Breed:        "1",
	}

	res, err := Request2Animal(req, "0")
	if err != nil {
		t.Fatalf("Wrong error!\nExpected: <nil>\nReceived: %v", err)
	}

	if res != exp {
		t.Errorf("Wrong result!\nExpected: %+v\nReceived: %+v", res, exp)
	}
}

func TestRequest2AnimalWithInvalidIdReturnsProperError(t *testing.T) {

	req := AnimalRequest{
		Series:       "1",
		BirthDate:    "2000-10-10",
		Species:      "1",
		UtilityType:  "1",
		Sex:          "1",
		Status:       "1",
		MotherSeries: "1",
		Breed:        "1",
	}
	exp := "cannot convert ownerId"

	_, err := Request2Animal(req, "grrrrr")
	if err != nil && err.Error() != exp {
		t.Fatalf("Wrong error!\nExpected: %v\nReceived: %v", err, exp)
	}
}

func TestRequest2AnimalWithInvalidDateReturnsProperError(t *testing.T) {

	req := AnimalRequest{
		Series:       "1",
		BirthDate:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Species:      "1",
		UtilityType:  "1",
		Sex:          "1",
		Status:       "1",
		MotherSeries: "1",
		Breed:        "1",
	}
	exp := "cannot parse date"

	_, err := Request2Animal(req, primitive.NewObjectID().Hex())
	if err != nil && err.Error() != exp {
		t.Fatalf("Wrong error!\nExpected: %v\nReceived: %v", err, exp)
	}
}
