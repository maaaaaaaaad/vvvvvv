package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	err := os.Unsetenv("APP_ENV")
	if err != nil {
		return
	}
	err = os.Unsetenv("PORT")
	if err != nil {
		return
	}
	err = os.Unsetenv("DB_HOST")
	if err != nil {
		return
	}
	cfg, err := Load()
	if err != nil {
		t.Fatalf(
			"load error: %v",
			err,
		)
	}
	if cfg.Env != EnvLocal || cfg.HTTPPort != 8080 {
		t.Fatalf(
			"unexpected cfg: %#v",
			cfg,
		)
	}
}

func TestLoad_InvalidEnv(t *testing.T) {
	err := os.Setenv(
		"APP_ENV",
		"bad",
	)
	if err != nil {
		return
	}
	defer func() {
		err := os.Unsetenv("APP_ENV")
		if err != nil {

		}
	}()
	_, err = Load()
	if err == nil {
		t.Fatalf("expected error")
	}
}
