package animals

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/wojbog/praktyki_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetAnimals(t *testing.T) {
	c := config(t)

	testUserId, _ := primitive.ObjectIDFromHex("0")

	exampleAnimal := models.Animal{
		OwnerId:      testUserId,
		Series:       "abcdef",
		BirthDate:    time.Date(2019, 10, 3, 0, 0, 0, 0, time.UTC),
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

	c.col.InsertMany(context.Background(), []interface{}{exampleAnimal, exampleAnimal2, exampleAnimal3})
	defer c.col.DeleteMany(context.Background(), bson.M{"ownerId": testUserId})

	testCases := []map[string]interface{}{
		{"series": "123456"},
		{"birthDate": bson.M{
			"$gte": time.Date(2005, 12, 12, 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC),
		}},
		{"birthDate": bson.M{
			"$gte": time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
		}},
		{"birthDate": bson.M{
			"$gte": time.Date(2005, 12, 12, 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(2007, 12, 31, 0, 0, 0, 0, time.UTC),
		}},
		{"species": "1"},
		{"utilityType": "2"},
		{"sex": "3"},
		{"status": "3"},
		{"motherSeries": "1"},
		{"breed": "3"},
	}

	exp := []struct {
		animals []models.Animal
		err     error
	}{
		{[]models.Animal{exampleAnimal2}, nil},
		{[]models.Animal{exampleAnimal, exampleAnimal2}, nil},
		{[]models.Animal{exampleAnimal}, nil},
		{nil, nil},
		{[]models.Animal{exampleAnimal}, nil},
		{[]models.Animal{exampleAnimal2}, nil},
		{[]models.Animal{exampleAnimal3}, nil},
		{[]models.Animal{exampleAnimal3}, nil},
		{[]models.Animal{exampleAnimal}, nil},
		{[]models.Animal{exampleAnimal3}, nil},
	}

	if len(testCases) != len(exp) {
		t.Fatal("Numbers of test cases and expectations are not equal!")
	}

	for i := 0; i < len(testCases); i++ {
		testCases[i]["ownerId"] = testUserId
		if res, err := c.GetAnimals(context.Background(), testCases[i]); err != nil {
			t.Errorf("Unexpected error occured: %s", err)
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
	c := config(t)
	if err := c.DeleteAnimal(context.Background(), models.AnimalFilters{}); err != AnimalNotexist {
		t.Error("no CanNotDeleteError")
	}
	id, _ := primitive.ObjectIDFromHex("1234")
	p := models.Animal{OwnerId: id, Series: "147258369"}

	c.col.InsertOne(context.Background(), p)
	if err := c.DeleteAnimal(context.Background(), models.AnimalFilters{OwnerId: id, Series: "147258369"}); err != nil {
		t.Error("can not delete")
	}
}
func TestInsertAnimalWithRedundantAnimalReturnAnimalExistsError(t *testing.T) {
	c := config(t)
	ctx := context.Background()

	testUserId, _ := primitive.ObjectIDFromHex("0")
	example := models.Animal{OwnerId: testUserId,
		Series:       "abcdef",
		BirthDate:    time.Date(2019, 10, 3, 0, 0, 0, 0, time.UTC),
		Species:      "1",
		UtilityType:  "1",
		Sex:          "1",
		Status:       "1",
		MotherSeries: "1",
		Breed:        "1"}

	c.col.InsertOne(ctx, example)
	defer c.col.DeleteMany(ctx, bson.M{"ownerId": testUserId})

	_, err := c.InsertAnimal(ctx, example)
	if err != AnimalExistsError {
		t.Errorf("Wrong error!\nExpected: %v\nReceived: %v", AnimalExistsError, err)
	}

	test := models.Animal{OwnerId: testUserId,
		Series:       "abcdef",
		BirthDate:    time.Date(2039, 10, 3, 0, 0, 0, 0, time.UTC),
		Species:      "2",
		UtilityType:  "qwer",
		Sex:          "qwer",
		Status:       "qwe",
		MotherSeries: "qwer",
		Breed:        "qwer"}

	_, err = c.InsertAnimal(ctx, test)
	if err != AnimalExistsError {
		t.Errorf("Wrong error!\nExpected: %v\nReceived: %v", AnimalExistsError, err)
	}

}
func TestInsertAnimalsWithUniqueAnimalInserAnimal(t *testing.T) {
	c := config(t)
	ctx := context.Background()

	testUserId, _ := primitive.ObjectIDFromHex("0")
	testCases := []models.Animal{
		{
			OwnerId:      testUserId,
			Series:       "abcdef",
			BirthDate:    time.Date(2019, 10, 3, 0, 0, 0, 0, time.UTC),
			Species:      "1",
			UtilityType:  "1",
			Sex:          "1",
			Status:       "1",
			MotherSeries: "1",
			Breed:        "1",
		}, {
			OwnerId:      testUserId,
			Series:       "123456",
			BirthDate:    time.Date(2010, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "2",
			UtilityType:  "2",
			Sex:          "2",
			Status:       "2",
			MotherSeries: "2",
			Breed:        "2",
		}, {
			OwnerId:      testUserId,
			Series:       "zyxwut",
			BirthDate:    time.Date(2000, 12, 3, 0, 0, 0, 0, time.UTC),
			Species:      "3",
			UtilityType:  "3",
			Sex:          "3",
			Status:       "3",
			MotherSeries: "3",
			Breed:        "3",
		},
	}

	exp := []models.AnimalRequest{
		{
			Series:       "abcdef",
			BirthDate:    "2019-10-03",
			Species:      "1",
			UtilityType:  "1",
			Sex:          "1",
			Status:       "1",
			MotherSeries: "1",
			Breed:        "1",
		}, {
			Series:       "123456",
			BirthDate:    "2010-12-03",
			Species:      "2",
			UtilityType:  "2",
			Sex:          "2",
			Status:       "2",
			MotherSeries: "2",
			Breed:        "2",
		}, {
			Series:       "zyxwut",
			BirthDate:    "2000-12-03",
			Species:      "3",
			UtilityType:  "3",
			Sex:          "3",
			Status:       "3",
			MotherSeries: "3",
			Breed:        "3",
		},
	}
	if len(testCases) != len(exp) {
		t.Fatal("Numbers of test cases and expectations are not equal!")
	}

	for i := 0; i < len(testCases); i++ {
		res, err := c.InsertAnimal(ctx, testCases[i])
		defer c.col.DeleteOne(ctx, bson.M{"ownerId": testUserId})
		if err != nil {
			t.Errorf("Wrong error!\nAnimal:%v\nExpected: nil\nReceived: %v\n\n", testCases[i], err)
		}

		if res != exp[i] {
			t.Errorf("Wrong animal result!\nExpected: %+v\nReceived: %+v\n\n", exp[i], res)
		}
	}

}

func config(t *testing.T) *Collection {
	str1 := os.Getenv("MONGO_URL")
	if str1 == "" {
		t.Fatal("NO MONGO URL")
	}
	str2 := os.Getenv("MONGO_DB")
	if str2 == "" {
		t.Fatal("NO MONGO DB")
	}
	mongo_URL := str1
	datebasename := str2

	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_URL))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}

	DatabseName := datebasename
	db := client.Database(DatabseName)

	col := *db.Collection("animals")

	userCol := NewCollection(&col)

	return userCol

}
