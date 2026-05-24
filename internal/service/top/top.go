package top

import (
	"context"
	"fmt"

	"github.com/Laefye/go-search/internal/repository"
	"github.com/Laefye/go-search/internal/service/dto"
)

type TopService struct {
	repo repository.QueryStatsRepository
}

func NewTopService(repo repository.QueryStatsRepository) *TopService {
	return &TopService{repo: repo}
}

func mapToDTO(entries []repository.SearchStats) []dto.QueryEntry {
	dtoEntries := make([]dto.QueryEntry, len(entries))
	for i, entry := range entries {
		dtoEntries[i] = dto.QueryEntry{
			Query: entry.Query,
			Count: entry.Count,
		}
	}
	return dtoEntries
}

func (s *TopService) GetTopQueries(ctx context.Context, limit int64) (*dto.TopQueriesResponse, error) {
	entries, err := s.repo.GetTopQueries(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top queries: %w", err)
	}

	return &dto.TopQueriesResponse{Top: mapToDTO(entries)}, nil
}
