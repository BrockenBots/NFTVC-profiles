package server

import (
	"context"
	"nftvc-profile/internal/service"
	"nftvc-profile/pkg/client"
	"nftvc-profile/pkg/config"
	"nftvc-profile/pkg/controllers"
	"nftvc-profile/pkg/logger"
	"nftvc-profile/pkg/middlewares"
	"nftvc-profile/pkg/mongodb"
	"nftvc-profile/pkg/repo"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	log               logger.Logger
	cfg               *config.Config
	echo              *echo.Echo
	profileController *controllers.ProfileController
	mongoClient       *mongo.Client
	middleware        *middlewares.MiddlewareManager
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{
		log:  log,
		cfg:  cfg,
		echo: echo.New(),
	}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	mongo, err := mongodb.NewMongoDbConn(ctx, s.cfg.Mongo)
	if err != nil {
		return err
	}
	s.mongoClient = mongo
	s.initMongoDBCollections(ctx)

	authClient := client.NewAuthClient(s.log, s.cfg)

	middlewareManager := middlewares.NewMiddlewareManager(s.log, s.cfg, authClient)
	s.middleware = middlewareManager
	profileRepo := repo.NewProfileRepo(s.log, s.cfg, s.mongoClient)
	profileService := service.NewProfileService(s.log, profileRepo, authClient)
	profileController := controllers.NewProfileController(s.log, s.cfg, profileService)

	s.profileController = profileController

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.log.Error("(HttpServer) err: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdownCtx()

	s.log.Infof("Server shutdown...")

	if err := s.echo.Shutdown(shutdownCtx); err != nil {
		s.log.Infof("Shutdown server with error: %v", err)
		return err
	}

	s.log.Infof("Server shutdown succesfuly")

	return nil
}
