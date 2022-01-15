package metric

import "github.com/anonimpopov/WTFTest/internal/models"

type SaverAndReader interface {
	Saver
	Reader
}

type Saver interface {
	SaveAction(metric models.Action) error
}

type Reader interface {
}

type Service struct {
	repo SaverAndReader
}

func New(repo SaverAndReader) *Service {
	return &Service{repo}
}
