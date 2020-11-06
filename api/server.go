package api

import (
	"context"
	"fmt"
	"github.com/amanbolat/furutsu/internal/config"
	"github.com/amanbolat/furutsu/services/authsrv"
	"github.com/amanbolat/furutsu/services/productsrv"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	router            *Router
	logger            *logrus.Logger
}

func NewServer(cfg config.Config, logger *logrus.Logger) (*Server, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DbConnString)
	if err != nil {
		return nil, err
	}

	productSrv := productsrv.NewProductService(conn)
	authSrv := authsrv.NewAuthService(conn)
	r := NewRouter(productSrv, authSrv)

	s := Server{
		router: r,
		logger: logger,
	}

	return &s, nil
}

func (s Server) Start(port int) {
	go func() {
		if err := s.router.e.Start(fmt.Sprintf(":%d", port)); err != nil {
			s.logger.Fatalf("failed to start server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.router.e.Shutdown(ctx); err != nil {
		s.logger.Fatal(err)
	}
}