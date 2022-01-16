package metric

import (
	"context"
	"github.com/anonimpopov/WTFTest/internal/models"
)

type Saver interface {
	SaveAction(context.Context, models.Action) error
}

type Service struct {
	repo1 Saver
	repo2 Saver
}

func New(repo1 Saver, repo2 Saver) *Service {
	return &Service{repo1, repo2}
}

func (s *Service) SaveMetric(ctx context.Context, action models.Action, version int) error {
	switch version {
	case 1:
		return s.repo1.SaveAction(ctx, action)
	case 2:
		return s.repo2.SaveAction(ctx, action)
	}

	return nil
}
