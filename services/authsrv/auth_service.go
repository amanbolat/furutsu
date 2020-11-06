package authsrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"time"
)

var JwtSecret = []byte("pan4HAPPENED8archaic2prolix")

type Claims struct {
	UserId string
	Username string
	jwt.StandardClaims
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Service struct {
	dbConn *pgx.Conn
}

func NewAuthService(conn *pgx.Conn) *Service {
	return &Service{dbConn: conn}
}

func (s Service) Login(creds Credentials, ctx context.Context) (string, error) {
	ds := datastore.NewUserDataStore(s.dbConn)
	u, err := ds.GetUserByUsername(creds.Username, ctx)
	if err != nil {
		return "", errors.New(fmt.Sprintf("wrong credentials: %v", err))
	}
	
	if u.Password != creds.Password {
		return "", errors.New("wrong credentials")
	}
	
	claims := &Claims{
		UserId:         u.Id,
		Username:       u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", errors.New(fmt.Sprintf("internal error: %v", err))
	}
	
	return tokenStr, nil
}

func (s Service) Register(u user.User, ctx context.Context) error {
	if err := u.Validate(); err != nil {
		return err
	}

	tx, err := s.dbConn.Begin(ctx)
	if err != nil {
		return err
	}

	ds := datastore.NewUserDataStore(tx)
	createdUser, err := ds.CreateUser(u, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	cartDs := datastore.NewCartDataStore(tx)

	_, err = cartDs.CreateCart(createdUser.Id, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}