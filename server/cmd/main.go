package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/amanbolat/furutsu/api"
	"github.com/amanbolat/furutsu/internal/config"
	"github.com/avast/retry-go"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var cfg config.Config
var logger *logrus.Logger

func main() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})

	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Fatalf("could not parse env vars: %v", err)
	}

	var m *migrate.Migrate
	err = retry.Do(func() error {
		m, err = migrate.New(
			fmt.Sprintf("file://%s", cfg.MigratesDir),
			cfg.DbConnString)
		if err != nil {
			logger.WithError(err).Warn("could not create migrate instance, retry in 5 seconds")
			return err
		}
		return nil
	}, retry.Attempts(5), retry.Delay(time.Second*3))
	if err != nil {
		logger.Fatalf("could not create migrate instance: %v", err)
	}

	m, err = migrate.New(
		fmt.Sprintf("file://%s", cfg.MigratesDir),
		cfg.DbConnString)
	if err != nil {
		logger.Fatalf("could not create migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Fatal(err)
	}

	server, err := api.NewServer(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	server.Start(cfg.Port)
}
