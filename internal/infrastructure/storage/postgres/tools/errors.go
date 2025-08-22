package tools

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsUniqueErr(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
