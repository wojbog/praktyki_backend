package animals

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	AnimalExistsError   = errors.New("animal exists")
	InternalServerError = errors.New("internal server error")
)

type Collection struct {
	col *mongo.Collection
}

func NewCollection(col *mongo.Collection) *Collection {
	return &Collection{col}
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

func (colAnim *Collection) InsertAnimal(ctx context.Context, animal models.Animal) (models.AnimalRequest, error) {
	// if err = colAnim.FindOne(ctx)

	if err := colAnim.col.FindOne(ctx, bson.M{"series": animal.Series}); err != nil {
		//if err isnt no document error
		if err.Err() != mongo.ErrNoDocuments {
			return models.AnimalRequest{}, AnimalExistsError
		}
	}

	res, err := colAnim.col.InsertOne(ctx, animal)
	if err != nil {
		log.Info("err insss")
		return models.AnimalRequest{}, InternalServerError
	}

	a := new(models.Animal)
	if err := colAnim.col.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(a); err != nil {
		log.Info("err fioooafter")
		log.Info(err)
		return models.AnimalRequest{}, InternalServerError
	}

	return models.Animal2Request(*a), nil

}
