package floodctrl

import (
	"context"
	"time"
)

type Repository interface {
	GetRequests(ctx context.Context, userID int64) (int, error)
	SetUserSession(ctx context.Context, userID int64, expiration time.Duration) error
	IncrUserRequests(ctx context.Context, userID int64) (int64, error)
}
