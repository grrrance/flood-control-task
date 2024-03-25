package repository

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
	"task/internal/floodctrl"
	"task/pkg/db"
	"time"
)

type floodRepo struct {
	redisClient *redis.Client
}

func NewFloodRepository(redisClient *redis.Client) floodctrl.Repository {
	return &floodRepo{redisClient: redisClient}
}

func (f *floodRepo) GetRequests(ctx context.Context, userID int64) (int, error) {
	r, err := f.redisClient.Get(ctx, convertIDToStr(userID)).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, errors.Wrap(db.NotFoundObject, "floodRepo.GetRequests.Get")
		}
		return 0, errors.Wrap(err, "floodRepo.GetRequests.Get")
	}

	return r, nil
}

func (f *floodRepo) SetUserSession(ctx context.Context, userID int64, expiration time.Duration) error {
	return errors.Wrap(f.redisClient.Set(ctx, convertIDToStr(userID), 1, expiration).Err(), "floodRepo.SetUserSession.Set")
}

func (f *floodRepo) IncrUserRequests(ctx context.Context, userID int64) (int64, error) {
	result, err := f.redisClient.Incr(ctx, convertIDToStr(userID)).Result()
	if err != nil {
		return 0, errors.Wrap(err, "floodRepo.IncrUserRequests.Incr")
	}
	return result, nil
}

func convertIDToStr(id int64) string {
	return strconv.FormatInt(id, 10)
}
