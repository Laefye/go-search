package cleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/Laefye/go-search/internal/repository"
)

type CleanerService struct {
	repo    *repository.QueryStatsRepository
	minutes int
}

func NewCleanerService(repo *repository.QueryStatsRepository, minutes int) *CleanerService {
	return &CleanerService{repo: repo, minutes: minutes}
}

func (s *CleanerService) Clean(ctx context.Context, now time.Time) error {
	cutoff := minuteBucket(now).Add(-time.Duration(s.minutes) * time.Minute)

	lastDeleted, err := s.repo.GetLastDeletedMinute(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last deleted minute: %w", err)
	}
	if lastDeleted.IsZero() {
		if err := s.repo.SetLastDeletedMinute(ctx, cutoff); err != nil {
			return fmt.Errorf("failed to set last deleted minute: %w", err)
		}
		return nil
	}

	for t := lastDeleted.Add(time.Minute); !t.After(cutoff); t = t.Add(time.Minute) {
		data, err := s.repo.GetQueries(ctx, t)
		if err != nil {
			return fmt.Errorf("failed to get queries: %w", err)
		}

		if len(data) > 0 {
			err = s.repo.DeleteQueries(ctx, t)
			if err != nil {
				return fmt.Errorf("failed to delete queries: %w", err)
			}

			for _, entry := range data {
				err = s.repo.DecrGlobalQuery(ctx, entry.Query, entry.Count)
				if err != nil {
					return fmt.Errorf("failed to decrement global query: %w", err)
				}
			}
		}

		if err := s.repo.SetLastDeletedMinute(ctx, t); err != nil {
			return fmt.Errorf("failed to set last deleted minute: %w", err)
		}
	}

	return nil
}

func minuteBucket(timestamp time.Time) time.Time {
	return time.Unix((timestamp.Unix()/60)*60, 0)
}
