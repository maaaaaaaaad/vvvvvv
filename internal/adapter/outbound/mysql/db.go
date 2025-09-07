package mysql

import (
	"database/sql"
	"fmt"
	"jello-mark-backend/internal/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(c config.DBConfig) (
	*sql.DB,
	error,
) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4,utf8&loc=UTC",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
	db, err := sql.Open(
		"mysql",
		dsn,
	)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetConnMaxLifetime(1 * time.Hour)
	deadline := time.Now().Add(60 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		if err := db.Ping(); err == nil {
			return db, nil
		} else {
			lastErr = err
			time.Sleep(1 * time.Second)
		}
	}
	_ = db.Close()
	return nil, lastErr
}
