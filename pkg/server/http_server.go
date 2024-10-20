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
	authGroup := s.echo.Group("api/profiles/", s.middleware.AuthMiddleware)

	authGroup.POST("/", s.profileController.Save)
	authGroup.GET("/", s.profileController.GetMe)
	authGroup.PUT("/", s.profileController.Update)

	s.echo.GET("/{id}", s.profileController.GetById)

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
