package authsrv

import (
	"context"
	"time"

	"github.com/amanbolat/furutsu/internal/apperr"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/user"
	"github.com/dgrijalva/jwt-go"
)

var JwtSecret = []byte("pan4HAPPENED8archaic2prolix")

const sessionDuration = time.Hour * 24

type Claims struct {
	Id                 string `json:"user_id"`
	Username           string `json:"username"`
	FullName           string `json:"full_name"`
	jwt.StandardClaims `json:""`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Service struct {
	dbConn *pgxpool.Pool
}

func NewAuthService(conn *pgxpool.Pool) *Service {
	return &Service{dbConn: conn}
}

func (s Service) Login(creds Credentials, ctx context.Context) (string, error) {
	ds := datastore.NewUserDataStore(s.dbConn)
	u, err := ds.GetUserByUsername(creds.Username, ctx)
	if err != nil {
		return "", apperr.With(err, "Wrong credentials", "Username or password is incorrect")
	}

	if u.Password != creds.Password {
		return "", apperr.With(err, "Wrong credentials", "Username or password is incorrect")
	}

	claims := &Claims{
		Id:       u.Id,
		Username: u.Username,
		FullName: u.FullName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(sessionDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", apperr.With(err, "Unknown error", "")
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
		return apperr.With(err, "User with this username does already exist", "")
	}

	cartDs := datastore.NewCartDataStore(tx)

	_, err = cartDs.CreateCart(createdUser.Id, ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
