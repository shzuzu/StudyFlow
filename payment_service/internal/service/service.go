//go:generate mockgen -source=service.go -destination=../mocks/payment_mocks.go -package=mocks

package service

import (
	"context"
	api2 "fileservice/pkg/api"
	"github.com/google/uuid"
	"paymentservice/internal/clients"
	errdefs "paymentservice/internal/errors"
	"paymentservice/internal/models"
	api3 "schedule_service/pkg/api"
	"time"
)

const maxRetries = 5                      // Максимальное количество попыток
const retryDelay = 500 * time.Millisecond // Задержка между попытками

type IPaymentRepo interface {
	CreateReceipt(ctx context.Context, receipt *models.PaymentReceiptCreateInput) (*models.PaymentReceipt, error)

	GetReceiptByID(ctx context.Context, id uuid.UUID) (*models.PaymentReceipt, error)

	UpdateReceipt(ctx context.Context, id uuid.UUID, isVerified bool) (*models.PaymentReceipt, error)

	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)

	GetReceiptByLessonID(ctx context.Context, lessonID uuid.UUID) (*models.PaymentReceipt, error)
}

type PaymentService struct {
	repo           IPaymentRepo
	userClient     clients.UserServiceClient
	fileClient     clients.FileServiceClient
	scheduleClient clients.ScheduleServiceClient
}

func NewPaymentService(
	repo IPaymentRepo,
	userClient clients.UserServiceClient,
	fileClient clients.FileServiceClient,
	scheduleClient clients.ScheduleServiceClient,
) *PaymentService {

	return &PaymentService{
		repo:           repo,
		userClient:     userClient,
		fileClient:     fileClient,
		scheduleClient: scheduleClient,
	}
}

func (s *PaymentService) SubmitPaymentReceipt(ctx context.Context, input *models.SubmitPaymentReceiptInput) (*models.PaymentReceipt, error) {
	if input.FileId == uuid.Nil || input.LessonId == uuid.Nil {
		return nil, errdefs.ErrInvalidArgument
	}

	getLessonRequest := &api3.GetLessonRequest{
		Id: input.LessonId.String(),
	}

	lesson, err := s.scheduleClient.GetLesson(ctx, getLessonRequest)
	if err != nil {
		return nil, err
	}

	if lesson.IsPaid {
		return nil, errdefs.ErrAlreadyExists
	}

	lesson.IsPaid = true
	updateLessonRequest := &api3.UpdateLessonRequest{
		Id:             lesson.Id,
		ConnectionLink: lesson.ConnectionLink,
		PriceRub:       lesson.PriceRub,
		PaymentInfo:    lesson.PaymentInfo,
	}
	lesson, err = s.scheduleClient.UpdateLesson(ctx, updateLessonRequest)
	if err != nil {
		return nil, err
	}
	newReceiptID := uuid.New()
	exists, err := s.repo.ExistsByID(ctx, newReceiptID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errdefs.ErrAlreadyExists
	}

	createReceiptInput := &models.PaymentReceiptCreateInput{
		LessonID:   input.LessonId,
		FileID:     input.FileId,
		IsVerified: true,
	}
	receipt, err := s.repo.CreateReceipt(ctx, createReceiptInput)
	if err != nil {
		return nil, errdefs.ErrNotFound
	}

	// отправить ивент уведомление

	return receipt, nil
}

func (s *PaymentService) GetPaymentInfo(ctx context.Context, input *models.GetPaymentInfoInput) (*models.PaymentInfo, error) {
	if input.LessonId == uuid.Nil {
		return nil, errdefs.ErrInvalidArgument
	}

	getLessonRequest := &api3.GetLessonRequest{
		Id: input.LessonId.String(),
	}

	lesson, err := retry[*api3.Lesson](maxRetries, retryDelay, func() (*api3.Lesson, error) {
		return s.scheduleClient.GetLesson(ctx, getLessonRequest)
	})
	if err != nil {
		return nil, err
	}
	if lesson.PriceRub == 0 {
		lesson.PriceRub = 0
	}
	paymentInfo := &models.PaymentInfo{
		LessonID:       input.LessonId,
		PriceRUB:       lesson.PriceRub,
		PaymentDetails: lesson.PaymentInfo,
	}
	return paymentInfo, nil
}
func (s *PaymentService) GetReceipt(ctx context.Context, input *models.GetReceiptInput) (*models.PaymentReceipt, error) {
	if input.ReceiptId == uuid.Nil {
		return nil, errdefs.ErrInvalidArgument
	}

	receipt, err := retry[*models.PaymentReceipt](maxRetries, retryDelay, func() (*models.PaymentReceipt, error) {
		return s.repo.GetReceiptByID(ctx, input.ReceiptId)
	})
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (s *PaymentService) VerifyReceipt(ctx context.Context, input *models.VerifyReceipt) (*models.PaymentReceipt, error) {
	if input.ReceiptId == uuid.Nil {
		return nil, errdefs.ErrInvalidArgument
	}
	receipt, err := retry[*models.PaymentReceipt](maxRetries, retryDelay, func() (*models.PaymentReceipt, error) {
		return s.repo.UpdateReceipt(ctx, input.ReceiptId, true)
	})
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (s *PaymentService) GetReceiptFile(ctx context.Context, input *models.GetReceiptFileInput) (*models.ReceiptFileUrl, error) {
	if input.ReceiptId == uuid.Nil {
		return nil, errdefs.ErrInvalidArgument
	}
	receipt, err := retry[*models.PaymentReceipt](maxRetries, retryDelay, func() (*models.PaymentReceipt, error) {
		return s.repo.GetReceiptByID(ctx, input.ReceiptId)
	})
	if err != nil {
		return nil, err
	}
	generateDownloadURLRequest := &api2.GenerateDownloadURLRequest{FileId: receipt.FileID.String()}
	url, err := s.fileClient.GenerateDownloadURL(ctx, generateDownloadURLRequest)
	if err != nil {
		return nil, err
	}
	receiptFileURL := &models.ReceiptFileUrl{
		URL: url.GetUrl(),
	}
	return receiptFileURL, nil
}

func retry[T any](attempts int, delay time.Duration, fn func() (T, error)) (T, error) {
	var zero T
	var err error
	var result T

	for i := 0; i < attempts; i++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}
		time.Sleep(delay)
	}
	return zero, err
}
