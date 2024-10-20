package service

import (
	"nftvc-profile/internal/repository"
	"nftvc-profile/pkg/client"
	"nftvc-profile/pkg/logger"
)

type ProfileService struct {
	log         logger.Logger
	profileRepo repository.ProfileRepository
	authClient  *client.AuthClient
}

func NewProfileService(log logger.Logger, profileRepo repository.ProfileRepository, authClient *client.AuthClient) *ProfileService {
	return &ProfileService{log: log, profileRepo: profileRepo, authClient: authClient}
}

// func (p *ProfileService)
