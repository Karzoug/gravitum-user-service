package service

import (
	"context"
	"errors"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"

	"github.com/Karzoug/gravitum-user-service/pkg/ucerr"

	"github.com/Karzoug/gravitum-user-service/internal/entity"
	repoerr "github.com/Karzoug/gravitum-user-service/internal/repo"
)

type UserService struct {
	repo   repository
	logger zerolog.Logger
}

// NewUserService creates a new user service.
func NewUserService(repo repository, logger zerolog.Logger) UserService {
	logger = logger.With().
		Str("component", "user service").
		Logger()

	return UserService{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new user.
func (us UserService) Create(ctx context.Context, u entity.User) (xid.ID, error) {
	if err := u.Validate(); err != nil {
		return xid.NilID(), ucerr.NewError(err, err.Error(), codes.InvalidArgument)
	}

	u.ID = xid.New()

	if err := us.repo.Create(ctx, u); err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordAlreadyExists):
			return xid.NilID(), ucerr.NewError(err, "user already exists", codes.AlreadyExists)
		default:
			return xid.NilID(), ucerr.NewInternalError(err)
		}
	}

	return u.ID, nil
}

// Update updates an existing user.
func (us UserService) Update(ctx context.Context, u entity.User) error {
	if err := u.Validate(); err != nil {
		return ucerr.NewError(err, err.Error(), codes.InvalidArgument)
	}

	if err := us.repo.Update(ctx, u); err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return ucerr.NewError(err, "user not found", codes.NotFound)
		default:
			return ucerr.NewInternalError(err)
		}
	}

	return nil
}

// Get returns an existing user.
func (us UserService) Get(ctx context.Context, id xid.ID) (entity.User, error) {
	u, err := us.repo.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoerr.ErrRecordNotFound):
			return entity.User{}, ucerr.NewError(err, "user not found", codes.NotFound)
		default:
			return entity.User{}, ucerr.NewInternalError(err)
		}
	}

	return u, nil
}

// Delete deletes an existing user.
func (us UserService) Delete(ctx context.Context, id xid.ID) error {
	if err := us.repo.Delete(ctx, id); err != nil {
		return ucerr.NewInternalError(err)
	}

	return nil
}
