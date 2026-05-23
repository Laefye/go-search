package consumer

import (
	"context"
	"fmt"

	"github.com/Laefye/go-search/internal/repository"
	"github.com/Laefye/go-search/internal/service/dto"
)

type ConsumerService struct {
	counter *repository.CounterRepository
}

func NewConsumerService(
	counter *repository.CounterRepository,
) *ConsumerService {
	return &ConsumerService{
		counter: counter,
	}
}

func (s *ConsumerService) Consume(ctx context.Context, event *dto.SearchQueryEvent) error {
	err := s.counter.AddQuery(ctx, event.Query, event.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to add query: %w", err)
	}
	return nil
}
