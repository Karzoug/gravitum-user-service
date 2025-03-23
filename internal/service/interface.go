package service

import (
	"context"

	"github.com/rs/xid"

	"github.com/Karzoug/gravitum-user-service/internal/entity"
)

type repository interface {
	Create(ctx context.Context, user entity.User) error
	Get(ctx context.Context, id xid.ID) (entity.User, error)
	Update(ctx context.Context, u entity.User) error
	Delete(ctx context.Context, id xid.ID) error
}
