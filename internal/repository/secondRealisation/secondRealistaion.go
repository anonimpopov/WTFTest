package secondRealisation

import (
	"context"
	"errors"
	"fmt"
	"github.com/anonimpopov/WTFTest/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type actionObject struct {
	Total     int `bson:"total" json:"total"`
	Countries map[string]map[string]int
}

type databaseObject struct {
	Total   int                      `bson:"total" json:"total"`
	Actions map[string]*actionObject `bson:"actions" json:"actions"`
}

type repository struct {
	db     *mongo.Collection
	itemId primitive.ObjectID
}

func New(db *mongo.Collection) *repository {
	return &repository{db, primitive.ObjectID{}}
}

func (r *repository) Init() error {
	res, err := r.db.InsertOne(context.TODO(), databaseObject{0, make(map[string]*actionObject)})
	if err != nil {
		return err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("cant cast interface")
	}

	r.itemId = id

	return nil
}

func (r *repository) SaveAction(ctx context.Context, action models.Action) error {
	update := bson.D{
		{"$inc", bson.D{
			{"total", 1},
			{fmt.Sprintf("actions.%v.total", action.Type), 1},
			{fmt.Sprintf("actions.%v.countries.%v.total", action.Type, action.Country), 1},
		}},
	}

	res, err := r.db.UpdateByID(ctx, r.itemId, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("cant modify (count 0)")
	}

	return nil
}
