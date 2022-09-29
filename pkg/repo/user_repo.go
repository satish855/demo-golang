package repo

import (
	"context"
	"github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/protos_go/commons/v1"
	"github.com/byteintellect/user_svc/pkg/domain"
	"gorm.io/gorm"
)

type UserRepo interface {
	db.BaseRepository
	CreateRelation(ctx context.Context, primaryUserId string, user *domain.User) (*domain.User, error)
	DeleteRelation(ctx context.Context, userId, relationId string) error
	GetRelations(ctx context.Context, userId string) ([]domain.User, error)
	GetByIdentity(ctx context.Context, identityType int32, identityValue string) (*domain.User, error)
}

type userGORMRepo struct {
	db.BaseRepository
}

func (u *userGORMRepo) CreateRelation(ctx context.Context, primaryUserId string, user *domain.User) (*domain.User, error) {
	if err, base := u.GetByExternalId(ctx, primaryUserId); err != nil {
		return nil, err
	} else {
		externalId := base.GetExternalId()
		user.ParentID = &externalId
		if err, cBase := u.Create(ctx, user); err != nil {
			return nil, err
		} else {
			return cBase.(*domain.User), nil
		}
	}
}

func (u *userGORMRepo) DeleteRelation(ctx context.Context, userId, relationId string) error {
	gDb := u.GetDb().(*gorm.DB)
	var existingRelation *domain.User
	if err := gDb.WithContext(ctx).Table("users").Where("external_id = ? AND parent_id = ?", relationId, userId).Find(&existingRelation).Error; err != nil {
		return err
	}
	existingRelation.Status = int(commonsv1.Status_STATUS_BLOCKED)
	if err, _ := u.Update(ctx, existingRelation.ExternalId, existingRelation); err != nil {
		return err
	}
	return nil
}

func (u *userGORMRepo) GetRelations(ctx context.Context, userId string) ([]domain.User, error) {
	gDb := u.GetDb().(*gorm.DB)
	var relations []domain.User
	if err := gDb.WithContext(ctx).Table("users").Where("parent_id = ?", userId).Find(&relations).Error; err != nil {
		return nil, err
	}
	return relations, nil
}

func (u *userGORMRepo) GetByIdentity(ctx context.Context, identityType int32, identityValue string) (*domain.User, error) {
	gDb := u.GetDb().(*gorm.DB)
	var res domain.User
	if err := gDb.WithContext(ctx).Table("users").Joins("INNER JOIN identities i ON i.user_id=users.external_id").Where("i.identity_type= ? AND i.identity_value =? AND i.status=0",
		identityType, identityValue).Find(&res).Error; err != nil || res.Id == 0 {
		return nil, err
	}
	return &res, nil
}

func NewUserGORMRepo(repo db.BaseRepository) UserRepo {
	return &userGORMRepo{repo}
}
