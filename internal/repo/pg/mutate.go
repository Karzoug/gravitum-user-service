package pg

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/xid"

	"github.com/Karzoug/gravitum-user-service/internal/entity"
	repoerr "github.com/Karzoug/gravitum-user-service/internal/repo"
)

func (r repo) Create(ctx context.Context, user entity.User) error {
	const (
		op          = "postgresql: create user"
		queryCreate = `
INSERT INTO users (id, username, name, image_url, status_text)
VALUES (@id, @username, @name, @image_url, @status_text)`
	)

	tag, err := r.db.Exec(ctx, queryCreate,
		pgx.NamedArgs{
			"id":          user.ID,
			"username":    user.Username,
			"name":        user.Name,
			"image_url":   user.ImageURL,
			"status_text": user.StatusText,
		})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if strings.HasPrefix(pgErr.Code, "23") && pgErr.TableName == "users" {
				return repoerr.ErrRecordAlreadyExists
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repoerr.ErrNoAffected
	}

	return nil
}

func (r repo) Update(ctx context.Context, user entity.User) error {
	const (
		op    = "postgresql: update user"
		query = `
UPDATE users
SET name = @name, username = @username, image_url = @image_url, status_text = @status_text
WHERE id = @id`
	)

	tag, err := r.db.Exec(ctx, query,
		pgx.NamedArgs{
			"id":          user.ID,
			"name":        user.Name,
			"username":    user.Username,
			"image_url":   user.ImageURL,
			"status_text": user.StatusText,
		})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return repoerr.ErrRecordNotFound
	}

	return nil
}

func (r repo) Delete(ctx context.Context, id xid.ID) error {
	const (
		op    = "postgresql: delete user"
		query = `DELETE FROM users WHERE id = @id`
	)

	_, err := r.db.Exec(ctx, query,
		pgx.NamedArgs{
			"id": id,
		})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
