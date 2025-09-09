package mysql

import (
	"context"
	"database/sql"
	"jello-mark-backend/internal/config"
	"jello-mark-backend/internal/domain"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestIntegration_UserRepo_MySQL(t *testing.T) {
	if os.Getenv("RUN_MYSQL_INTEGRATION") != "1" {
		t.Skip("integration disabled")
	}
	dbc := config.DBConfig{
		Host: getenv(
			"DB_HOST",
			"127.0.0.1",
		),
		Port: getint(
			"DB_PORT",
			3306,
		),
		User: getenv(
			"DB_USER",
			"root",
		),
		Password: getenv(
			"DB_PASSWORD",
			"",
		),
		Name: getenv(
			"DB_NAME",
			"app",
		),
		MaxOpenConns: 5,
		MaxIdleConns: 2,
	}
	db, err := NewDB(dbc)
	if err != nil {
		t.Fatalf(
			"db error: %v",
			err,
		)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	ctx := context.Background()
	_, _ = db.ExecContext(
		ctx,
		"CREATE DATABASE IF NOT EXISTS "+dbc.Name,
	)
	_, _ = db.ExecContext(
		ctx,
		"CREATE TABLE IF NOT EXISTS users (id VARCHAR(64) PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at DATETIME(6) NOT NULL)",
	)
	repo := NewUserRepo(db)
	uid := domain.UserID("it-" + time.Now().UTC().Format("20060102150405.000000"))
	u := domain.User{
		ID:        uid,
		Name:      "it-user",
		CreatedAt: time.Now().UTC(),
	}
	if err := repo.Create(
		ctx,
		u,
	); err != nil {
		t.Fatalf(
			"create error: %v",
			err,
		)
	}
	list, err := repo.List(
		ctx,
		20,
	)
	if err != nil {
		t.Fatalf(
			"list error: %v",
			err,
		)
	}
	found := false
	for _, x := range list {
		if x.ID == uid {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("user not found in list")
	}
}

func getenv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}

func getint(
	k string,
	d int,
) int {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return d
	}
	return i
}
