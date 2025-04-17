package data

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"payment_service/internal/errors"
)

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func handleError(err error) error {
	if isUniqueViolation(err) {
		return errors1.ErrAlreadyExists
	}
	if isNotFound(err) {
		return errors1.ErrReceiptNotFound
	}
	return fmt.Errorf("repository error: %w", err)
}
