package datastore_test

// func TestUserDataStore_CreateUser(t *integration_tests.T) {
// 	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@127.0.0.1:5444/furutsu")
// 	assert.NoError(t, err)
//
// 	ds := datastore.NewUserDataStore(conn)
//
// 	u := user.User{
// 		Username: "some_user3",
// 		Password: "pass",
// 		FullName: "some suer",
// 	}
//
// 	createdUser, err := ds.CreateUser(u, context.Background())
// 	assert.NoError(t, err)
// 	assert.Equal(t, u.Username, createdUser.Username)
// 	assert.Equal(t, u.Password, createdUser.Password)
// 	assert.Equal(t, u.FullName, createdUser.FullName)
//
// 	t.Log(createdUser)
// }
