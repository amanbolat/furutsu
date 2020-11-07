package datastore_test

import (
	"context"
	"github.com/jackc/pgx/v4"
	"testing"

	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestUserDataStore_CreateUser(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@127.0.0.1:5444/furutsu")
	assert.NoError(t, err)

	ds := datastore.NewUserDataStore(conn)

	u := user.User{
		Username: "some_user3",
		Password: "pass",
		FullName: "some suer",
	}

	createdUser, err := ds.CreateUser(u, context.Background())
	assert.NoError(t, err)
	assert.Equal(t, u.Username, createdUser.Username)
	assert.Equal(t, u.Password, createdUser.Password)
	assert.Equal(t, u.FullName, createdUser.FullName)

	t.Log(createdUser)
}
