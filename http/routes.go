package http

import (
	"go-service-template/config"
	"go-service-template/http/handlers"
	"go-service-template/http/middleware"

	"github.com/sirupsen/logrus"
)

func (s *Server) SetupRoutes(authConfig *config.JwtAuthConfig) error {
	// middleware
	auth, err := middleware.NewJwtAuth(authConfig)
	if err != nil {
		logrus.Error("can't init jwt auth middleware")
		return err
	}

	// apiNoAuth := s.app.Group("/api/v1")

	// handlers that require authentication
	apiAuth := s.app.Group("/api/v1", auth.Handle)

	apiAuth.Get("/some", handlers.HandleSomeServiceCall(s.someService)).Name("Get something")
	apiAuth.Post("/some", handlers.HandleSomeServiceCall(s.someService)).Name("Post something")
	apiAuth.Put("/some", handlers.HandleSomeServiceCall(s.someService)).Name("Put something")
	apiAuth.Delete("/some", handlers.HandleSomeServiceCall(s.someService)).Name("Delete something")

	return nil
}
