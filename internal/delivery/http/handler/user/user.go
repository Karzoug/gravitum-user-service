package user

import (
	"context"
	"net/http"

	"github.com/rs/xid"

	gen "github.com/Karzoug/gravitum-user-service/internal/delivery/http/gen/user/v1"
	"github.com/Karzoug/gravitum-user-service/internal/delivery/http/httperr"
	"github.com/Karzoug/gravitum-user-service/internal/entity"
)

//nolint:revive,stylecheck // codegen name
func (h handlers) GetUsersId(ctx context.Context, request gen.GetUsersIdRequestObject) (gen.GetUsersIdResponseObject, error) {
	id, err := xid.FromString(request.Id)
	if err != nil {
		return nil, httperr.NewError("invalid id", http.StatusBadRequest)
	}

	u, err := h.userService.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return gen.GetUsersId200JSONResponse(toGenUser(u)), nil
}

func (h handlers) PostUsers(ctx context.Context, request gen.PostUsersRequestObject) (gen.PostUsersResponseObject, error) {
	if request.Body == nil {
		return nil, httperr.NewError("request body is nil", http.StatusBadRequest)
	}

	u := froGenMutableUser(*request.Body)
	id, err := h.userService.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	u.ID = id

	return gen.PostUsers201JSONResponse(toGenUser(u)), nil
}

//nolint:revive,stylecheck // codegen name
func (h handlers) PutUsersId(ctx context.Context, request gen.PutUsersIdRequestObject) (gen.PutUsersIdResponseObject, error) {
	if request.Body == nil {
		return nil, httperr.NewError("request body is nil", http.StatusBadRequest)
	}
	id, err := xid.FromString(request.Id)
	if err != nil {
		return nil, httperr.NewError("invalid id", http.StatusBadRequest)
	}

	u := froGenMutableUser(*request.Body)
	u.ID = id

	if err := h.userService.Update(ctx, u); err != nil {
		return nil, err
	}

	return gen.PutUsersId200JSONResponse(toGenUser(u)), nil
}

//nolint:revive,stylecheck // codegen name
func (h handlers) DeleteUsersId(ctx context.Context, request gen.DeleteUsersIdRequestObject) (gen.DeleteUsersIdResponseObject, error) {
	id, err := xid.FromString(request.Id)
	if err != nil {
		return nil, httperr.NewError("invalid id", http.StatusBadRequest)
	}

	if err := h.userService.Delete(ctx, id); err != nil {
		return nil, err
	}

	return gen.DeleteUsersId204Response{}, nil
}

func toGenUser(u entity.User) gen.User {
	return gen.User{
		Username:   u.Username,
		Name:       u.Name,
		ImageUrl:   u.ImageURL,
		StatusText: u.StatusText,
		Id:         u.ID.String(),
	}
}

func froGenMutableUser(u gen.MutableUser) entity.User {
	return entity.User{
		Username:   u.Username,
		Name:       u.Name,
		ImageURL:   u.ImageUrl,
		StatusText: u.StatusText,
	}
}
