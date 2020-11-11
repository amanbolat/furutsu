package datastore

import (
	"context"
	"github.com/amanbolat/furutsu/internal/user"
	"github.com/georgysavva/scany/pgxscan"
)

type UserDataStore struct {
	querier pgxscan.Querier
}

func NewUserDataStore(q pgxscan.Querier) *UserDataStore {
	return &UserDataStore{querier: q}
}

func (s UserDataStore) GetUserByUsername(uname string, ctx context.Context) (user.User, error) {
	var u user.User
	err := pgxscan.Get(ctx, s.querier, &u, `select * from "user" where username=$1`, uname)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (s UserDataStore) CreateUser(u user.User, ctx context.Context) (user.User, error) {
	rows, err := s.querier.Query(ctx, `INSERT INTO "user" (username, full_name, password) VALUES ($1, $2, $3) RETURNING *`, u.Username, u.FullName, u.Password)
	if err != nil {
		return user.User{}, err
	}
	defer rows.Close()

	var createdUser user.User
	err = pgxscan.ScanOne(&createdUser, rows)
	if err != nil {
		return user.User{}, err
	}

	return createdUser, nil
}
