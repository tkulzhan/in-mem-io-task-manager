package service_test

import (
	"context"
	"errors"
	"testing"

	"in-mem-io-task-manager/internal/domain/service"
	mock_service "in-mem-io-task-manager/internal/domain/service/mocks"
	"in-mem-io-task-manager/internal/infrastructure/http/server/generated"
	"in-mem-io-task-manager/internal/infrastructure/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newLogger() *logger.Logger {
	return logger.New("debug")
}

func newService(t *testing.T, repo *mock_service.MockTaskRepository) *service.TaskManagerService {
	conf := service.Configuration{
		L:                  newLogger(),
		TaskRepository:     repo,
		MaxConcurrentTasks: 2,
	}
	svc, err := service.NewTaskManagerService(conf)
	require.NoError(t, err)
	return svc
}

func newMockTask(ctrl *gomock.Controller) *mock_service.MockTask {
	m := mock_service.NewMockTask(ctrl)
	m.EXPECT().GetID().AnyTimes().Return("mock-task-id")
	return m
}

func TestTaskManagerService_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_service.NewMockTaskRepository(ctrl)
	mockTask := newMockTask(ctrl)

	svc := newService(t, mockRepo)

	tests := []struct {
		name      string
		setupMock func()
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockRepo.EXPECT().
					CreateTask(gomock.Any(), gomock.AssignableToTypeOf(service.RunningTask{})).
					Return(nil)
				mockTask.EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo fails",
			setupMock: func() {
				mockRepo.EXPECT().
					CreateTask(gomock.Any(), gomock.AssignableToTypeOf(service.RunningTask{})).
					Return(errors.New("repo error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := svc.CreateTask(context.Background(), mockTask)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTaskManagerService_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_service.NewMockTaskRepository(ctrl)
	mockTask := newMockTask(ctrl)

	svc := newService(t, mockRepo)

	tests := []struct {
		name      string
		setupMock func()
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() {
				rt := service.RunningTask{
					Task:   mockTask,
					Cancel: func() {},
				}
				mockRepo.EXPECT().GetTaskById(gomock.Any(), "id").Return(rt, nil)
				mockRepo.EXPECT().DeleteTask(gomock.Any(), "id").Return(nil)
			},
			wantErr: false,
		},
		{
			name: "get fails",
			setupMock: func() {
				mockRepo.EXPECT().GetTaskById(gomock.Any(), "id").Return(service.RunningTask{}, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "delete fails",
			setupMock: func() {
				rt := service.RunningTask{
					Task:   mockTask,
					Cancel: func() {},
				}
				mockRepo.EXPECT().GetTaskById(gomock.Any(), "id").Return(rt, nil)
				mockRepo.EXPECT().DeleteTask(gomock.Any(), "id").Return(errors.New("delete error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := svc.DeleteTask(context.Background(), "id")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTaskManagerService_GetTaskById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_service.NewMockTaskRepository(ctrl)
	mockTask := newMockTask(ctrl)

	svc := newService(t, mockRepo)

	tests := []struct {
		name      string
		setupMock func()
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() {
				rt := service.RunningTask{Task: mockTask}
				mockRepo.EXPECT().GetTaskById(gomock.Any(), "id").Return(rt, nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			setupMock: func() {
				mockRepo.EXPECT().GetTaskById(gomock.Any(), "id").Return(service.RunningTask{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := svc.GetTaskById(context.Background(), "id")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTaskManagerService_ToTaskObject(t *testing.T) {
	svc, err := service.NewTaskManagerService(service.Configuration{
		L:                  newLogger(),
		TaskRepository:     &mock_service.MockTaskRepository{},
		MaxConcurrentTasks: 1,
	})
	require.NoError(t, err)

	desc := "desc"
	tests := []struct {
		name    string
		req     generated.CreateTaskRequest
		wantErr bool
	}{
		{
			name: "default",
			req: generated.CreateTaskRequest{
				Data: map[string]interface{}{
					"title": "Test Task",
				},
			},
			wantErr: false,
		},
		{
			name: "unsupported",
			req: generated.CreateTaskRequest{
				Type: ptr("other"),
			},
			wantErr: true,
		},
		{
			name: "invalid input",
			req: generated.CreateTaskRequest{
				Data: map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "default with description",
			req: generated.CreateTaskRequest{
				Data: map[string]interface{}{
					"title":       "Task",
					"description": desc,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := svc.ToTaskObject(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, task)
			}
		})
	}
}

func ptr(s string) *string { return &s }
