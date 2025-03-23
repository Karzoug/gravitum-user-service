package pg

import (
	"github.com/Karzoug/gravitum-user-service/pkg/postgresql"
)

type repo struct {
	db postgresql.DB
}

func NewUserRepo(db postgresql.DB) repo {
	return repo{db: db}
}
