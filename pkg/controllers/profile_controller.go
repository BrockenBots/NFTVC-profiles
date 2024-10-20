package controllers

import (
	"context"
	"fmt"
	"net/http"
	"nftvc-profile/internal/model"
	"nftvc-profile/internal/repository"
	"nftvc-profile/pkg/client"
	"nftvc-profile/pkg/config"
	"nftvc-profile/pkg/logger"
	"nftvc-profile/pkg/requests"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

type ProfileController struct {
	log         logger.Logger
	cfg         *config.Config
	profileRepo repository.ProfileRepository
	authClient  *client.AuthClient
	validate    *validator.Validate
	s3Client    *client.S3Client
}

func NewProfileController(log logger.Logger, cfg *config.Config, profileRepo repository.ProfileRepository, authClient *client.AuthClient, validator *validator.Validate, s3Client *client.S3Client) *ProfileController {
	return &ProfileController{log: log, cfg: cfg, profileRepo: profileRepo, authClient: authClient, validate: validator, s3Client: s3Client}
}

// GetById godoc
// @Summary Get Profile by ID
// @Description Get a profile by its ID
// @Tags profiles
// @Param id path string true "Profile ID"
// @Success 200 {object} model.Profile
// @Failure 400 {object} response.ErrorResponse
// @Router /api/profiles/{id} [get]
func (p *ProfileController) GetById(ctx echo.Context) error {
	p.log.Debugf("GetById")
	id := ctx.Param("id")
	p.log.Debugf("ID: %s", id)
	profile, err := p.profileRepo.GetById(context.Background(), id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Error by get profile by id: %v", err)})
	}

	return ctx.JSON(http.StatusOK, profile)
}

// Save godoc
// @Summary Save Profile
// @Description Save a new profile
// @Tags profiles
// @Accept json
// @Produce json
// @Param profile body requests.SaveProfileRequest true "Profile Data"
// @Success 200 {object} response.SaveProfileResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/profiles/ [post]
func (p *ProfileController) Save(ctx echo.Context) error {
	p.log.Debugf("(SaveProfile)")
	var req requests.SaveProfileRequest
	if err := p.decodeRequest(ctx, &req); err != nil {
		p.log.Debugf("Failed to decode request SaveProfileRequest: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Validation error: %v", err)})
	}

	accountId := ctx.Get("accountId").(string)
	token := ctx.Get("token").(string)

	if exist := p.profileRepo.ProfileExist(context.Background(), accountId); exist {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "profile already exist"})
	}

	if _, err := p.profileRepo.GetByLogin(context.Background(), req.Login); err == nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Login already exist"})
	}

	profileUuid, _ := uuid.NewV7()
	photoPath := ""
	if req.Photo != "" {
		photoBytes, err := p.decodeBase64(req.Photo)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "photo is invalid"})
		}
		photoPath, err = p.s3Client.UploadFile(photoBytes, req.PhotoTitle)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to save photo: %v", err)})
		}
	}

	res, err := p.authClient.ChangeRole(req.Role, token)
	if err != nil {
		p.log.Debugf("Failed to change role: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Failed to change role: %v", err)})
	}

	profile := model.NewProfile(profileUuid.String(), req.Login, req.Name, req.Email, photoPath, req.Description, accountId)
	if err := p.profileRepo.Create(context.Background(), profile); err != nil {
		p.log.Debugf("Failed to create profile")
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to save profile: %v", err)})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"access_token":  res.AccessToken,
		"refresh_token": res.RefreshToken,
	})
}

// GetMe godoc
// @Summary Get current user's profile
// @Description Retrieve the current user's profile
// @Tags profiles
// @Success 200 {object} model.Profile
// @Failure 400 {object} response.ErrorResponse
// @Router /api/profiles/me [get]
func (p *ProfileController) GetMe(ctx echo.Context) error {
	p.log.Debugf("(GetMe)")
	accountId := ctx.Get("accountId").(string)
	profile, err := p.profileRepo.GetByAccountId(context.Background(), accountId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Profile doesnt exist: %v", err)})
	}

	return ctx.JSON(http.StatusOK, profile)
}

// Update godoc
// @Summary Update Profile
// @Description Update an existing profile
// @Tags profiles
// @Accept json
// @Produce json
// @Param profile body requests.UpdateProfileRequest true "Profile Data"
// @Success 200 {object} model.Profile
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/profiles/ [put]
func (p *ProfileController) Update(ctx echo.Context) error {
	p.log.Debugf("(Update)")
	var req requests.UpdateProfileRequest
	if err := p.decodeRequest(ctx, &req); err != nil {
		p.log.Debugf("Failed to decode request UpdateProfileRequest: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Validation error: %v", err)})
	}

	accountId := ctx.Get("accountId").(string)

	currentProfile, err := p.profileRepo.GetByAccountId(context.Background(), accountId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Profile doesn't exist: %v", err)})
	}

	if req.Login != "" && req.Login != currentProfile.Login {
		if _, err := p.profileRepo.GetByLogin(context.Background(), req.Login); err == nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Login already exist"})
		}
		currentProfile.Login = req.Login
	}

	if req.Name != "" && req.Name != currentProfile.Name {
		currentProfile.Name = req.Name
	}

	if req.Email != "" && req.Email != currentProfile.Email {
		currentProfile.Email = req.Email
	}

	if req.Description != "" && req.Description != currentProfile.Description {
		currentProfile.Description = req.Description
	}

	if req.Photo != "" && req.PhotoTitle != "" {
		photoBytes, err := p.decodeBase64(req.Photo)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "photo is invalid"})
		}
		photoPath, err := p.s3Client.UploadFile(photoBytes, req.PhotoTitle)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to update photo: %v", err)})
		}
		currentProfile.Photo = photoPath
	}

	if err := p.profileRepo.Update(context.Background(), currentProfile); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to update profile: %v", err)})
	}

	return ctx.JSON(http.StatusOK, currentProfile)
}

func (p *ProfileController) GetByWalletAddress(ctx echo.Context) error {
	var req requests.GetByWalletAddressRequest
	if err := p.decodeRequest(ctx, &req); err != nil {
		p.log.Debugf("Failed to decode request UpdateProfileRequest: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Validation error: %v", err)})
	}

	profile, err := p.profileRepo.GetByWalletAddress(context.Background(), req.WalletAddress)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Error by get profile by id: %v", err)})
	}

	return ctx.JSON(http.StatusOK, profile)
}
