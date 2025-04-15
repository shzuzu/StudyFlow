package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"payment_service/internal/domain"
	"payment_service/internal/repository"
	"payment_service/internal/service/mapper"
	pb "payment_service/proto"
	"time"
)

type PaymentServiceServer struct {
	pb.UnimplementedPaymentServiceServer
	repo *repository.PaymentRepo
}

func NewPaymentServiceServer(repo *repository.PaymentRepo) *PaymentServiceServer {
	return &PaymentServiceServer{repo: repo}
}

// GetPaymentInfo возвращает реквизиты и цену урока по lesson_id.
// 4. Возвращает PaymentInfo с данными из schedule_service и user_service.
func (s *) GetPaymentInfo(ctx context.Context, lessonID string) (*payment.PaymentInfo, error) {
	// 1. Извлекаем user_id из метаданных gRPC (API Gateway передаёт user_id и другие данные в metadata)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "нет метаданных")
	}
	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "user_id не найден")
	}
	studentID := userIDs[0] // пользователь, вызывающий данный метод, является учеником

	// 2. Получаем информацию об уроке из schedule_service.
	lessonResp, err := s.scheduleClient.GetLesson(ctx, &schedulepb.GetLessonRequest{
		LessonId: lessonID,
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "урок не найден: %v", err)
	}
	lesson := lessonResp.GetLesson()
	if lesson == nil {
		return nil, status.Error(codes.NotFound, "данные урока отсутствуют")
	}
	tutorID := lesson.GetTutorId()

	// 3. Обращаемся к user_service для получения контекста отношений.
	userCtxResp, err := s.userClient.ResolveTutorStudentContext(ctx, &userpb.ResolveTutorStudentContextRequest{
		TutorId:   tutorID,
		StudentId: studentID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось получить контекст отношений: %v", err)
	}
	// Проверяем, что отношения активны.
	if userCtxResp.GetRelationshipStatus() != "active" {
		return nil, status.Errorf(codes.PermissionDenied, "статус отношений '%s' не даёт права на оплату", userCtxResp.GetRelationshipStatus())
	}

	// Определяем цену урока: приоритетно берем цену из schedule_service, если она указана, иначе — из контекста отношений.
	var price int32
	if lesson.GetPriceRub() > 0 {
		price = lesson.GetPriceRub()
	} else {
		price = userCtxResp.GetLessonPriceRub()
	}

	// 4. Формируем и возвращаем объект PaymentInfo.
	return &payment.PaymentInfo{
		LessonId:    lessonID,
		PriceRub:    price,
		PaymentInfo: userCtxResp.GetPaymentInfo(),
	}, nil
}

// SubmitPaymentReceipt создаёт чек оплаты по данным из запроса и сохраняет его в БД.
func (s *PaymentServiceServer) SubmitPaymentReceipt(ctx context.Context, req *pb.SubmitPaymentReceiptRequest) (*pb.Receipt, error) {
	lessonUUID, err := uuid.Parse(req.LessonId)
	if err != nil {
		return nil, fmt.Errorf("invalid lesson_id: %v", err)
	}
	fileUUID, err := uuid.Parse(req.FileId)
	if err != nil {
		return nil, fmt.Errorf("invalid file_id: %v", err)
	}

	now := time.Now().UTC()
	receipt := &domain.PaymentReceipt{
		ID:         uuid.New(),
		LessonID:   lessonUUID,
		FileID:     fileUUID,
		IsVerified: false,
		CreatedAt:  now,
		EditedAt:   nil,
	}

	if err := s.repo.CreateReceipt(ctx, receipt); err != nil {
		return nil, fmt.Errorf("failed to create receipt: %v", err)
	}

	return mapper.MapToProtoReceipt(receipt), nil
}

// GetReceipt Получает чек по id из БД
func (s *PaymentServiceServer) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.Receipt, error) {
	receiptUUID, err := uuid.Parse(req.ReceiptId)
	if err != nil {
		return nil, fmt.Errorf("invalid receipt_id: %v", err)
	}

	receipt, err := s.repo.GetReceipt(ctx, receiptUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("receipt not found")
		}
		return nil, fmt.Errorf("failed to get receipt: %v", err)
	}

	return mapper.MapToProtoReceipt(receipt), nil
}
