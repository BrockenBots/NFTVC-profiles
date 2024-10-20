package repo

import (
	"context"
	"nftvc-profile/internal/model"
	"nftvc-profile/pkg/config"
	"nftvc-profile/pkg/logger"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	_, err := p.getProfilesCollections().InsertOne(ctx, profile, &options.InsertOneOptions{})
	if err != nil && !strings.Contains(err.Error(), "no documents") {
		p.log.Debugf("(ProfileMongoRepo) error: %v", err)
		return err
	}

	return nil
}
func (p *ProfileMongoRepo) Update(ctx context.Context, profile *model.Profile) error {
	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(false)

	err := p.getProfilesCollections().FindOneAndUpdate(ctx, bson.M{"_id": profile.Id}, bson.M{"$set": profile}, ops).Err()

	if err != nil {
		p.log.Debugf("(ProfileMongoRepo) error: %v", err)
		return err
	}
	return nil
}
func (p *ProfileMongoRepo) GetById(ctx context.Context, profileId string) (*model.Profile, error) {
	var profile model.Profile
	err := p.getProfilesCollections().FindOne(ctx, bson.M{"_id": profileId}).Decode(&profile)
	if err != nil {
		p.log.Debugf("(ProfileMongoRepo) error: %v", err)
		return nil, err
	}
	return &profile, nil
}
func (p *ProfileMongoRepo) GetByWalletAddress(ctx context.Context, walletAddress string) (*model.Profile, error) {
	var profile model.Profile
	err := p.getProfilesCollections().FindOne(ctx, bson.M{"walletAddress": walletAddress}).Decode(&profile)
	if err != nil {
		p.log.Debugf("(ProfileMongoRepo) error: %v", err)
		return nil, err
	}
	return &profile, nil
}
func (p *ProfileMongoRepo) Delete(ctx context.Context, profileId string) error {
	_, err := p.getProfilesCollections().DeleteOne(ctx, bson.M{"_id": profileId})
	if err != nil {
		p.log.Debugf("(ProfileMongoRepo) error: %v", err)
		return err
	}
	return nil
}

func (p *ProfileMongoRepo) GetByAccountId(ctx context.Context, accountId string) (*model.Profile, error) {
	var profile model.Profile
	err := p.getProfilesCollections().FindOne(ctx, bson.M{"accountId": accountId}).Decode(&profile)
	if err != nil {
		p.log.Debugf("(ProfileMongoRepo) error: %v", err)
		return nil, err
	}

	return &profile, nil
}

func (p *ProfileMongoRepo) GetByLogin(ctx context.Context, login string) (*model.Profile, error) {
	var profile model.Profile
	err := p.getProfilesCollections().FindOne(ctx, bson.M{"login": login}).Decode(&profile)
	if err != nil {
		p.log.Debugf("(ProfileMongoRepo) error: %v", err)
		return nil, err
	}

	return &profile, nil
}

func (p *ProfileMongoRepo) ProfileExist(ctx context.Context, accountId string) bool {
	var profile model.Profile
	err := p.getProfilesCollections().FindOne(ctx, bson.M{"accountId": accountId}).Decode(&profile)
	if err != nil {
		p.log.Debugf("(ProfileMongoRepo) error: %v", err)
		return false
	}

	return true
}

func (p *ProfileMongoRepo) getProfilesCollections() *mongo.Collection {
	return p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Profiles)
}
