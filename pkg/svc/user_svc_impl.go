package svc

import (
	"context"
	"errors"
	"github.com/byteintellect/go_commons/cache"
	"github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/protos_go/users/v1"
	"github.com/byteintellect/user_svc/pkg/domain"
	"github.com/byteintellect/user_svc/pkg/repo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserSvcImpl struct {
	usersv1.UnimplementedUserServiceServer
	repo.UserRepo
	identitySvc IdentitySvc
	addressSvc  AddressSvc
	cache       cache.BaseCache
	lgr         *zap.Logger
}

func NewUserServiceServer(userRepo db.BaseRepository, identitySvc IdentitySvc, addressSvc AddressSvc, cache cache.BaseCache, lgr *zap.Logger) usersv1.UserServiceServer {
	return &UserSvcImpl{
		UserRepo:    repo.NewUserGORMRepo(userRepo),
		identitySvc: identitySvc,
		addressSvc:  addressSvc,
		cache:       cache,
		lgr:         lgr,
	}
}

func (u *UserSvcImpl) CreateUser(ctx context.Context, request *usersv1.CreateUserRequest) (*usersv1.CreateUserResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		user := domain.NewUser(request.Request)
		if err, user := u.Create(ctx, user); err != nil {
			grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "400"))
			return nil, status.Errorf(codes.InvalidArgument, "user %v", err)
		} else {
			grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "201"))
			return &usersv1.CreateUserResponse{Response: user.(*domain.User).ToDto().(*usersv1.UserDto)}, nil
		}
	}
}

func (u *UserSvcImpl) UpdateUser(ctx context.Context, request *usersv1.UpdateUserRequest) (*usersv1.UpdateUserResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		user := domain.NewUser(request.Request)
		if err, user := u.Update(ctx, request.UserId, user); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "user %v", err)
		} else {
			return &usersv1.UpdateUserResponse{Response: user.(*domain.User).ToDto().(*usersv1.UserDto)}, nil
		}
	}
}

func (u *UserSvcImpl) GetUser(ctx context.Context, request *usersv1.GetUserByIdRequest) (*usersv1.GetUserByIdResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		base, err := u.cache.Get(ctx, request.UserId)
		if err == nil && base.GetExternalId() != "" {
			return &usersv1.GetUserByIdResponse{
				Response: base.(*domain.User).ToDto().(*usersv1.UserDto),
			}, nil
		}
		if err, user := u.UserRepo.GetByExternalId(ctx, request.UserId); err != nil {
			grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "404"))
			return nil, status.Errorf(codes.NotFound, "user %v", err)
		} else {
			grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "200"))
			err := u.cache.Put(ctx, user)
			if err != nil {
				u.lgr.Error("user record cache insertion error", zap.Error(err))
			}
			return &usersv1.GetUserByIdResponse{Response: user.(*domain.User).ToDto().(*usersv1.UserDto)}, nil
		}
	}
}

func (u *UserSvcImpl) BlockUser(ctx context.Context, request *usersv1.BlockUserRequest) (*usersv1.BlockUserResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		if err, base := u.UserRepo.GetByExternalId(ctx, request.UserId); err != nil {
			return nil, err
		} else {
			user := base.(*domain.User)
			user.Block()
			if err, _ := u.Update(ctx, request.UserId, user); err != nil {
				return nil, err
			}
			return &usersv1.BlockUserResponse{
				Status: true,
			}, nil
		}
	}
}

func (u *UserSvcImpl) CreateUserIdentity(ctx context.Context, request *usersv1.CreateUserIdentityRequest) (*usersv1.CreateUserIdentityResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		if cIdentity, err := u.identitySvc.CreateForUser(ctx, request.UserId, request.Request); err != nil {
			return nil, err
		} else {
			return &usersv1.CreateUserIdentityResponse{
				Response: cIdentity,
			}, nil
		}
	}
}

func (u *UserSvcImpl) UpdateUserIdentity(ctx context.Context, request *usersv1.UpdateUserIdentityRequest) (*usersv1.UpdateUserIdentityResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		res, err := u.identitySvc.UpdateIdentity(ctx, request.UserId, request.IdentityId, request.Request)
		if err != nil {
			return nil, err
		}
		return &usersv1.UpdateUserIdentityResponse{
			Response: res,
		}, nil
	}
}

func (u *UserSvcImpl) GetUserIdentities(ctx context.Context, request *usersv1.GetUserIdentitiesRequest) (*usersv1.GetUserIdentitiesResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		res, err := u.identitySvc.GetIdentitiesForUser(ctx, request.UserId)
		if err != nil {
			return nil, err
		}
		return &usersv1.GetUserIdentitiesResponse{
			Response: res,
		}, nil
	}
}

func (u *UserSvcImpl) CreateUserRelation(ctx context.Context, request *usersv1.CreateUserRelationRequest) (*usersv1.CreateUserRelationResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		user := domain.NewUser(request.Request)
		if res, err := u.UserRepo.CreateRelation(ctx, request.PrimaryUserId, user); err != nil {
			return nil, err
		} else {
			return &usersv1.CreateUserRelationResponse{Response: res.ToDto().(*usersv1.UserDto)}, nil
		}
	}
}

func (u *UserSvcImpl) DeleteUserRelation(ctx context.Context, request *usersv1.DeleteUserRelationRequest) (*usersv1.DeleteUserRelationResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		if err := u.UserRepo.DeleteRelation(ctx, request.PrimaryUserId, request.RelationId); err != nil {
			return nil, err
		}
		return &usersv1.DeleteUserRelationResponse{
			Status: true,
		}, nil
	}
}

func (u *UserSvcImpl) CreateUserAddress(ctx context.Context, request *usersv1.CreateUserAddressRequest) (*usersv1.CreateUserAddressResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		if address, err := u.addressSvc.CreateAddress(ctx, request.UserId, request.Request); err != nil || address == nil {
			return nil, err
		} else {
			return &usersv1.CreateUserAddressResponse{
				Response: address,
			}, nil
		}
	}
}

func (u *UserSvcImpl) UpdateUserAddress(ctx context.Context, request *usersv1.UpdateUserAddressRequest) (*usersv1.UpdateUserAddressResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		if updatedAddress, err := u.addressSvc.UpdateAddress(ctx, request.UserId, request.AddressId, request.Request); err != nil || updatedAddress == nil {
			return nil, err
		} else {
			return &usersv1.UpdateUserAddressResponse{Response: updatedAddress}, nil
		}
	}
}

func (u *UserSvcImpl) GetUserAddresses(ctx context.Context, request *usersv1.GetUserAddressesRequest) (*usersv1.GetUserAddressesResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		if addresses, err := u.addressSvc.GetUserAddresses(ctx, request.UserId); err != nil || len(addresses) == 0 {
			return nil, err
		} else {
			return &usersv1.GetUserAddressesResponse{
				Response: addresses,
			}, nil
		}
	}
}

func (u *UserSvcImpl) GetUserByIdentity(ctx context.Context, request *usersv1.GetUserByIdentityRequest) (*usersv1.CreateUserResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		if user, err := u.UserRepo.GetByIdentity(ctx, int32(request.Type), request.Value); err != nil || user == nil {
			return nil, err
		} else {
			return &usersv1.CreateUserResponse{
				Response: user.ToDto().(*usersv1.UserDto),
			}, nil
		}
	}
}
