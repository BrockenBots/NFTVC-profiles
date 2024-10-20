package repo

import (
	"context"
	"fmt"
	"nftvc-profile/internal/model"
	"nftvc-profile/pkg/config"
	"nftvc-profile/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileMongoRepo struct {
	log logger.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewProfileRepo(log logger.Logger, cfg *config.Config, db *mongo.Client) *ProfileMongoRepo {
	return &ProfileMongoRepo{log: log, cfg: cfg, db: db}
}

func (p *ProfileMongoRepo) Create(ctx context.Context, profile *model.Profile) error {
	return fmt.Errorf("not impl")
}
func (p *ProfileMongoRepo) Update(ctx context.Context, profile *model.Profile) error {
	return fmt.Errorf("not impl")
}
func (p *ProfileMongoRepo) GetById(ctx context.Context, profileId string) (*model.Profile, error) {
	return nil, fmt.Errorf("not impl")
}
func (p *ProfileMongoRepo) GetByWalletAddress(ctx context.Context, walletAddress string) (*model.Profile, error) {
	return nil, fmt.Errorf("not impl")
}
func (p *ProfileMongoRepo) Delete(ctx context.Context, profileId string) error {
	return fmt.Errorf("not impl")
}

func (p *ProfileMongoRepo) getProfilesCollections() *mongo.Collection {
	return p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Profiles)
}
