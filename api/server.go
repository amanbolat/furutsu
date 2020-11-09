package api

import (
	"context"
	"fmt"
	"github.com/amanbolat/furutsu/services/ordersrv"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"os/signal"
	"time"

	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/config"
	"github.com/amanbolat/furutsu/services/authsrv"
	"github.com/amanbolat/furutsu/services/cartsrv"
	"github.com/amanbolat/furutsu/services/productsrv"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *Router
	logger *logrus.Logger
}

func NewServer(cfg config.Config, logger *logrus.Logger) (*Server, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.DbConnString)
	if err != nil {
		return nil, err
	}
	connConfig.MaxConns = 10

	conn, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, err
	}

	productSrv := productsrv.NewProductService(conn)
	authSrv := authsrv.NewAuthService(conn)
	cartSrv := cartsrv.NewCartService(datastore.NewPgxConn(conn))
	orderSrv := ordersrv.NewService(datastore.NewPgxConn(conn))
	r := NewRouter(RouterConfig{
		ProductService: productSrv,
		AuthService:    authSrv,
		CartService:    cartSrv,
		OrderService:   orderSrv,
	})

	s := Server{
		router: r,
		logger: logger,
	}

	return &s, nil
}

func (s Server) Start(port int) {
	go func() {
		if err := s.router.e.Start(fmt.Sprintf(":%d", port)); err != nil {
			s.logger.WithError(err).Warn()
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.router.e.Shutdown(ctx); err != nil {
		s.logger.WithError(err).Warn()
	}
}
