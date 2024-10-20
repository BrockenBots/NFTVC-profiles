package server

import (
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *server) runHttpServer() error {
	s.echo.Use(s.middleware.CORS())

	s.mapRoutes()

	return s.echo.Start(s.cfg.Http.Port)
}

func (s *server) mapRoutes() {
	s.echo.GET("api/profiles/:id", s.profileController.GetById)

	authGroup := s.echo.Group("api/profiles", s.middleware.AuthMiddleware)

	authGroup.POST("/", s.profileController.Save)
	authGroup.GET("/me", s.profileController.GetMe)
	authGroup.PATCH("/", s.profileController.Update)

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
