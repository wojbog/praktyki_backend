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
	CanNotDeleteError = errors.New("can not delete")
	AnimalNotexist    = errors.New("animal not exists")
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

//DeleteAnimal delete animal in DB
//return error if can not delete
func (colAnim *Collection) DeleteAnimal(ctx context.Context, filter models.AnimalFilters) error {

	if result, err := colAnim.col.DeleteOne(ctx, bson.M{"ownerId": filter.OwnerId, "series": filter.Series}); err != nil {
		return CanNotDeleteError
	} else if result.DeletedCount == 0 {
		return AnimalNotexist
	} else {
		log.Info("deleted animal ownerid:" + filter.OwnerId.Hex() + " series: " + filter.Series)
		return nil
	}
}

//InsertAnimal
//return inserted Animal
//if Animalexists return AnimalExistsError
func (colAnim *Collection) InsertAnimal(ctx context.Context, animal models.Animal) (models.AnimalRequest, error) {

	if err := colAnim.col.FindOne(ctx, bson.M{"series": animal.Series}); err != nil {
		//if err isnt no document error
		if err.Err() != mongo.ErrNoDocuments {
			return models.AnimalRequest{}, AnimalExistsError
		}
	}

	res, err := colAnim.col.InsertOne(ctx, animal)
	if err != nil {
		return models.AnimalRequest{}, InternalServerError
	}

	a := new(models.Animal)
	if err := colAnim.col.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(a); err != nil {
		log.Info(err)
		return models.AnimalRequest{}, InternalServerError
	}

	return models.Animal2Request(*a), nil

}

func (colAnim *Collection) UpdateAnimal(ctx context.Context, animal models.Animal) error {

	if result,err:=colAnim.col.ReplaceOne(ctx,bson.M{"ownerId":animal.OwnerId,"series":animal.Series},animal);err!=nil{
		return InternalServerError
	} else if result.MatchedCount==0 {
		return AnimalNotexist
	} else {
		log.Info("updated animal ownerid:" + animal.OwnerId.Hex() + " series: " + animal.Series)
		return nil
	}
	

}
