package domain

import (
	"errors"
	"time"
)

type UserID string

type User struct {
	ID        UserID
	Name      string
	CreatedAt time.Time
}

func NewUser(
	id UserID,
	name string,
	now time.Time,
) (
	User,
	error,
) {
	if name == "" {
		return User{}, errors.New("name required")
	}
	return User{
		ID:        id,
		Name:      name,
		CreatedAt: now.UTC(),
	}, nil
}
