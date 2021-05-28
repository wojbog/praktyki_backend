package animals

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	col *mongo.Collection
}

var (
	CanNotDeleteError = errors.New("can not delete")
	AnimalNotexist=errors.New("animal not exists")
)

//GetAnimals returns array of model.Animal
//Use filter to filter objects
func (colAnim *Collection) GetAnimals(ctx context.Context, filter interface{}) ([]models.Animal, error) {
	var animals []models.Animal

	if cursor, err := colAnim.col.Find(ctx, filter); err != nil {
		log.Info(err)
		return nil, err
	} else {
		if err = cursor.All(ctx, &animals); err != nil {
			return nil, err
		}
	}
	return animals, nil
}

//DeleteAnimal delete animal in DB
//return error if can not delete
func (colAnim *Collection) DeleteAnimal(ctx context.Context, filter models.AnimalFilters) error {

	if result, err := colAnim.col.DeleteOne(ctx, bson.M{"ownerId": filter.OwnerId, "series": filter.Series}); err != nil {
		return CanNotDeleteError
	} else if result.DeletedCount == 0 {
		return AnimalNotexist
	} else {
		log.Info("delete animal ownerid:" + filter.OwnerId.Hex() + " series: " + filter.Series)
		return nil
	}
}

func NewCollection(col *mongo.Collection) *Collection {
	return &Collection{col}
}
