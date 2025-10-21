package server

import (
	"xaxaton/internal/handler"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	server *fiber.App
}

func NewServer() *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: fiber.DefaultErrorHandler,
	})

	app.Use(RecoverMiddleware)

	return &Server{
		server: app,
	}
}

func RecoverMiddleware(ctx *fiber.Ctx) error {
	defer func() {
		if err := recover(); err != nil {
			ctx.Context().Logger().Printf("%v", err)
			ctx.Context().SetStatusCode(fiber.ErrInternalServerError.Code)
		}
	}()

	return ctx.Next()
}

func (s *Server) Run(addr string) error {
	return s.server.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown()
}

func (s *Server) InitRoutes(handler *handler.Handler) {
	s.server.Get("/", handler.Template)
}
