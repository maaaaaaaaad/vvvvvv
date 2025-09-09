package memory

import (
	"context"
	"errors"
	"jello-mark-backend/internal/domain"
	"jello-mark-backend/internal/ports/outbound"
	"sync"
)

type UserRepo struct {
	mu    sync.RWMutex
	items []domain.User
}

var _ outbound.UserRepository = (*UserRepo)(nil)

func NewUserRepo() *UserRepo {
	return &UserRepo{
		items: make(
			[]domain.User,
			0,
			16,
		),
	}
}

func (r *UserRepo) Create(
	ctx context.Context,
	u domain.User,
) error {
	if u.ID == "" {
		return errors.New("id required")
	}
	r.mu.Lock()
	r.items = append(
		r.items,
		u,
	)
	r.mu.Unlock()
	return nil
}

func (r *UserRepo) List(
	ctx context.Context,
	limit int,
) (
	[]domain.User,
	error,
) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if limit > len(r.items) {
		limit = len(r.items)
	}
	res := make(
		[]domain.User,
		limit,
	)
	copy(
		res,
		r.items[:limit],
	)
	return res, nil
}
