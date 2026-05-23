package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Laefye/go-search/internal/service/dto"
	"github.com/redis/go-redis/v9"
)

type QueryStatsRepository struct {
	rdb *redis.Client
}

const keyFormat = "search.query.%d"
const globalKey = "search.query.global"
const lastDeletedMinuteKey = "search.last_deleted_minute"

func NewQueryStatsRepository(rdb *redis.Client) *QueryStatsRepository {
	return &QueryStatsRepository{rdb: rdb}
}

func minuteUnix(timestamp time.Time) int64 {
	return timestamp.Unix() / 60
}

func normalizeMinute(timestamp time.Time) time.Time {
	return time.Unix(minuteUnix(timestamp)*60, 0)
}

func formatKey(timestamp time.Time) string {
	return fmt.Sprintf(keyFormat, minuteUnix(timestamp))
}

func (c *QueryStatsRepository) AddQuery(ctx context.Context, query string, timestamp time.Time) error {
	key := formatKey(timestamp)

	err := c.rdb.ZIncrBy(ctx, key, 1, query).Err()
	if err != nil {
		return err
	}

	err = c.rdb.ZIncrBy(ctx, globalKey, 1, query).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *QueryStatsRepository) GetQueries(ctx context.Context, timestamp time.Time) ([]dto.QueryEntry, error) {
	key := formatKey(timestamp)

	results, err := c.rdb.ZRevRangeWithScores(ctx, key, 0, -1).Result()
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

func (c *QueryStatsRepository) GetGlobalQueriesTop(ctx context.Context, limit int64) ([]dto.QueryEntry, error) {
	results, err := c.rdb.ZRevRangeWithScores(ctx, globalKey, 0, limit-1).Result()
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

func (c *QueryStatsRepository) DecrGlobalQuery(ctx context.Context, query string, count int) error {
	err := c.rdb.ZIncrBy(ctx, globalKey, float64(-count), query).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *QueryStatsRepository) DeleteQueries(ctx context.Context, timestamp time.Time) error {
	key := formatKey(timestamp)
	return c.rdb.Del(ctx, key).Err()
}

func (c *QueryStatsRepository) SetLastDeletedMinute(ctx context.Context, timestamp time.Time) error {
	timestamp = normalizeMinute(timestamp)
	return c.rdb.Set(ctx, lastDeletedMinuteKey, minuteUnix(timestamp), 0).Err()
}

func (c *QueryStatsRepository) GetLastDeletedMinute(ctx context.Context) (time.Time, error) {
	result, err := c.rdb.Get(ctx, lastDeletedMinuteKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	var lastDeletedMinute int64
	_, err = fmt.Sscanf(result, "%d", &lastDeletedMinute)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(lastDeletedMinute*60, 0), nil
}

func (c *QueryStatsRepository) DeleteUntil(ctx context.Context, from time.Time, to time.Time) error {
	from = normalizeMinute(from)
	to = normalizeMinute(to)
	for t := from; !t.After(to); t = t.Add(time.Minute) {
		err := c.DeleteQueries(ctx, t)
		if err != nil {
			return err
		}
	}

	return nil
}
