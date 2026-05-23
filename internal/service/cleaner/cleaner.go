package cleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/Laefye/go-search/internal/repository"
)

type CleanerService struct {
	repo    *repository.CounterRepository
	minutes int
}

func NewCleanerService(repo *repository.CounterRepository, minutes int) *CleanerService {
	return &CleanerService{repo: repo, minutes: minutes}
}

func (s *CleanerService) Clean(ctx context.Context, timestamp time.Time) error {
	timestamp = timestamp.Add(-time.Duration(s.minutes) * time.Minute)

	data, err := s.repo.GetQueries(ctx, timestamp)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	err = s.repo.DeleteQueries(ctx, timestamp)
	if err != nil {
		return fmt.Errorf("failed to delete queries: %w", err)
	}

	for _, entry := range data {
		err = s.repo.DecrGlobalQuery(ctx, entry.Query, entry.Count)
		if err != nil {
			return fmt.Errorf("failed to decrement global query: %w", err)
		}
	}

	return nil
}
