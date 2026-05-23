package repository

import "context"

type GuardRepository interface {
	IncrementRequest(ctx context.Context, userID string) error
	ShouldSkip(ctx context.Context, userID string) (bool, error)
}
