package api

import (
	"context"
	"github.com/amanbolat/furutsu/services/productsrv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
)

type Router struct {
	e *echo.Echo
}

func NewRouter(productSrv *productsrv.Service) *Router {
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

	// authGroup := e.Group("")
	// authGroup.Use()


	return &Router{e: e}
}

type JSONResponse struct {
	Data interface{} `json:"data"`
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
