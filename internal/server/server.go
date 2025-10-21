package server

import (
	"xaxaton/internal/handler"
	"xaxaton/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	server *fiber.App
}

func NewServer() *Server {
	return &Server{
		server: fiber.New(fiber.Config{
			ErrorHandler: fiber.DefaultErrorHandler,
		}),
	}
}

func (s *Server) Run(addr string) error {
	return s.server.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown()
}

func (s *Server) InitRoutes(handler *handler.Handler, middleware middleware.Middleware) {
	s.server.Use(middleware.RecoverMiddleware())
}
