package metricBatch

import (
	"context"
	"github.com/anonimpopov/WTFTest/internal/models"
)

type SaverAndReader interface {
	Saver
	Reader
}
type Saver interface {
	SaveAction(models.Action) error
}

type Reader interface {
	GetMetrics(context.Context, int64, int64) ([]byte, error)
}

type Service struct {
	repo SaverAndReader
}

func New(repo SaverAndReader) *Service {
	return &Service{repo}
}

func (s *Service) SaveMetric(action models.Action) error {
	return s.repo.SaveAction(action)
}

func (s *Service) GetMetrics(ctx context.Context, from int64, to int64) ([]byte, error) {
	return s.repo.GetMetrics(ctx, from, to)
}
