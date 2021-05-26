package animals

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	col *mongo.Collection
}

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

func NewCollection(col *mongo.Collection) *Collection {
	return &Collection{col}
}
