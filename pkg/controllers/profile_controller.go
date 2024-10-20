package controllers

import (
	"fmt"
	"nftvc-profile/internal/service"
	"nftvc-profile/pkg/config"
	"nftvc-profile/pkg/logger"

	"github.com/labstack/echo/v4"
)

type ProfileController struct {
	log            logger.Logger
	cfg            *config.Config
	profileService *service.ProfileService
}

func NewProfileController(log logger.Logger, cfg *config.Config, profileService *service.ProfileService) *ProfileController {
	return &ProfileController{log: log, cfg: cfg, profileService: profileService}
}

func (p *ProfileController) GetById(ctx echo.Context) error {
	return fmt.Errorf("not impl")
}

func (p *ProfileController) Save(ctx echo.Context) error {
	return fmt.Errorf("not impl")
}

func (p *ProfileController) GetMe(ctx echo.Context) error {
	return fmt.Errorf("not impl")
}

func (p *ProfileController) Update(ctx echo.Context) error {
	return fmt.Errorf("not impl")
}

func (p *ProfileController) GetByWalletAddress(ctx echo.Context) error {
	return fmt.Errorf("not impl")
}
