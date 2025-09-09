package application

import (
	"context"
	"jello-mark-backend/internal/domain"
	"jello-mark-backend/internal/ports/inbound"
	"jello-mark-backend/internal/ports/outbound"
)

type userService struct {
	repo  outbound.UserRepository
	clock outbound.Clock
	idgen func() domain.UserID
}

func NewUserService(
	repo outbound.UserRepository,
	clock outbound.Clock,
	idgen func() domain.UserID,
) inbound.UserService {
	return &userService{
		repo:  repo,
		clock: clock,
		idgen: idgen,
	}
}

func (s *userService) Create(
	ctx context.Context,
	name string,
) (
	domain.User,
	error,
) {
	u, err := domain.NewUser(
		s.idgen(),
		name,
		s.clock.Now(),
	)
	if err != nil {
		return domain.User{}, err
	}
	if err := s.repo.Create(
		ctx,
		u,
	); err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (s *userService) List(
	ctx context.Context,
	limit int,
) (
	[]domain.User,
	error,
) {
	if limit <= 0 {
		limit = 100
	}
	return s.repo.List(
		ctx,
		limit,
	)
}
