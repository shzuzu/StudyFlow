package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"payment_service/internal/domain"
	"time"
)

type PaymentRepo struct {
	db *sql.DB
}

func (r *PaymentRepo) CreateReceipt(ctx context.Context, receipt *domain.PaymentReceipt) error {
	query := `
        INSERT INTO receipts (id, lesson_id, file_id, is_verified)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.ExecContext(
		ctx,
		query,
		receipt.ID,
		receipt.LessonID,
		receipt.FileID,
		receipt.IsVerified,
		receipt.CreatedAt,
		receipt.EditedAt)
	return err
}

func (r *PaymentRepo) GetReceipt(ctx context.Context, id uuid.UUID) (*domain.PaymentReceipt, error) {
	query := `
        SELECT id, lesson_id, file_id, is_verified, created_at, edited_at
        FROM receipts
        WHERE id = $1
    `

	var receipt domain.PaymentReceipt
	var editedAt *time.Time
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&receipt.ID,
		&receipt.LessonID,
		&receipt.FileID,
		&receipt.IsVerified,
		&receipt.CreatedAt,
		&editedAt,
	)
	receipt.EditedAt = editedAt
	return &receipt, err
}

func (r *PaymentRepo) VerifyReceipt(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE receipts
		SET is_verified = TRUE, edited_at =
		CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PaymentRepo) DeleteReceipt(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM receipts
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PaymentRepo) GetReceiptByLessonID(ctx context.Context, lessonID uuid.UUID) (*domain.PaymentReceipt, error) {
	query := `
        SELECT id, lesson_id, file_id, is_verified, created_at, edited_at
        FROM receipts
        WHERE lesson_id = $1
    `
	var receipt domain.PaymentReceipt
	var editedAt *time.Time
	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(
		&receipt.ID, &receipt.LessonID, &receipt.FileID, &receipt.IsVerified, &receipt.CreatedAt, &editedAt,
	)
	if err != nil {
		return nil, err
	}
	receipt.EditedAt = editedAt
	return &receipt, nil
}
