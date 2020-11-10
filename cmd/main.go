package main

import (
	"github.com/amanbolat/furutsu/api"
	"github.com/amanbolat/furutsu/internal/config"
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

	server, err := api.NewServer(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	server.Start(cfg.Port)
}
