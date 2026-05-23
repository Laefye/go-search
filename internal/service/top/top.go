package top

import (
	"context"
	"fmt"

	"github.com/Laefye/go-search/internal/repository"
	"github.com/Laefye/go-search/internal/service/dto"
)

type TopService struct {
	repo *repository.CounterRepository
}

func NewTopService(repo *repository.CounterRepository) *TopService {
	return &TopService{repo: repo}
}

func (s *TopService) GetTopQueries(ctx context.Context, limit int64) (*dto.TopQueriesResponse, error) {
	entries, err := s.repo.GetGlobalQueriesTop(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top queries: %w", err)
	}

	return &dto.TopQueriesResponse{Top: entries}, nil
}
