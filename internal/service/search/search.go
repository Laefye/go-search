package search

import (
	"context"
	"errors"
	"strings"

	"github.com/Laefye/go-search/internal/rabbitmq"
)

var ErrInvalidQuery = errors.New("invalid query")

type SearchService struct {
	publisher *rabbitmq.Publisher
}

func NewSearchService(publisher *rabbitmq.Publisher) *SearchService {
	return &SearchService{publisher: publisher}
}

func (s *SearchService) Publish(ctx context.Context, query string) error {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return ErrInvalidQuery
	}

	return s.publisher.PublishQuery(ctx, trimmed)
}
