package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/Laefye/go-search/internal/rabbitmq/events"
	"github.com/Laefye/go-search/internal/repository"
)

type ConsumerService struct {
	counter repository.QueryStatsRepository
	guard   repository.GuardRepository
}

func NewConsumerService(
	counter repository.QueryStatsRepository,
	guard repository.GuardRepository,
) *ConsumerService {
	return &ConsumerService{
		counter: counter,
		guard:   guard,
	}
}

func (s *ConsumerService) Consume(ctx context.Context, event *events.SearchQueryEvent) error {
	err := s.guard.IncrementRequest(ctx, event.UserID)
	if err != nil {
		return fmt.Errorf("failed to add query: %w", err)
	}

	shouldSkip, err := s.guard.ShouldSkip(ctx, event.UserID)
	if err != nil {
		return fmt.Errorf("failed to check if query should be skipped: %w", err)
	}

	if shouldSkip {
		log.Printf("Skipping query from user %s due to rate limiting\n", event.UserID)
		return nil
	}

	err = s.counter.IncrementQuery(ctx, event.Query, event.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to increment query count: %w", err)
	}

	return nil
}
