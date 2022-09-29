package svc

import (
	"context"
	"errors"
	"github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/protos_go/users/v1"
	"github.com/byteintellect/user_svc/pkg/domain"
	"github.com/byteintellect/user_svc/pkg/repo"
)

type IdentitySvc interface {
	CreateForUser(ctx context.Context, userId string, identity *usersv1.IdentityDto) (*usersv1.IdentityDto, error)
	GetIdentitiesForUser(ctx context.Context, userId string) ([]*usersv1.IdentityDto, error)
	UpdateIdentity(ctx context.Context, userId, identityId string, identity *usersv1.IdentityDto) (*usersv1.IdentityDto, error)
}

type identitySvc struct {
	repo.IdentityRepo
}

func (i *identitySvc) CreateForUser(ctx context.Context, userId string, dto *usersv1.IdentityDto) (*usersv1.IdentityDto, error) {
	identity := domain.NewIdentity(dto)
	identity.UserID = userId
	if base, err := i.IdentityRepo.GetExistingIdentity(ctx, identity); err == nil && base != nil && base.Id != 0 {
		return nil, errors.New("identity already exists")
	}
	if err, base := i.IdentityRepo.Create(ctx, identity); err != nil {
		return nil, err
	} else {
		return base.(*domain.Identity).ToDto().(*usersv1.IdentityDto), nil
	}
}

func (i *identitySvc) GetIdentitiesForUser(ctx context.Context, userId string) ([]*usersv1.IdentityDto, error) {
	if identities, err := i.IdentityRepo.GetIdentitiesForUser(ctx, userId); err != nil {
		return nil, err
	} else {
		var res []*usersv1.IdentityDto
		for _, identity := range identities {
			res = append(res, identity.ToDto().(*usersv1.IdentityDto))
		}
		return res, nil
	}
}

func (i *identitySvc) UpdateIdentity(ctx context.Context, userId, identityId string, dto *usersv1.IdentityDto) (*usersv1.IdentityDto, error) {
	identity := domain.NewIdentity(dto)
	if err, base := i.IdentityRepo.GetByExternalId(ctx, identityId); err != nil {
		return nil, err
	} else if eIdentity := base.(*domain.Identity); eIdentity.UserID != userId {
		return nil, errors.New("user id mismatch")
	} else {
		if err, base := i.IdentityRepo.Update(ctx, identityId, identity); err != nil {
			return nil, err
		} else {
			return base.(*domain.Identity).ToDto().(*usersv1.IdentityDto), nil
		}
	}
}

func NewIdentitySvc(identityRepo db.BaseRepository) IdentitySvc {
	return &identitySvc{
		repo.NewIdentityRepo(identityRepo),
	}
}
