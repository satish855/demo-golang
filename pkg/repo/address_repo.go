package repo

import (
	"context"
	"github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/user_svc/pkg/domain"
	"gorm.io/gorm"
)

type AddressRepo interface {
	db.BaseRepository
	GetUserAddresses(ctx context.Context, userId string) ([]domain.Address, error)
}

type addressGORMRepo struct {
	db.BaseRepository
}

func (a *addressGORMRepo) GetUserAddresses(ctx context.Context, userId string) ([]domain.Address, error) {
	gDb := a.GetDb().(*gorm.DB)
	var res []domain.Address
	if err := gDb.WithContext(ctx).Table("addresses").Where("user_id = ?", userId).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func NewAddressGORMRepo(repo db.BaseRepository) AddressRepo {
	return &addressGORMRepo{
		repo,
	}
}
