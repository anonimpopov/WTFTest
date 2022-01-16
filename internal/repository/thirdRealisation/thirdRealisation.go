package thirdRealisation

import (
	"context"
	"encoding/json"
	"github.com/anonimpopov/WTFTest/internal/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type ActionObject struct {
	Total     int                       `bson:"total" json:"total"`
	Countries map[string]map[string]int `bson:"countries" json:"countries"`
}

type DatabaseObject struct {
	Total     int                      `bson:"total" json:"total"`
	Actions   map[string]*ActionObject `bson:"actions" json:"actions"`
	Timestamp primitive.DateTime       `bson:"timestamp" json:"timestamp"`
}

type cache struct {
	mx sync.Mutex
	DatabaseObject
}

type Repository struct {
	db    *mongo.Collection
	cache *cache
}

func New(db *mongo.Collection) *Repository {
	ob := &cache{DatabaseObject: DatabaseObject{Total: 0, Actions: make(map[string]*ActionObject)}}
	return &Repository{db, ob}
}

func (r *Repository) Init() chan<- bool {
	stop := make(chan bool)
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-stop:
				return
			case _ = <-ticker.C:
				r.cache.mx.Lock()

				if r.cache.Total != 0 {
					r.cache.Timestamp = primitive.NewDateTimeFromTime(time.Now())
					_, err := r.db.InsertOne(context.TODO(), r.cache.DatabaseObject)
					if err != nil {
						logrus.Errorf("error during saving batch: %v", err)
					}
					r.cache = &cache{DatabaseObject: DatabaseObject{Total: 0, Actions: make(map[string]*ActionObject)}}
				}

				r.cache.mx.Unlock()
			}

		}
	}()

	return stop
}

func (r *Repository) SaveAction(action models.Action) {
	r.cache.mx.Lock()
	defer r.cache.mx.Unlock()

	r.cache.Total++
	if _, ok := r.cache.Actions[action.Type]; !ok {
		r.cache.Actions[action.Type] = &ActionObject{Total: 1, Countries: make(map[string]map[string]int)}
		r.cache.Actions[action.Type].Countries[action.Country] = make(map[string]int)
		r.cache.Actions[action.Type].Countries[action.Country]["total"] = 1
		return
	}

	r.cache.Actions[action.Type].Total++
	if _, ok := r.cache.Actions[action.Type].Countries[action.Country]; !ok {
		r.cache.Actions[action.Type].Countries[action.Country] = make(map[string]int)
		r.cache.Actions[action.Type].Countries[action.Country]["total"] = 1
		return
	}

	r.cache.Actions[action.Type].Countries[action.Country]["total"]++
}

func (r *Repository) GetMetrics(ctx context.Context, from int64, to int64) ([]byte, error) {
	fromTime := primitive.NewDateTimeFromTime(time.Unix(from, 0))
	toTime := primitive.NewDateTimeFromTime(time.Unix(to, 0))

	filter := bson.D{
		{"timestamp", bson.D{
			{"$gt", fromTime},
			{"$lt", toTime},
		}},
	}

	cursor, err := r.db.Find(ctx, filter)
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			logrus.Errorf("mongo cursor close error: %v", err)
		}
	}(cursor, ctx)

	if err != nil {
		return nil, err
	}

	var metrics []DatabaseObject
	for cursor.Next(ctx) {
		var element DatabaseObject
		if err := cursor.Decode(&element); err != nil {
			return nil, err
		}

		metrics = append(metrics, element)
	}
	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}
	return jsonMetrics, nil
}
