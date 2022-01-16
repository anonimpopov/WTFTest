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

type actionObject struct {
	Total     int                       `bson:"total" json:"total"`
	Countries map[string]map[string]int `bson:"countries" json:"countries"`
}

type databaseObject struct {
	Total     int                      `bson:"total" json:"total"`
	Actions   map[string]*actionObject `bson:"actions" json:"actions"`
	Timestamp primitive.DateTime       `bson:"timestamp" json:"timestamp"`
}

type cache struct {
	mx sync.Mutex
	databaseObject
}

type repository struct {
	db    *mongo.Collection
	cache *cache
}

func New(db *mongo.Collection) *repository {
	ob := &cache{databaseObject: databaseObject{Total: 0, Actions: make(map[string]*actionObject)}}
	return &repository{db, ob}
}

func (r *repository) Init() chan<- bool {
	stop := make(chan bool)
	ticker := time.NewTicker(time.Hour)

	go func() {
		for {
			select {
			case <-stop:
				return
			case _ = <-ticker.C:
				r.cache.mx.Lock()

				if r.cache.Total != 0 {
					r.cache.Timestamp = primitive.NewDateTimeFromTime(time.Now())
					_, err := r.db.InsertOne(context.TODO(), r.cache.databaseObject)
					if err != nil {
						logrus.Errorf("error during saving batch: %v", err)
					}
					r.cache = &cache{databaseObject: databaseObject{Total: 0, Actions: make(map[string]*actionObject)}}
				}

				r.cache.mx.Unlock()
			}

		}
	}()

	return stop
}

func (r *repository) SaveAction(action models.Action) {
	r.cache.mx.Lock()
	defer r.cache.mx.Unlock()

	r.cache.Total++
	if _, ok := r.cache.Actions[action.Type]; !ok {
		r.cache.Actions[action.Type] = &actionObject{Total: 1, Countries: make(map[string]map[string]int)}
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

func (r *repository) GetMetrics(ctx context.Context, from int64, to int64) ([]byte, error) {
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

	var metrics []databaseObject
	for cursor.Next(ctx) {
		var element databaseObject
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
