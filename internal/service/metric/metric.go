package metric

import (
	"context"
	"github.com/anonimpopov/WTFTest/internal/models"
)

type SaverAndReader interface {
	Saver
	Reader
}

type Saver interface {
	SaveAction(context.Context, models.Action) error
}

type Reader interface {
}

type Service struct {
	repo SaverAndReader
}

func New(repo SaverAndReader) *Service {
	return &Service{repo}
}

func (s *Service) SaveMetric(ctx context.Context, action models.Action) error {
	return s.repo.SaveAction(ctx, action)
}
