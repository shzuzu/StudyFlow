package service

import (
	"context"
	api2 "fileservice/pkg/api"
	"github.com/google/uuid"
	errdefs "paymentservice/internal/errors"
	"paymentservice/internal/models"
	api3 "schedule_service/pkg/api"
	api4 "userservice/pkg/api"
)

type IPaymentRepo interface {
	CreateReceipt(ctx context.Context, receipt *models.PaymentReceiptCreateInput) (*models.PaymentReceipt, error)

	GetReceiptByID(ctx context.Context, id uuid.UUID) (*models.PaymentReceipt, error)

	UpdateReceipt(ctx context.Context, id uuid.UUID, isVerified bool) (*models.PaymentReceipt, error)

	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)

	GetReceiptByLessonID(ctx context.Context, lessonID uuid.UUID) (*models.PaymentReceipt, error)
}

type PaymentService struct {
	repo           IPaymentRepo
	userClient     api4.UserServiceClient
	fileClient     api2.FileServiceClient
	scheduleClient api3.ScheduleServiceClient
}

func NewPaymentService(
	repo IPaymentRepo,
	userClient api4.UserServiceClient,
	fileClient api2.FileServiceClient,
	scheduleClient api3.ScheduleServiceClient,
) *PaymentService {

	return &PaymentService{
		repo:           repo,
		userClient:     userClient,
		fileClient:     fileClient,
		scheduleClient: scheduleClient,
	}
}

func (s *PaymentService) SubmitPaymentReceipt(ctx context.Context, input *models.SubmitPaymentReceiptInput) (*models.PaymentReceipt, error) {

	getLessonRequest := &api3.GetLessonRequest{
		Id: input.LessonId.String(),
	}

	lesson, err := s.scheduleClient.GetLesson(ctx, getLessonRequest)
	if err != nil {
		return nil, err
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

	receipt, err := s.repo.GetReceiptByLessonID(ctx, input.LessonId)
	if err != nil {
		return nil, errdefs.ErrNotFound
	}

	// отправить ивент уведомление

	return receipt, nil
}

func (s *PaymentService) GetPaymentInfo(ctx context.Context, input *models.GetPaymentInfoInput) (*models.PaymentInfo, error) {
	getLessonRequest := &api3.GetLessonRequest{
		Id: input.LessonId.String(),
	}

	lesson, err := s.scheduleClient.GetLesson(ctx, getLessonRequest)
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
	receipt, err := s.repo.GetReceiptByID(ctx, input.ReceiptId)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (s *PaymentService) VerifyReceipt(ctx context.Context, input *models.VerifyReceipt) (*models.PaymentReceipt, error) {
	receipt, err := s.repo.UpdateReceipt(ctx, input.ReceiptId, true)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (s *PaymentService) GetReceiptFile(ctx context.Context, input *models.GetReceiptFileInput) (*models.ReceiptFileUrl, error) {
	receipt, err := s.repo.GetReceiptByID(ctx, input.ReceiptId)
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
