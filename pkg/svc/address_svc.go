package svc

import (
	"context"
	"errors"
	"github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/protos_go/users/v1"
	"github.com/byteintellect/user_svc/pkg/domain"
	"github.com/byteintellect/user_svc/pkg/repo"
)

type AddressSvc interface {
	GetUserAddresses(ctx context.Context, userId string) ([]*usersv1.AddressDto, error)
	CreateAddress(ctx context.Context, userId string, address *usersv1.AddressDto) (*usersv1.AddressDto, error)
	UpdateAddress(ctx context.Context, userId string, addressId string, address *usersv1.AddressDto) (*usersv1.AddressDto, error)
}

type addressSvcImpl struct {
	repo.AddressRepo
}

func (a *addressSvcImpl) GetUserAddresses(ctx context.Context, userId string) ([]*usersv1.AddressDto, error) {
	if addresses, err := a.AddressRepo.GetUserAddresses(ctx, userId); err != nil {
		return nil, err
	} else {
		var res []*usersv1.AddressDto
		for _, address := range addresses {
			res = append(res, address.ToDto().(*usersv1.AddressDto))
		}
		return res, nil
	}
}

func (a *addressSvcImpl) CreateAddress(ctx context.Context, userId string, dto *usersv1.AddressDto) (*usersv1.AddressDto, error) {
	address := domain.NewAddress(dto)
	address.UserID = userId
	if err, base := a.Create(ctx, address); err != nil {
		return nil, err
	} else {
		return base.(*domain.Address).ToDto().(*usersv1.AddressDto), nil
	}
}

func (a *addressSvcImpl) UpdateAddress(ctx context.Context, userId, addressId string, dto *usersv1.AddressDto) (*usersv1.AddressDto, error) {
	if err, base := a.AddressRepo.GetByExternalId(ctx, addressId); err != nil {
		return nil, err
	} else if eAddress := base.(*domain.Address); eAddress.UserID != userId {
		return nil, errors.New("user id mismatch")
	} else {
		address := domain.NewAddress(dto)
		if err, base := a.Update(ctx, addressId, address); err != nil {
			return nil, err
		} else {
			return base.ToDto().(*usersv1.AddressDto), nil
		}
	}
}

func NewAddressSvc(aRepo db.BaseRepository) AddressSvc {
	return &addressSvcImpl{
		repo.NewAddressGORMRepo(aRepo),
	}
}
