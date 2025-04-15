package schedule

import (
	"context"
	schedulepb "paymen/proto/schedule/v1"
)

type Client interface {
	GetLesson(ctx context.Context, lessonID string) (*schedulepb.Lesson, error)
	MarkLessonAsPaid(ctx context.Context, lessonID string) error
}

type client struct {
	grpcClient schedulepb.ScheduleServiceClient
}

func New(grpcClient schedulepb.ScheduleServiceClient) Client {
	return &client{grpcClient: grpcClient}
}

func (c *client) GetLesson(ctx context.Context, lessonID string) (*schedulepb.Lesson, error) {
	resp, err := c.grpcClient.GetLesson(ctx, &schedulepb.GetLessonRequest{LessonId: lessonID})
	if err != nil {
		return nil, err
	}
	return resp.Lesson, nil
}

func (c *client) MarkLessonAsPaid(ctx context.Context, lessonID string) error {
	_, err := c.grpcClient.MarkLessonAsPaid(ctx, &schedulepb.MarkLessonAsPaidRequest{LessonId: lessonID})
	return err
}
