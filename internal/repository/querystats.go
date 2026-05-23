package repository

import (
	"context"
	"time"
)

type SearchStats struct {
	Query string
	Count int
}

type QueryStatsRepository interface {
	IncrementQuery(ctx context.Context, query string, timestamp time.Time) error
	GetTopQueries(ctx context.Context, limit int64) ([]SearchStats, error)
	Delete(ctx context.Context, to time.Time) error
}
