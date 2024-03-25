package usecase

import (
	"context"
	"errors"
	"sync"
	"task/config"
	"task/internal/floodctrl"
	"task/pkg/db"
)

type floodUC struct {
	repo floodctrl.Repository
	mx   sync.Mutex
	cfg  *config.Config
}

func NewFloodUC(repository floodctrl.Repository, cfg *config.Config) floodctrl.FloodControl {
	return &floodUC{repo: repository, mx: sync.Mutex{}, cfg: cfg}
}

func (f *floodUC) Check(ctx context.Context, userID int64) (bool, error) {
	f.mx.Lock()
	defer f.mx.Unlock()
	r, err := f.repo.GetRequests(ctx, userID)
	if err != nil {
		if errors.Is(err, db.NotFoundObject) {
			err = f.repo.SetUserSession(ctx, userID, f.cfg.Flood.TimeLimit)
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, err
	}

	if r+1 > f.cfg.Flood.MaxRequests {
		return false, nil
	}

	_, err = f.repo.IncrUserRequests(ctx, userID)

	return true, err
}
