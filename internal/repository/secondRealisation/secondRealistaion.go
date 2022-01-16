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

type ActionObject struct {
	Total     int `bson:"total"`
	Countries map[string]map[string]int
}

type DatabaseObject struct {
	Total   int                     `bson:"total"`
	Actions map[string]ActionObject `bson:"actions"`
}

type Repository struct {
	db     *mongo.Collection
	itemId primitive.ObjectID
}

func New(db *mongo.Collection) *Repository {
	return &Repository{db, primitive.ObjectID{}}
}

func (r *Repository) InitItem() error {
	res, err := r.db.InsertOne(context.TODO(), DatabaseObject{0, make(map[string]ActionObject)})
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

func (r *Repository) SaveAction(ctx context.Context, action models.Action) error {
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
