package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Laefye/go-search/internal/service/dto"
	"github.com/redis/go-redis/v9"
)

type CounterRepository struct {
	rbd *redis.Client
}

const keyFormat = "search.query.%d"
const globalKey = "search.query.global"

func NewRedisClient(rdb *redis.Client) *CounterRepository {
	return &CounterRepository{rbd: rdb}
}

func formatKey(timestamp time.Time) string {
	return fmt.Sprintf(keyFormat, timestamp.Unix()/60)
}

func (c *CounterRepository) AddQuery(ctx context.Context, query string, timestamp time.Time) error {
	key := formatKey(timestamp)

	err := c.rbd.ZIncrBy(ctx, key, 1, query).Err()
	if err != nil {
		return err
	}

	err = c.rbd.ZIncrBy(ctx, globalKey, 1, query).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *CounterRepository) GetQueries(ctx context.Context, timestamp time.Time) ([]dto.QueryEntry, error) {
	key := formatKey(timestamp)

	results, err := c.rbd.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var queryEntries []dto.QueryEntry
	for _, z := range results {
		queryEntries = append(queryEntries, dto.QueryEntry{
			Query: z.Member.(string),
			Count: int(z.Score),
		})
	}

	return queryEntries, nil
}

func (c *CounterRepository) GetGlobalQueriesTop(ctx context.Context, limit int64) ([]dto.QueryEntry, error) {
	results, err := c.rbd.ZRevRangeWithScores(ctx, globalKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	var entries []dto.QueryEntry
	for _, z := range results {
		entries = append(entries, dto.QueryEntry{
			Query: z.Member.(string),
			Count: int(z.Score),
		})
	}

	return entries, nil
}

func (c *CounterRepository) DecrGlobalQuery(ctx context.Context, query string, count int) error {
	err := c.rbd.ZIncrBy(ctx, globalKey, float64(-count), query).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *CounterRepository) DeleteQueries(ctx context.Context, timestamp time.Time) error {
	key := formatKey(timestamp)
	return c.rbd.Del(ctx, key).Err()
}
