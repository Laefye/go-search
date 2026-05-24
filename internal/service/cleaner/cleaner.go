package cleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/Laefye/go-search/internal/repository"
)

type CleanerService struct {
	repo repository.QueryStatsRepository
}

func NewCleanerService(repo repository.QueryStatsRepository) *CleanerService {
	return &CleanerService{repo: repo}
}

const windowMinutes = 5

func (s *CleanerService) Clean(ctx context.Context, now time.Time) error {
	to := repository.NormalizeMinute(now).Add(-time.Duration(windowMinutes) * time.Minute)

	err := s.repo.Delete(ctx, to)
	if err != nil {
		return fmt.Errorf("failed to delete old queries: %w", err)
	}

	return nil
}
