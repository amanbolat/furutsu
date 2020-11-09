package integration_test

import (
	"context"
	"fmt"
	"github.com/amanbolat/furutsu/datastore"
	"github.com/amanbolat/furutsu/internal/cart"
	"github.com/amanbolat/furutsu/internal/discount"
	"github.com/amanbolat/furutsu/internal/product"
	"github.com/amanbolat/furutsu/internal/user"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var (
	conn        *pgxpool.Pool
	allProducts []product.Product

	appleId  string
	pearId   string
	bananaId string
	orangeId string

	testUser = user.User{
		Username: "test_user",
		Password: "password",
		FullName: "John Doe",
	}

	orangeCoupon discount.Coupon
)

func TestMain(m *testing.M) {
	fmt.Println(os.Getenv("TEST_DB_URL"))
	retryCount := 30

	for {
		var err error
		conn, err = pgxpool.Connect(context.Background(), os.Getenv("TEST_DB_URL"))
		if err != nil {
			if retryCount == 0 {
				logrus.WithError(err).Fatal("no more retries")
			}

			logrus.WithError(err).Infof("could not connect to db, wait 2 seconds. %d retries left", retryCount)
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	file, err := ioutil.ReadFile("../sql/bootstrap.sql")
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	_, err = conn.Exec(context.Background(), string(file))
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	m.Run()
}

func TestCreateAndListProducts(t *testing.T) {
	pds := datastore.NewProductDataStore(datastore.NewPgxConn(conn))
	list := []product.Product{
		{
			Name:        "apple",
			Price:       200,
			Description: "apple",
		},
		{
			Name:        "pear",
			Price:       300,
			Description: "pear",
		},
		{
			Name:        "banana",
			Price:       400,
			Description: "banana",
		},
		{
			Name:        "orange",
			Price:       500,
			Description: "orange",
		},
	}

	for _, p := range list {
		_, err := pds.CreateProduct(p, context.Background())
		require.NoError(t, err)
	}

	found, err := pds.ListProducts(context.Background())
	require.NoError(t, err)
	assert.Len(t, found, len(list))
	allProducts = append(allProducts, found...)
	assert.Len(t, allProducts, len(list))
}

func TestFindProductsById(t *testing.T) {
	pds := datastore.NewProductDataStore(datastore.NewPgxConn(conn))
	for _, p := range allProducts {
		found, err := pds.GetProductById(p.ID, context.Background())
		require.NoError(t, err)
		switch found.Name {
		case "apple":
			appleId = found.ID
		case "pear":
			pearId = found.ID
		case "banana":
			bananaId = found.ID
		case "orange":
			orangeId = found.ID
		}
	}
}

func TestCreateUser(t *testing.T) {
	uds := datastore.NewUserDataStore(datastore.NewPgxConn(conn))

	created, err := uds.CreateUser(testUser, context.Background())
	require.NoError(t, err)
	assert.Equal(t, testUser.Username, created.Username)
	assert.Equal(t, testUser.Password, created.Password)
	assert.Equal(t, testUser.FullName, created.FullName)

	t.Run("find created user", func(t *testing.T) {
		found, err := uds.GetUserByUsername(testUser.Username, context.Background())
		require.NoError(t, err)
		assert.Equal(t, testUser.Username, found.Username)
		testUser.Id = found.Id
	})
}

func TestCreateCoupon(t *testing.T) {
	dds := datastore.NewDiscountDataStore(datastore.NewPgxConn(conn))

	c := discount.Coupon{
		Code: "cup123",
		Name: "10 oranges discount 30%",
		Rule: discount.RuleItemsAll{
			ProductID: orangeId,
			Amount:    10,
		},
		Percent: 30,
		Expire:  time.Now().AddDate(1, 0, 0),
	}

	created, err := dds.CreateCoupon(c, context.Background())
	require.NoError(t, err)
	assert.Equal(t, c.Code, created.Code)
	assert.Equal(t, c.Name, created.Name)
	assert.Equal(t, c.Rule, created.Rule)
	assert.Equal(t, c.Percent, created.Percent)
	assert.Equal(t, c.Expire, created.Expire)
	orangeCoupon = created
}

func TestCreateDiscounts(t *testing.T) {
	dds := datastore.NewDiscountDataStore(datastore.NewPgxConn(conn))

	d1 := discount.Discount{
		Name: "7 apples 10%",
		Rule: discount.RuleItemsAll{
			ProductID: appleId,
			Amount:    7,
		},
		Percent: 10,
	}

	d2 := discount.Discount{
		Name: "4 pears, 2 bananas - 10%",
		Rule: discount.RuleItemsSet{
			ItemSet: map[string]int{pearId: 4, bananaId: 2},
		},
		Percent: 30,
	}

	createdD1, err := dds.CreateDiscount(d1, context.Background())
	require.NoError(t, err)
	assert.Equal(t, d1, createdD1)
	createdD2, err := dds.CreateDiscount(d2, context.Background())
	require.NoError(t, err)
	assert.Equal(t, d2, createdD2)

	t.Run("find created discounts", func(t *testing.T) {
		found, err := dds.ListDiscounts(context.Background())
		require.NoError(t, err)
		assert.Len(t, found, 2)
	})
}

func TestCreateCart(t *testing.T) {
	cds := datastore.NewCartDataStore(datastore.NewPgxConn(conn))
	created, err := cds.CreateCart(testUser.Id, context.Background())
	require.NoError(t, err)

	assert.NotEmpty(t, created.Id)
	assert.Empty(t, created.Items)
	assert.Empty(t, created.DiscountSets)
	assert.Empty(t, created.Coupons)
	assert.Empty(t, created.NonDiscountSet)
}

func TestGetCartForUser(t *testing.T) {
	cds := datastore.NewCartDataStore(datastore.NewPgxConn(conn))
	found, err := cds.GetCartForUser(testUser.Id, context.Background())
	require.NoError(t, err)
	require.NotEqual(t, cart.Cart{}, found)
}

func TestCreateAndGetCartItem(t *testing.T) {
	cds := datastore.NewCartDataStore(datastore.NewPgxConn(conn))
	userCart, err := cds.GetCartForUser(testUser.Id, context.Background())
	require.NoError(t, err)
	assert.Equal(t, testUser.Id, userCart.UserId)
	err = cds.CreateCartItem(userCart.Id, appleId, 20, context.Background())
	require.NoError(t, err)
	err = cds.CreateCartItem(userCart.Id, bananaId, 20, context.Background())
	require.NoError(t, err)
	err = cds.CreateCartItem(userCart.Id, pearId, 20, context.Background())
	require.NoError(t, err)
	err = cds.CreateCartItem(userCart.Id, orangeId, 20, context.Background())
	require.NoError(t, err)

	updatedCart, err := cds.GetCartForUser(testUser.Id, context.Background())
	require.NoError(t, err)
	assert.Len(t, updatedCart.Items, 4)
}

func TestUpdateCartItems(t *testing.T) {
	cds := datastore.NewCartDataStore(datastore.NewPgxConn(conn))
	userCart, err := cds.GetCartForUser(testUser.Id, context.Background())
	require.NoError(t, err)
	err = cds.SetCartItemAmount(userCart.Id, appleId, 10, context.Background())
	require.NoError(t, err)
	err = cds.DeleteCartItem(userCart.Id, orangeId, context.Background())
	require.NoError(t, err)

	updatedCart, err := cds.GetCartForUser(testUser.Id, context.Background())
	require.NoError(t, err)

	assert.Len(t, updatedCart.Items, 3)
}

func TestCartAttachDetachCoupon(t *testing.T) {
	cds := datastore.NewCartDataStore(datastore.NewPgxConn(conn))
	err := cds.AttachCouponToCart(testUser.Id, orangeCoupon.ID, context.Background())
	require.NoError(t, err)
	userCart, err := cds.GetCartForUser(testUser.Id, context.Background())
	require.NoError(t, err)
	assert.Len(t, userCart.Coupons, 1)
	assert.Equal(t, userCart.Coupons[0].GetName(), orangeCoupon.Name)
	assert.Equal(t, userCart.Coupons[0].GetExpireTime(), orangeCoupon.Expire)
	assert.Equal(t, userCart.Coupons[0].GetPercentage(), orangeCoupon.Percent)
}

func TestClearCart(t *testing.T) {
	cds := datastore.NewCartDataStore(datastore.NewPgxConn(conn))
	userCart, err := cds.GetCartForUser(testUser.Id, context.Background())
	require.NoError(t, err)
	err = cds.ClearCart(userCart.Id, context.Background())
	require.NoError(t, err)
	clearedCart, err := cds.GetCartForUser(testUser.Id, context.Background())
	require.NoError(t, err)
	assert.Len(t, clearedCart.Items, 0)
}
