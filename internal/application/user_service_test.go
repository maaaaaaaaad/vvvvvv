package application

import (
	"context"
	"jello-mark-backend/internal/adapter/outbound/memory"
	"jello-mark-backend/internal/domain"
	"testing"
	"time"
)

type fixedClock struct{ t time.Time }

func (f fixedClock) Now() time.Time { return f.t }

func TestUserService_CreateAndList(t *testing.T) {
	repo := memory.NewUserRepo()
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
