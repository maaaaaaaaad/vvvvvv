package mysql

import (
	"context"
	"database/sql"
	"errors"
	"jello-mark-backend/internal/domain"
	"jello-mark-backend/internal/ports/outbound"
)

type UserRepo struct {
	db *sql.DB
}

// Compile-time interface check
var _ outbound.UserRepository = (*UserRepo)(nil)

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(
	ctx context.Context,
	u domain.User,
) error {
	if u.ID == "" {
		return errors.New("id required")
	}
	q := "INSERT INTO users(id, name, created_at) VALUES(?, ?, ?)"
	_, err := r.db.ExecContext(
		ctx,
		q,
		string(u.ID),
		u.Name,
		u.CreatedAt.UTC(),
	)
	return err
}

func (r *UserRepo) List(
	ctx context.Context,
	limit int,
) (
	[]domain.User,
	error,
) {
	q := "SELECT id, name, created_at FROM users ORDER BY created_at DESC LIMIT ?"
	rows, err := r.db.QueryContext(
		ctx,
		q,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make(
		[]domain.User,
		0,
		limit,
	)
	for rows.Next() {
		var id string
		var name string
		var createdAt sql.NullTime
		if err := rows.Scan(
			&id,
			&name,
			&createdAt,
		); err != nil {
			return nil, err
		}
		u := domain.User{
			ID:   domain.UserID(id),
			Name: name,
		}
		if createdAt.Valid {
			u.CreatedAt = createdAt.Time.UTC()
		}
		res = append(
			res,
			u,
		)
	}
	return res, rows.Err()
}
