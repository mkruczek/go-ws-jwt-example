package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mkruczek/go-ws-jwt-example/internal/api"
)

func main() {

	e := echo.New()
	s := api.NewServer(e)
	s.RegisterRouters()

	s.Start()
}
