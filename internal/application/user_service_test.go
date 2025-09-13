package application

import (
	"context"
	"errors"
	"jello-mark-backend/internal/domain"
	"sync"
	"testing"
	"time"
)

type fixedClock struct{ t time.Time }

type testUserRepo struct {
	mu    sync.RWMutex
	items []domain.User
}

func newTestUserRepo() *testUserRepo {
	return &testUserRepo{
		items: make(
			[]domain.User,
			0,
			16,
		),
	}
}

func (r *testUserRepo) Create(
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

func (r *testUserRepo) List(
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

func (f fixedClock) Now() time.Time { return f.t }

func TestUserService_CreateAndList(t *testing.T) {
	repo := newTestUserRepo()
	fc := fixedClock{
		t: time.Date(
			2024,
			1,
			2,
			3,
			4,
			5,
			0,
			time.UTC,
		),
	}
	id := domain.UserID("id-1")
	idgen := func() domain.UserID { return id }
	svc := NewUserService(
		repo,
		fc,
		idgen,
	)
	ctx := context.Background()
	u, err := svc.Create(
		ctx,
		"bob",
	)
	if err != nil {
		t.Fatalf(
			"create error: %v",
			err,
		)
	}
	if u.ID != id || u.Name != "bob" || !u.CreatedAt.Equal(fc.t) {
		t.Fatalf(
			"unexpected user: %#v",
			u,
		)
	}
	list, err := svc.List(
		ctx,
		10,
	)
	if err != nil {
		t.Fatalf(
			"list error: %v",
			err,
		)
	}
	if len(list) != 1 || list[0].ID != id {
		t.Fatalf(
			"unexpected list: %#v",
			list,
		)
	}
}
