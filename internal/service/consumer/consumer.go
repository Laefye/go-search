package consumer

import (
	"context"
	"fmt"

	"github.com/Laefye/go-search/internal/repository"
	"github.com/Laefye/go-search/internal/service/dto"
)

type ConsumerService struct {
	counter repository.QueryStatsRepository
}

func NewConsumerService(
	counter repository.QueryStatsRepository,
) *ConsumerService {
	return &ConsumerService{
		counter: counter,
	}
}

func (s *ConsumerService) Consume(ctx context.Context, event *dto.SearchQueryEvent) error {
	err := s.counter.IncrementQuery(ctx, event.Query, event.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to add query: %w", err)
	}
	return nil
}
