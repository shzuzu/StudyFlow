package handler_test

import (
	"context"
	errdefs "paymentservice/internal/errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"paymentservice/internal/handler"
	"paymentservice/internal/models"
	pb "paymentservice/pkg/api"
)

type MockService struct {
	mock.Mock
	handler.PaymentService
}

func (m *MockService) GetPaymentInfo(ctx context.Context, input *models.GetPaymentInfoInput) (*models.PaymentInfo, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.PaymentInfo), args.Error(1)
}

func TestGetPaymentInfo_InvalidUUID(t *testing.T) {
	s := handler.NewPaymentServiceServer(nil)
	_, err := s.GetPaymentInfo(context.Background(), &pb.GetPaymentInfoRequest{LessonId: "invalid"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetPaymentInfo_Success(t *testing.T) {
	lessonID := uuid.New()
	mockSvc := new(MockService)
	s := handler.NewPaymentServiceServer(mockSvc)
	ctx := context.Background()
	input := &models.GetPaymentInfoInput{LessonId: lessonID}
	mockSvc.On("GetPaymentInfo", ctx, input).Return(&models.PaymentInfo{
		LessonID:       lessonID,
		PriceRUB:       1000, // int32
		PaymentDetails: "details",
	}, nil)

	resp, err := s.GetPaymentInfo(ctx, &pb.GetPaymentInfoRequest{LessonId: lessonID.String()})
	assert.NoError(t, err)
	assert.Equal(t, lessonID.String(), resp.LessonId)
	assert.Equal(t, int32(1000), resp.PriceRub)
	assert.Equal(t, "details", resp.PaymentInfo)

}

func TestGetPaymentInfo_NotFound(t *testing.T) {
	lessonID := uuid.New()
	mockSvc := new(MockService)
	s := handler.NewPaymentServiceServer(mockSvc)
	ctx := context.Background()
	input := &models.GetPaymentInfoInput{LessonId: lessonID}

	// Возвращаем ошибку ErrNotFound
	mockSvc.On("GetPaymentInfo", ctx, input).Return(nil, errdefs.ErrNotFound)

	_, err := s.GetPaymentInfo(ctx, &pb.GetPaymentInfoRequest{LessonId: lessonID.String()})

	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	mockSvc.AssertExpectations(t)
}
