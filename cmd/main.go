package main

import (
	"flag"
	"github.com/amanbolat/furutsu/api"
	"github.com/amanbolat/furutsu/internal/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

var cfg config.Config
var logger *logrus.Logger

func main()  {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	cfgFile := flag.String("config", "config", "path to the configuration file")
	flag.Parse()

	err := godotenv.Load(*cfgFile)
	if err != nil {
		logger.Fatalf("could not load env file\n%v", err)
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		logger.Fatalf("could not parse env vars: %v", err)
	}

	server, err := api.NewServer(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	server.Start(cfg.Port)
}
