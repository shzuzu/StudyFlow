package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"payment_service/internal/domain"
	pb "payment_service/proto"
)

func MapToProtoReceipt(r *domain.PaymentReceipt) *pb.Receipt {
	var editedAt *timestamppb.Timestamp
	if r.EditedAt != nil {
		editedAt = timestamppb.New(*r.EditedAt)
	}

	return &pb.Receipt{
		Id:         r.ID.String(),
		LessonId:   r.LessonID.String(),
		FileId:     r.FileID.String(),
		IsVerified: r.IsVerified,
		CreatedAt:  timestamppb.New(r.CreatedAt),
		EditedAt:   editedAt,
	}
}
