package firstRealisation

import (
	"context"
	"errors"
	"fmt"
	"github.com/anonimpopov/WTFTest/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type databaseObject struct {
	Total     int                       `bson:"total" json:"total"`
	Actions   map[string]map[string]int `bson:"actions" json:"actions"`
	Countries map[string]map[string]int `bson:"countries" json:"countries"`
}

type repository struct {
	db     *mongo.Collection
	itemId primitive.ObjectID
}

func New(db *mongo.Collection) *repository {
	return &repository{db, primitive.ObjectID{}}
}

func (r *repository) Init() error {
	res, err := r.db.InsertOne(context.TODO(), databaseObject{0, make(map[string]map[string]int), make(map[string]map[string]int)})
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
			{fmt.Sprintf("countries.%v.total", action.Country), 1},
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
