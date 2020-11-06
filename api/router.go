package api

import (
	"context"
	"fmt"
	"github.com/amanbolat/furutsu/internal/user"
	"github.com/amanbolat/furutsu/services/authsrv"
	"github.com/amanbolat/furutsu/services/productsrv"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
)

type Router struct {
	e *echo.Echo
}

func NewRouter(productSrv *productsrv.Service, authSrv *authsrv.Service) *Router {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)
		e.DefaultHTTPErrorHandler(err, c)
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

	e.GET("/product", ListProducts(productSrv))
	e.GET("/product/{id}", GetProductById(productSrv))
	e.POST("/auth/login", Login(authSrv))
	e.POST("/auth/register", Register(authSrv))
	authGroup := e.Group("")
	authGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		ErrorHandler: func(err error) error {
			return echo.ErrUnauthorized
		},
		SuccessHandler: func(c echo.Context) {
			claims := c.Get("jwt_token").(*jwt.Token).Claims.(*authsrv.Claims)
			c.Set("user", claims)
		},
		SigningKey:              authsrv.JwtSecret,
		Claims:                  &authsrv.Claims{},
		ContextKey: "jwt_token",
	}))
	authGroup.GET("/cart", GetCart())

	return &Router{e: e}
}

type JSONResponse struct {
	Data interface{} `json:"data"`
}

func GetCart() echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := c.Get("user").(*authsrv.Claims)
		fmt.Printf("claims: %v\n", claims)
		return c.JSON(http.StatusOK, "CART")
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

func GetProductById(srv *productsrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		product, err := srv.GetProductById(id, context.Background())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, JSONResponse{Data: product})
	}
}

func Login(srv *authsrv.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		creds := authsrv.Credentials{}
		err := c.Bind(&creds)
		c.Logger().Error(err)
		if err != nil {
			return echo.ErrUnauthorized
		}

		token, err := srv.Login(creds, context.Background())
		c.Logger().Error(err)
		if err != nil {
			return echo.ErrUnauthorized
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