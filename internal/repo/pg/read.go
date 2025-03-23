package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"

	"github.com/Karzoug/gravitum-user-service/internal/entity"
	repoerr "github.com/Karzoug/gravitum-user-service/internal/repo"
)

func (r repo) Get(ctx context.Context, id xid.ID) (entity.User, error) {
	const (
		op    = "postgresql: gen one user"
		query = `
SELECT username, name, image_url, status_text, updated_at
FROM users
WHERE id = @id`
	)

	row, err := r.db.Query(ctx, query,
		pgx.NamedArgs{
			"id": id,
		})
	if err != nil {
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	u, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[entity.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerr.ErrRecordNotFound
		}
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	u.ID = id

	return u, nil
}
