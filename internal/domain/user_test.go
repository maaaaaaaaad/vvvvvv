package domain

import (
	"testing"
	"time"
)

func TestNewUser_Success(t *testing.T) {
	now := time.Date(
		2024,
		1,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	u, err := NewUser(
		UserID("id-1"),
		"alice",
		now,
	)
	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
	if u.ID != UserID("id-1") || u.Name != "alice" || !u.CreatedAt.Equal(now) {
		t.Fatalf(
			"unexpected user: %#v",
			u,
		)
	}
}

func TestNewUser_EmptyName(t *testing.T) {
	_, err := NewUser(
		UserID("id-1"),
		"",
		time.Now(),
	)
	if err == nil {
		t.Fatalf("expected error")
	}
}
