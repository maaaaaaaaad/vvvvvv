package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"jello-mark-backend/internal/adapter/inbound/graphql"
	mem "jello-mark-backend/internal/adapter/outbound/memory"
	mysqlrepo "jello-mark-backend/internal/adapter/outbound/mysql"
	"jello-mark-backend/internal/adapter/outbound/system"
	"jello-mark-backend/internal/application"
	"jello-mark-backend/internal/config"
	"jello-mark-backend/internal/domain"
	"jello-mark-backend/internal/ports/outbound"
	"jello-mark-backend/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(
			"config error: %v",
			err,
		)
	}
	printStartupInfo(cfg)
	var userRepo outbound.UserRepository
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
				log.Printf(
					"error closing database: %v",
					err,
				)
			}
		}(db)
		userRepo = mysqlrepo.NewUserRepo(db)
	} else {
		userRepo = mem.NewUserRepo()
	}
	idgen := func() domain.UserID {
		return domain.UserID(uuid.NewString())
	}
	clk := system.NewClock()
	usvc := application.NewUserService(
		userRepo,
		clk,
		idgen,
	)

	h := graphql.NewGraphQLHandler(
		usvc,
		cfg.Env == config.EnvLocal,
	)
	mux := server.NewMux(h)

	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}
	go func() {
		printServerInfo(cfg)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(
			err,
			http.ErrServerClosed,
		) {
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

func printStartupInfo(cfg config.AppConfig) {
	log.Printf("Jello Mark Backend Starting...")
	log.Printf(
		"Environment: %s",
		cfg.Env,
	)
	log.Printf(
		"HTTP Port: %d",
		cfg.HTTPPort,
	)
	log.Printf(
		"Database: %s@%s:%d/%s",
		cfg.DB.User,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)
}

func printServerInfo(cfg config.AppConfig) {
	baseURL := fmt.Sprintf(
		"http://localhost%s",
		cfg.Addr(),
	)

	log.Println("Server is running!")
	log.Printf(
		"Health Check: %s/health",
		baseURL,
	)
	if cfg.Env == config.EnvLocal {
		log.Printf(
			"GraphQL Playground: %s/graphql",
			baseURL,
		)
	}
}
