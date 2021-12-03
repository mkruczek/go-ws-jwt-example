package api

import (
	"github.com/labstack/echo/v4"
	"github.com/mkruczek/go-ws-jwt-example/internal/handlers"
)

type HttpServer struct {
	main *echo.Echo
}

func NewServer(e *echo.Echo) *HttpServer {
	return &HttpServer{main: e}
}

func (s HttpServer) RegisterRouters() {
	s.main.GET("/", handlers.Home)
	s.main.GET("/ws", handlers.WsEndpoint)
}

func (s HttpServer) Start() {
	s.main.Start(":8888")
}
