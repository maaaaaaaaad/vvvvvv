package main

import (
	"context"
	"database/sql"
	"fmt"
	"jello-mark-backend/internal/adapter/inbound/graphql"
	mem "jello-mark-backend/internal/adapter/outbound/memory"
	mysqlrepo "jello-mark-backend/internal/adapter/outbound/mysql"
	"jello-mark-backend/internal/application"
	"jello-mark-backend/internal/config"
	"jello-mark-backend/internal/domain"
	"jello-mark-backend/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(
			"config error: %v",
			err,
		)
	}
	var userRepo domain.UserRepository
	if cfg.Env == config.EnvProd || cfg.Env == config.EnvLocal {
		db, err := mysqlrepo.NewDB(cfg.DB)
		if err != nil {
			log.Fatalf(
				"db error: %v",
				err,
			)
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {

			}
		}(db)
		userRepo = mysqlrepo.NewUserRepo(db)
	} else {
		userRepo = mem.NewUserRepo()
	}
	var counter uint64
	idgen := func() domain.UserID {
		n := atomic.AddUint64(
			&counter,
			1,
		)
		return domain.UserID(
			fmt.Sprintf(
				"%d-%d",
				time.Now().UTC().UnixNano(),
				n,
			),
		)
	}
	clk := application.SystemClock{}
	usvc := application.NewUserService(
		userRepo,
		clk,
		idgen,
	)

	h := graphql.Handler{Users: usvc}
	mux := server.NewMux(h)

	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}
	go func() {
		log.Printf(
			"listening on %s",
			cfg.Addr(),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(
				"server error: %v",
				err,
			)
		}
	}()

	quit := make(
		chan os.Signal,
		1,
	)
	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-quit
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf(
			"shutdown error: %v",
			err,
		)
	}
}
