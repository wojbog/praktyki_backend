package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator"
	"github.com/mitchellh/mapstructure"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/animals"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	animalsCol *animals.Collection
}

func NewService(animalsCol *animals.Collection) *Service {
	return &Service{animalsCol}
}

///////Errors
var NoDataAnimalError = errors.New("Animal must be specified")

type ValidationError struct {
	InvalidFields []string `json:"invalidFields"`
}

func (e *ValidationError) Error() string {
	return "validation error"
}

//GetAnimals returns array of filter.ownerId's animals.
//It converts models.AnimalFilters to compatible with db map with filters
//If no animals return null
func (s *Service) GetAnimals(ctx context.Context, filter models.AnimalFilters) ([]models.Animal, error) {
	const layoutISO = "2006-01-02"

	maxDate, err := time.Parse(layoutISO, filter.MaxBirthDate)
	if err != nil {
		maxDate, _ = time.Parse(layoutISO, "9999-12-31")
	}
	minDate, _ := time.Parse(layoutISO, filter.MinBirthDate)

	filter.MinBirthDate = ""
	filter.MaxBirthDate = ""

	mapFilter := make(map[string]interface{})

	if err := mapstructure.Decode(filter, &mapFilter); err != nil {
		return nil, err
	}

	mapFilter["birthDate"] = bson.M{
		"$gte": minDate,
		"$lt":  maxDate,
	}

	animals, err := s.animalsCol.GetAnimals(ctx, mapFilter)
	return animals, err
}

//AddNewAnimal
//Validate animal and call insertAnimal
func (s *Service) AddNewAnimal(ctx context.Context, animal models.Animal) (models.AnimalRequest, error) {
	validate := validator.New()
	if err := validate.Struct(animal); err != nil {
		var TabErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			TabErrors = append(TabErrors, err.Field())
		}
		return models.AnimalRequest{}, &ValidationError{InvalidFields: TabErrors}
	}

	animalRes, err := s.animalsCol.InsertAnimal(ctx, animal)
	return animalRes, err
}

//DeleteAnimal return error if DeleteAnimal in repository return error
func (s *Service) DeleteAnimal(ctx context.Context, filter models.AnimalFilters) error {

	if err := s.animalsCol.DeleteAnimal(ctx, filter); err != nil {
		return err
	}
	return nil
}

//UpdateAnimal return error if UpdateAnimal in repository return error
func (s *Service) UpdateAnimal(ctx context.Context, animal models.Animal) error {

	if err := s.animalsCol.UpdateAnimal(ctx, animal); err != nil {
		return err
	}
	return nil
}