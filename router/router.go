package router

import "github.com/labstack/echo"

type Router interface {
	Healthcheck(c echo.Context) error
	Event(c echo.Context) error
}
