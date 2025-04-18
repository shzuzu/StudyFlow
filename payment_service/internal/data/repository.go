package data

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	errors1 "payment_service/internal/errors"
	"payment_service/internal/models"
	"time"
)

type PaymentRepo struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) CreateReceipt(ctx context.Context, receipt *models.PaymentReceiptCreateInput) (*models.PaymentReceipt, error) {

	query := `
		INSERT INTO receipts (id, lesson_id, file_id, is_verified, created_at, edited_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, lesson_id, file_id, is_verified, created_at, edited_at
	`
	now := time.Now()
	pr := &models.PaymentReceipt{}
	err := pgxscan.Get(ctx, r.db, pr, query,
		receipt.ID,
		receipt.LessonID,
		receipt.FileID,
		false,
		now,
		now,
	)
	if err != nil {
		return nil, handleError(err)
	}
	return pr, nil
}

func (r *PaymentRepo) GetReceiptByID(ctx context.Context, id uuid.UUID) (*models.PaymentReceipt, error) {
	query := `
		SELECT id, lesson_id, file_id, is_verified, created_at, edited_at FROM receipts WHERE id = $1
	`

	pr := &models.PaymentReceipt{}
	err := pgxscan.Get(ctx, r.db, pr, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors1.ErrReceiptNotFound
		}
		return nil, handleError(err)
	}

	return pr, nil
}

func (r *PaymentRepo) UpdateReceipt(ctx context.Context, receipt *models.PaymentReceiptUpdateInput) error {
	query := `
		UPDATE receipts 
		SET is_verified = $1, edited_at = $2 
		WHERE id = $3
	`
	now := time.Now()

	_, err := r.db.Exec(ctx, query,
		receipt.IsVerified,
		now,
		receipt.ID,
	)
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (r *PaymentRepo) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM receipts WHERE id = $1)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, handleError(err)
	}
	return exists, nil
}

func (r *PaymentRepo) GetReceiptByLessonID(ctx context.Context, lessonID uuid.UUID) (*models.PaymentReceipt, error) {
	query := `
		SELECT id, lesson_id, file_id, is_verified, created_at, edited_at FROM receipts WHERE lesson_id = $1
	`

	pr := &models.PaymentReceipt{}
	err := pgxscan.Get(ctx, r.db, pr, query, lessonID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors1.ErrReceiptNotFound
		}
		return nil, handleError(err)
	}
	return pr, nil
}
