package service

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wojbog/praktyki_backend/models"
	"github.com/wojbog/praktyki_backend/repository/animals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetAnimals(t *testing.T) {
	s, c := config()

	testUserId, _ := primitive.ObjectIDFromHex("0")

	exampleAnimal := models.Animal{
		OwnerId:      testUserId,
		Series:       "abcdef",
		BirthDate:    time.Date(2020, 10, 3, 0, 0, 0, 0, time.UTC),
		Species:      "1",
		UtilityType:  "1",
		Sex:          "1",
		Status:       "1",
		MotherSeries: "1",
		Breed:        "1",
	}

	exampleAnimal2 := models.Animal{
		OwnerId:      testUserId,
		Series:       "123456",
		BirthDate:    time.Date(2010, 12, 3, 0, 0, 0, 0, time.UTC),
		Species:      "2",
		UtilityType:  "2",
		Sex:          "2",
		Status:       "2",
		MotherSeries: "2",
		Breed:        "2",
	}

	exampleAnimal3 := models.Animal{
		OwnerId:      testUserId,
		Series:       "zyxwut",
		BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
		Species:      "3",
		UtilityType:  "3",
		Sex:          "3",
		Status:       "3",
		MotherSeries: "3",
		Breed:        "3",
	}

	c.InsertMany(context.Background(), []interface{}{exampleAnimal, exampleAnimal2, exampleAnimal3})
	defer c.DeleteMany(context.Background(), bson.M{"ownerId": testUserId})

	testCases := []models.AnimalFilters{
		{Series: "123456"},
		{
			OwnerId:      testUserId,
			MinBirthDate: "2006-10-10",
			MaxBirthDate: "2030-10-10",
		},
		{
			OwnerId:      testUserId,
			MinBirthDate: "sie",
			MaxBirthDate: "ma co tam",
		},
		{
			OwnerId:      testUserId,
			MinBirthDate: "ugh hackers everywhere",
			MaxBirthDate: "2005-12-19",
		},
		{
			OwnerId:      testUserId,
			MinBirthDate: "2007-01-01",
			MaxBirthDate: "im like an analphabetic",
		},
		{OwnerId: testUserId, MotherSeries: "3"},
		{OwnerId: testUserId, Sex: "2"},
		{OwnerId: testUserId, Breed: "chleb"},
		{OwnerId: testUserId, UtilityType: "coco jambo"},
	}

	exp := []struct {
		animals []models.Animal
		err     error
	}{
		{[]models.Animal{exampleAnimal2}, nil},
		{[]models.Animal{exampleAnimal, exampleAnimal2}, nil},
		{[]models.Animal{exampleAnimal, exampleAnimal2, exampleAnimal3}, nil},
		{[]models.Animal{exampleAnimal3}, nil},
		{[]models.Animal{exampleAnimal, exampleAnimal2}, nil},
		{[]models.Animal{exampleAnimal3}, nil},
		{[]models.Animal{exampleAnimal2}, nil},
		{nil, nil},
		{nil, nil},
	}

	if len(testCases) != len(exp) {
		t.Fatal("Numbers of test cases and expectations are not equal!")
	}

	for i := 0; i < len(testCases); i++ {
		if res, err := s.GetAnimals(context.Background(), testCases[i]); err != nil {
			t.Errorf("Unespected error occured: %s", err)
		} else {
			if !reflect.DeepEqual(res, exp[i].animals) {
				t.Errorf("Wrong query result!\nFilter: %+v\nExpected: %v\nReceived: %v", testCases[i], exp[i].animals, res)
			}
			if err != exp[i].err {
				t.Errorf("Wrong error response!\nFilter: %+v\nExpected: %v\nReceived: %v", testCases[i], exp[i].err, err)
			}
		}
	}
}

func TestDeleteAnimal(t *testing.T) {
	s, c := config()
	if err := s.DeleteAnimal(context.Background(), models.AnimalFilters{}); err != animals.AnimalNotexist {
		t.Error("no AnimalNotexist")
	}
	id, _ := primitive.ObjectIDFromHex("1234")
	p := models.Animal{OwnerId: id, Series: "147258369"}

	c.InsertOne(context.Background(), p)
	if err := s.DeleteAnimal(context.Background(), models.AnimalFilters{OwnerId: id, Series: "147258369"}); err != nil {
		t.Error("can not delete")
	}
}

func TestUpdateAnimal(t *testing.T) {
	s, c := config()
	if err := s.UpdateAnimal(context.Background(), models.Animal{}); err != animals.AnimalNotexist {
		t.Error("no AnimalNotexist")
	}
	id, _ := primitive.ObjectIDFromHex("1234")
	p := models.Animal{OwnerId: id, Series: "147258369"}

	c.InsertOne(context.Background(), p)
	if err := s.UpdateAnimal(context.Background(), models.Animal{OwnerId: id, Series: "147258369"}); err != nil {
		t.Error("can not update")
	}
	c.DeleteOne(context.Background(), bson.M{"ownerId": p.OwnerId, "series": p.Series})
}

func TestAddNewAnimalWithInvalidInputReturnValidationError(t *testing.T) {
	s, c := config()
	ctx := context.Background()

	testUserId, _ := primitive.ObjectIDFromHex("0")
	testCases := []models.Animal{
		{
			OwnerId:      testUserId,
			Series:       "1",
			BirthDate:    time.Date(2019, 10, 3, 0, 0, 0, 0, time.UTC),
			Species:      "pig",
			UtilityType:  "",
			Sex:          "female",
			Status:       "sold",
			MotherSeries: "1",
			Breed:        "MM",
		}, {
			OwnerId:      testUserId,
			Series:       "2",
			BirthDate:    time.Date(2010, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "",
			UtilityType:  "meat",
			Sex:          "male",
			Status:       "ezzz",
			MotherSeries: "2",
			Breed:        "SM",
		}, {
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "cattle",
			UtilityType:  "commbbbined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RaaaaaaaW",
		},
		{
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "cattle",
			UtilityType:  "commdddbbbined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RW",
		}, {
			OwnerId:      testUserId,
			Series:       "",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "cattle",
			UtilityType:  "commbbbined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RW",
		}, {
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "",
			UtilityType:  "combined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RW",
		},
		{
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "cattle",
			UtilityType:  "",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RW",
		}, {
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "cattle",
			UtilityType:  "combined",
			Sex:          "",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RW",
		}, {
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "cattle",
			UtilityType:  "combined",
			Sex:          "female",
			Status:       "",
			MotherSeries: "3",
			Breed:        "RW",
		},
		{
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "pig",
			UtilityType:  "commbbbined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "",
			Breed:        "RW",
		},
		{
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "pig",
			UtilityType:  "commbbbined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "",
		}, {
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "cattle",
			UtilityType:  "commbbbined",
			Sex:          "haha",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RW",
		},
		{
			OwnerId:      testUserId,
			Series:       "3@#",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "piig",
			UtilityType:  "commbbbined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RW",
		}, {
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "",
			UtilityType:  "commbbbined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3!@##!",
			Breed:        "RW",
		},
	}

	defer c.DeleteMany(ctx, bson.M{"ownerId": testUserId})
	for i := 0; i < len(testCases); i++ {
		_, err := s.AddNewAnimal(ctx, testCases[i])
		exp := "validation error"

		if err.Error() != exp {
			t.Errorf("Wrong error!\nExpected: %+v\nReceived: %+v\n", exp, err)
		}
	}
}
func TestAddNewAnimalWithValidInputReturnCreatedAnimal(t *testing.T) {
	s, c := config()
	ctx := context.Background()

	testUserId, _ := primitive.ObjectIDFromHex("0")
	testCases := []models.Animal{
		{
			OwnerId:      testUserId,
			Series:       "1",
			BirthDate:    time.Date(2019, 10, 3, 0, 0, 0, 0, time.UTC),
			Species:      "pig",
			UtilityType:  "milk",
			Sex:          "female",
			Status:       "sold",
			MotherSeries: "1",
			Breed:        "MM",
		}, {
			OwnerId:      testUserId,
			Series:       "2",
			BirthDate:    time.Date(2010, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "cattle",
			UtilityType:  "meat",
			Sex:          "male",
			Status:       "current",
			MotherSeries: "2",
			Breed:        "SM",
		}, {
			OwnerId:      testUserId,
			Series:       "3",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "cattle",
			UtilityType:  "combined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RW",
		},
	}

	exp := []models.AnimalRequest{
		{
			Series:       "1",
			BirthDate:    "2019-10-03",
			Species:      "pig",
			UtilityType:  "milk",
			Sex:          "female",
			Status:       "sold",
			MotherSeries: "1",
			Breed:        "MM",
		}, {
			Series:       "2",
			BirthDate:    "2010-12-03",
			Species:      "cattle",
			UtilityType:  "meat",
			Sex:          "male",
			Status:       "current",
			MotherSeries: "2",
			Breed:        "SM",
		}, {
			Series:       "3",
			BirthDate:    "2000-12-03",
			Species:      "cattle",
			UtilityType:  "combined",
			Sex:          "female",
			Status:       "carrion",
			MotherSeries: "3",
			Breed:        "RW",
		},
	}
	if len(testCases) != len(exp) {
		t.Fatal("Numbers of test cases and expectations are not equal!")
	}

	defer c.DeleteMany(ctx, bson.M{"ownerId": testUserId})
	for i := 0; i < len(testCases); i++ {
		res, err := s.AddNewAnimal(ctx, testCases[i])

		if err != nil {
			t.Errorf("Wrong error!\nExpected: <nil>\nReceived: %+v\n", err)
		}
		if res != exp[i] {
			t.Errorf("Wrong result!\nExpected: %+v\nReceived: %+v\n\n", exp[i], res)
		}
	}
}

func config() (*Service, *mongo.Collection) {
	str1 := os.Getenv("MONGO_URL")
	if str1 == "" {
		log.Fatal("NO MONGO URL")
	}
	str2 := os.Getenv("MONGO_DB")
	if str2 == "" {
		log.Fatal("NO MONGO DB")
	}
	mongo_URL := str1
	datebasename := str2

	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_URL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("connect to DB")
	}
	DatabseName := datebasename
	db := client.Database(DatabseName)

	acol := *db.Collection("animals")
	animalCol := animals.NewCollection(&acol)

	s := NewService(animalCol)

	return s, &acol
}
