package metric

import (
	"context"
	"github.com/anonimpopov/WTFTest/internal/models"
)

type Saver interface {
	SaveAction(context.Context, models.Action) error
}

type Service struct {
	repo Saver
}

func New(repo Saver) *Service {
	return &Service{repo}
}

func (s *Service) SaveMetric(ctx context.Context, action models.Action) error {
	return s.repo.SaveAction(ctx, action)
}
