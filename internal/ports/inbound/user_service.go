package inbound

import (
	"context"
	"jello-mark-backend/internal/domain"
)

type UserService interface {
	Create(
		ctx context.Context,
		name string,
	) (
		domain.User,
		error,
	)
	List(
		ctx context.Context,
		limit int,
	) (
		[]domain.User,
		error,
	)
}
