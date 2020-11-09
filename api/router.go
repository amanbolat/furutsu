package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/amanbolat/furutsu/internal/apperr"
	"github.com/amanbolat/furutsu/services/ordersrv"
	"net/http"
	"os"

	"github.com/amanbolat/furutsu/internal/user"
	"github.com/amanbolat/furutsu/services/authsrv"
	"github.com/amanbolat/furutsu/services/cartsrv"
	"github.com/amanbolat/furutsu/services/productsrv"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	e *echo.Echo
}

type RouterConfig struct {
	ProductService *productsrv.Service
	AuthService    *authsrv.Service
	CartService    *cartsrv.Service
	OrderService   *ordersrv.Service
}

func NewRouter(cfg RouterConfig) *Router {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)
		if errors.As(err, &apperr.Error{}) {
			err = c.JSON(echo.ErrInternalServerError.Code, err)
			if err != nil {
				c.Logger().Error(err)
			}
		} else {
			e.DefaultHTTPErrorHandler(err, c)
		}
	}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}",` +
			`"method":"${method}","uri":"${uri}","status":${status},` +
			`"latency_human":"${latency_human}"` + "\n",
		Output: os.Stdout,
	}), middleware.Recover(), middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:80"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, "X-API-REQUEST-ID"},
		AllowCredentials: true,
	}))

	e.GET("/product", ListProducts(cfg.ProductService))
	e.POST("/auth/login", Login(cfg.AuthService))
	e.POST("/auth/register", Register(cfg.AuthService))

	authGroup := e.Group("")
	authGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		ErrorHandler: func(err error) error {
			return echo.ErrUnauthorized
		},
		SuccessHandler: func(c echo.Context) {
			claims := c.Get("jwt_token").(*jwt.Token).Claims.(*authsrv.Claims)
			c.Set("user", claims)
		},
		SigningKey: authsrv.JwtSecret,
		Claims:     &authsrv.Claims{},
		ContextKey: "jwt_token",
	}))
	authGroup.GET("/cart", GetCart(cfg.CartService))
	authGroup.POST("/cart/product", SetCartItemAmount(cfg.CartService))
	authGroup.POST("/cart/coupon", ApplyCouponToCart(cfg.CartService))
	authGroup.DELETE("/cart/coupon/:coupon_code", DetachCouponFromCart(cfg.CartService))
	authGroup.PUT("/order", CreateOrder(cfg.OrderService))
	authGroup.GET("/order", ListOrders(cfg.OrderService))
	authGroup.GET("/order/:id", GetOrderById(cfg.OrderService))
	authGroup.POST("/payment/pay/:order_id", nil)
	authGroup.GET("/coupon", nil)

	return &Router{e: e}
}

type JSONResponse struct {
	Data interface{} `json:"data"`
}

func GetOrderById(srv *ordersrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("user").(*authsrv.Claims)
		id := c.Param("id")

		o, err := srv.GetOrderById(id, claims.Id, context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{Data: o})
	}
}

func ListOrders(srv *ordersrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("user").(*authsrv.Claims)

		ol, err := srv.ListOrders(claims.Id, context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{Data: ol})
	}
}

func CreateOrder(srv *ordersrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("user").(*authsrv.Claims)

		o, err := srv.CreateOrder(claims.Id, context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{Data: o})
	}
}

func DetachCouponFromCart(srv *cartsrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("user").(*authsrv.Claims)
		code := c.Param("coupon_code")

		userCart, err := srv.DetachCoupon(claims.Id, code, context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{Data: userCart})
	}
}

type ApplyCouponRequest struct {
	Code string `json:"code"`
}

func ApplyCouponToCart(srv *cartsrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("user").(*authsrv.Claims)
		var req ApplyCouponRequest
		err := c.Bind(&req)
		if err != nil {
			return err
		}

		userCart, err := srv.ApplyCoupon(claims.Id, req.Code, context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{Data: userCart})
	}
}

type SetCartItemRequest struct {
	ProductId string `json:"product_id"`
	Amount    int    `json:"amount"`
}

func SetCartItemAmount(srv *cartsrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("user").(*authsrv.Claims)
		var req SetCartItemRequest
		err := c.Bind(&req)
		if err != nil {
			return echo.ErrInternalServerError
		}
		updatedCart, err := srv.SetItemAmount(req.ProductId, claims.Id, req.Amount, context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{
			Data: updatedCart,
		})
	}
}

func GetCart(srv *cartsrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("user").(*authsrv.Claims)

		userCart, err := srv.GetCart(claims.Id, context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{Data: userCart})
	}
}

func ListProducts(srv *productsrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		products, err := srv.ListProducts(context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{Data: products})
	}
}

func Login(srv *authsrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		creds := authsrv.Credentials{}
		err := c.Bind(&creds)
		c.Logger().Error(err)
		if err != nil {
			return err
		}

		token, err := srv.Login(creds, context.Background())
		c.Logger().Error(err)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{Data: token})
	}
}

func Register(srv *authsrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := user.User{}
		err := c.Bind(&u)
		if err != nil {
			return echo.ErrBadRequest
		}
		fmt.Println(u)
		err = srv.Register(u, context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{
			Data: "ok",
		})
	}
}
