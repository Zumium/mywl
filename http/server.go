package http

import (
	"github.com/labstack/echo"
)

type Server struct {
	httpserver *echo.Echo
	listenAddr string
}

func (s *Server) Start() {
	s.httpserver.Start(s.listenAddr)
}
