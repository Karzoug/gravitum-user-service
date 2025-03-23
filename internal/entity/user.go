package entity

import (
	"time"

	"github.com/rs/xid"
)

type User struct {
	ID         xid.ID    `db:"id"`
	Username   string    `db:"username" validate:"required,min=3,max=50"`
	Name       string    `db:"name" validate:"required,min=1,max=50"`
	ImageURL   string    `db:"image_url" validate:"omitempty,url,max=255"`
	StatusText string    `db:"status_text" validate:"omitempty,max=200"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (u User) Validate() error {
	return validatorError(validate.Struct(u))
}
