package clients

import (
	"common_library/logging"
	"context"
	api2 "fileservice/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	api3 "schedule_service/pkg/api"
	api4 "userservice/pkg/api"
)

func New(ctx context.Context, url string) (*grpc.ClientConn, func()) {
	client, err := grpc.NewClient(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		if logger, ok := logging.GetFromContext(ctx); ok {
			logger.Fatal(ctx, "cannot create user grpc client", zap.Error(err))
		}
	}

	closeFunc := func() {
		err := client.Close()
		if err != nil {
			if logger, ok := logging.GetFromContext(ctx); ok {
				logger.Fatal(ctx, "cannot close user grpc client", zap.Error(err))
			}
		}
	}

	return client, closeFunc
}

type FileServiceClient interface {
	GenerateDownloadURL(ctx context.Context, req *api2.GenerateDownloadURLRequest, opts ...grpc.CallOption) (*api2.DownloadURL, error)
}

type ScheduleServiceClient interface {
	GetLesson(ctx context.Context, req *api3.GetLessonRequest, opts ...grpc.CallOption) (*api3.Lesson, error)
	CreateLesson(ctx context.Context, req *api3.CreateLessonRequest, opts ...grpc.CallOption) (*api3.Lesson, error)
	UpdateLesson(ctx context.Context, req *api3.UpdateLessonRequest, opts ...grpc.CallOption) (*api3.Lesson, error)
	CancelLesson(ctx context.Context, req *api3.CancelLessonRequest, opts ...grpc.CallOption) (*api3.Lesson, error)
}

type UserServiceClient interface {
	RegisterViaTelegram(ctx context.Context, req *api4.RegisterViaTelegramRequest, opts ...grpc.CallOption) (*api4.User, error)

	AuthorizeByAuthHeader(ctx context.Context, req *api4.AuthorizeByAuthHeaderRequest, opts ...grpc.CallOption) (*api4.User, error)

	GetMe(ctx context.Context, req *api4.Empty, opts ...grpc.CallOption) (*api4.User, error)

	GetUser(ctx context.Context, req *api4.GetUserRequest, opts ...grpc.CallOption) (*api4.UserPublic, error)

	UpdateUser(ctx context.Context, req *api4.UpdateUserRequest, opts ...grpc.CallOption) (*api4.User, error)

	UpdateTutorProfile(ctx context.Context, req *api4.UpdateTutorProfileRequest, opts ...grpc.CallOption) (*api4.TutorProfile, error)

	GetTutorProfileByUserId(ctx context.Context, req *api4.GetTutorProfileByUserIdRequest, opts ...grpc.CallOption) (*api4.TutorProfile, error)

	GetTutorStudent(ctx context.Context, req *api4.GetTutorStudentRequest, opts ...grpc.CallOption) (*api4.TutorStudent, error)

	CreateTutorStudent(ctx context.Context, req *api4.CreateTutorStudentRequest, opts ...grpc.CallOption) (*api4.TutorStudent, error)

	UpdateTutorStudent(ctx context.Context, req *api4.UpdateTutorStudentRequest, opts ...grpc.CallOption) (*api4.TutorStudent, error)

	DeleteTutorStudent(ctx context.Context, req *api4.DeleteTutorStudentRequest, opts ...grpc.CallOption) (*api4.Empty, error)

	ListTutorStudents(ctx context.Context, req *api4.ListTutorStudentsRequest, opts ...grpc.CallOption) (*api4.ListTutorStudentsResponse, error)

	ListTutorsForStudent(ctx context.Context, req *api4.ListTutorsForStudentRequest, opts ...grpc.CallOption) (*api4.ListTutorsForStudentResponse, error)

	ResolveTutorStudentContext(ctx context.Context, req *api4.ResolveTutorStudentContextRequest, opts ...grpc.CallOption) (*api4.ResolvedTutorStudentContext, error)

	AcceptInvitationFromTutor(ctx context.Context, req *api4.AcceptInvitationFromTutorRequest, opts ...grpc.CallOption) (*api4.Empty, error)
}
