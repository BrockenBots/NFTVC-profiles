package repository

import (
	"context"
	"nftvc-profile/internal/model"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *model.Profile) error
	Update(ctx context.Context, profile *model.Profile) error
	GetById(ctx context.Context, profileId string) (*model.Profile, error)
	GetByWalletAddress(ctx context.Context, walletAddress string) (*model.Profile, error)
	GetByAccountId(ctx context.Context, accountId string) (*model.Profile, error)
	GetByLogin(ctx context.Context, login string) (*model.Profile, error)
	ProfileExist(ctx context.Context, accountId string) bool
	Delete(ctx context.Context, profileId string) error
}
