package outbound

import (
	"context"
	"jello-mark-backend/internal/domain"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		u domain.User,
	) error
	List(
		ctx context.Context,
		limit int,
	) (
		[]domain.User,
		error,
	)
}
