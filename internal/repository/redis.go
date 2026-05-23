package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueryStatsRepository struct {
	rdb *redis.Client
}

const keyFormat = "search.query.%d"
const globalKey = "search.query.global"
const lastDeletedMinuteKey = "search.last_deleted_minute"

func NewRedisQueryStatsRepository(rdb *redis.Client) *RedisQueryStatsRepository {
	return &RedisQueryStatsRepository{rdb: rdb}
}

func minuteUnix(timestamp time.Time) int64 {
	return timestamp.Unix() / 60
}

func NormalizeMinute(timestamp time.Time) time.Time {
	return time.Unix(minuteUnix(timestamp)*60, 0)
}

func formatKey(timestamp time.Time) string {
	return fmt.Sprintf(keyFormat, minuteUnix(timestamp))
}

type SearchStats struct {
	Query string
	Count int
}

func zsetToEntries(results []redis.Z) []SearchStats {
	entries := make([]SearchStats, 0, len(results))
	for _, z := range results {
		query, ok := z.Member.(string)
		if !ok {
			continue
		}
		entries = append(entries, SearchStats{
			Query: query,
			Count: int(z.Score),
		})
	}
	return entries
}

func (c *RedisQueryStatsRepository) IncrementQuery(ctx context.Context, query string, timestamp time.Time) error {
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

func (c *RedisQueryStatsRepository) getQueries(ctx context.Context, key string) ([]SearchStats, error) {
	results, err := c.rdb.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	return zsetToEntries(results), nil
}

func (c *RedisQueryStatsRepository) GetTopQueries(ctx context.Context, limit int64) ([]SearchStats, error) {
	results, err := c.rdb.ZRevRangeWithScores(ctx, globalKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	return zsetToEntries(results), nil
}

func (c *RedisQueryStatsRepository) deleteMinute(ctx context.Context, key string) error {
	queries, err := c.getQueries(ctx, key)
	if err != nil {
		return err
	}

	for _, entry := range queries {
		err = c.rdb.ZIncrBy(ctx, globalKey, -float64(entry.Count), entry.Query).Err()
		if err != nil {
			return err
		}
	}

	err = c.rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisQueryStatsRepository) setLastDeletedMinute(ctx context.Context, timestamp time.Time) error {
	timestamp = NormalizeMinute(timestamp)
	return c.rdb.Set(ctx, lastDeletedMinuteKey, minuteUnix(timestamp), 0).Err()
}

func (c *RedisQueryStatsRepository) getLastDeletedMinute(ctx context.Context) (time.Time, error) {
	result, err := c.rdb.Get(ctx, lastDeletedMinuteKey).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return time.Unix(result*60, 0), nil
}

func (c *RedisQueryStatsRepository) Delete(ctx context.Context, to time.Time) error {
	from, err := c.getLastDeletedMinute(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last deleted minute: %w", err)
	}

	if from.IsZero() {
		from = to
	}

	for t := from.Add(time.Minute); !t.After(to); t = t.Add(time.Minute) {
		err := c.deleteMinute(ctx, formatKey(t))
		if err != nil {
			return fmt.Errorf("failed to delete minute: %w", err)
		}
	}

	err = c.setLastDeletedMinute(ctx, to)
	if err != nil {
		return fmt.Errorf("failed to set last deleted minute: %w", err)
	}

	return nil
}
