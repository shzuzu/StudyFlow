package service_test

import (
	"common_library/ctxdata"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"schedule_service/internal/database/repo"
	"schedule_service/internal/service/service"
	pb "schedule_service/pkg/api"
	"schedule_service/pkg/mocks"
)

func setup(t *testing.T) (*service.ScheduleServer, *mocks.MockRepository, *mocks.MockIUserClient, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mocks.NewMockRepository(ctrl)
	mockUserClient := mocks.NewMockIUserClient(ctrl)
	srv := service.NewScheduleServer(mockRepo, mockUserClient)

	return srv, mockRepo, mockUserClient, ctrl
}

func TestGetSlot(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv, mockRepo, _, _ := setup(t)
		ctx := ctxdata.WithUserID(context.Background(), "123")

		slotID := "5d2d6f89-5bfa-4d85-bd3e-fcb4d6f56f20"
		now := time.Now()

		mockRepo.EXPECT().GetSlot(ctx, slotID).Return(&repo.Slot{
			ID:        slotID,
			TutorID:   "123",
			StartsAt:  now,
			EndsAt:    now.Add(time.Hour),
			IsBooked:  false,
			CreatedAt: now,
			EditedAt:  &now,
		}, nil)

		resp, err := srv.GetSlot(ctx, &pb.GetSlotRequest{Id: slotID})
		require.NoError(t, err)
		require.Equal(t, slotID, resp.Id)
		require.Equal(t, "123", resp.TutorId)
		require.False(t, resp.IsBooked)
		require.WithinDuration(t, now, resp.StartsAt.AsTime(), time.Second)
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		srv, _, _, _ := setup(t)
		ctx := context.Background()

		_, err := srv.GetSlot(ctx, &pb.GetSlotRequest{Id: "abc"})
		require.Error(t, err)
		st, _ := status.FromError(err)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		srv, _, _, _ := setup(t)
		ctx := ctxdata.WithUserID(context.Background(), "1")

		_, err := srv.GetSlot(ctx, &pb.GetSlotRequest{Id: "not-a-uuid"})
		require.Error(t, err)
		st, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("Slot Not Found", func(t *testing.T) {
		srv, mockRepo, _, _ := setup(t)
		ctx := ctxdata.WithUserID(context.Background(), "123")

		slotID := "de305d54-75b4-431b-adb2-eb6b9e546014"
		mockRepo.EXPECT().GetSlot(ctx, slotID).Return(nil, service.ErrSlotNotFound)

		_, err := srv.GetSlot(ctx, &pb.GetSlotRequest{Id: slotID})
		require.Error(t, err)
		st, _ := status.FromError(err)
		require.Equal(t, codes.NotFound, st.Code())
	})
}
