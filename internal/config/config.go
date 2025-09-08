package config

import (
	"fmt"
	"os"
	"strconv"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvProd  Env = "prod"
)

type DBConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
}

type AppConfig struct {
	Env      Env
	HTTPPort int
	DB       DBConfig
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getint(
	key string,
	def int,
) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func Load() (
	AppConfig,
	error,
) {
	env := Env(
		getenv(
			"APP_ENV",
			string(EnvLocal),
		),
	)
	if env != EnvLocal && env != EnvProd {
		return AppConfig{}, fmt.Errorf(
			"invalid APP_ENV: %s",
			env,
		)
	}
	cfg := AppConfig{
		Env: env,
		HTTPPort: getint(
			"PORT",
			8080,
		),
		DB: DBConfig{
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
				"mad",
			),
			Password: getenv(
				"DB_PASSWORD",
				"",
			),
			Name: getenv(
				"DB_NAME",
				"main",
			),
			MaxOpenConns: getint(
				"DB_MAX_OPEN_CONNS",
				10,
			),
			MaxIdleConns: getint(
				"DB_MAX_IDLE_CONNS",
				5,
			),
		},
	}
	return cfg, nil
}

func (c AppConfig) Addr() string {
	return fmt.Sprintf(
		":%d",
		c.HTTPPort,
	)
}
