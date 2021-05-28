package service

import (
	"context"
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
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

var NoDataAnimalError = errors.New("Animal must be specified")

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

	if animals, err := s.animalsCol.GetAnimals(ctx, mapFilter); err != nil {
		log.Info(err.Error())
		return nil, err
	} else {
		return animals, nil
	}
}

//DeleteAnimal return error if DeleteAnimal in repository return error
func (s *Service) DeleteAnimal(ctx context.Context, filter models.AnimalFilters) error {

	if err := s.animalsCol.DeleteAnimal(ctx, filter); err != nil {
		return err
	}
	return nil
}
