package repo

import (
	"context"
	"github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/user_svc/pkg/domain"
	"gorm.io/gorm"
)

type IdentityRepo interface {
	db.BaseRepository
	GetExistingIdentity(ctx context.Context, identity *domain.Identity) (*domain.Identity, error)
	GetIdentitiesForUser(ctx context.Context, userId string) ([]domain.Identity, error)
}

type identityGormRepo struct {
	db.BaseRepository
}

func (i *identityGormRepo) GetExistingIdentity(ctx context.Context, identity *domain.Identity) (*domain.Identity, error) {
	gDb := i.GetDb().(*gorm.DB)
	var existingIdentity domain.Identity
	if err := gDb.WithContext(ctx).Table("identities").Where("identity_type = ? AND identity_value = ?", int32(identity.IdentityType.Number()), identity.IdentityValue).Find(&existingIdentity).Error; err != nil {
		return nil, err
	}
	return &existingIdentity, nil
}

func (i *identityGormRepo) GetIdentitiesForUser(ctx context.Context, userId string) ([]domain.Identity, error) {
	gDb := i.GetDb().(*gorm.DB)
	var res []domain.Identity
	if err := gDb.WithContext(ctx).Table("identities").Where("user_id = ?", userId).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func NewIdentityRepo(repo db.BaseRepository) IdentityRepo {
	return &identityGormRepo{
		repo,
	}
}
