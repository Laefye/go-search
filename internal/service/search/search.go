package search

import (
	"context"
	"errors"
	"strings"

	"github.com/Laefye/go-search/internal/rabbitmq"
	"github.com/Laefye/go-search/internal/service/dto"
)

var ErrInvalidQuery = errors.New("invalid query")

type SearchService struct {
	publisher *rabbitmq.Publisher
}

func NewSearchService(publisher *rabbitmq.Publisher) *SearchService {
	return &SearchService{publisher: publisher}
}

func (s *SearchService) Publish(ctx context.Context, query dto.SearchQueryEvent) error {
	trimmed := strings.TrimSpace(query.Query)
	if trimmed == "" {
		return ErrInvalidQuery
	}

	return s.publisher.PublishQuery(ctx, query)
}
