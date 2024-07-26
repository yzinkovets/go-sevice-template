package http

import (
	"go-service-template/config"
	"go-service-template/services"

	"github.com/gofiber/fiber/v3"
)

type Server struct {
	cfg         *config.ServerConfig
	app         *fiber.App
	someService *services.SomeService
}

func NewServer(cfg *config.ServerConfig, someService *services.SomeService) (*Server, error) {
	s := &Server{
		cfg:         cfg,
		someService: someService,
	}

	// Initialize a new Fiber app
	s.app = fiber.New()

	// health check
	s.app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	// main routes
	if err := s.SetupRoutes(&cfg.JwtAuthConfig); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) Start() error {
	return s.app.Listen(s.cfg.Addr)
}

func (s *Server) Shutdown() {
	s.app.Shutdown()
}
