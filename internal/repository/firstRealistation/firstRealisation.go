package firstRealistation

import (
	"github.com/anonimpopov/WTFTest/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db *mongo.Collection
}

func New(db *mongo.Collection) *Repository {
	return &Repository{db}
}

func (r *Repository) SaveAction(metric models.Action) error {
	return nil
}
