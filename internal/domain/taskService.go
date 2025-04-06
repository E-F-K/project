package domain

import (
	"context"
	"errors"
)

var (
	_ TaskInterface = (*TaskService)(nil)
)

var (
	ErrToDoServiceCreateTask = errors.Join(errToDoService, errors.New("create task failed"))
)

type TaskService struct {
	provider ConnectionProvider
	taskRepo TasksRepository
}

func NewTaskService(provider ConnectionProvider, taskRepo TasksRepository) *TaskService {
	return &TaskService{
		provider: provider,
		taskRepo: taskRepo,
	}
}

// Close implements TaskInterface.
func (s *TaskService) Close() error {
	return s.provider.Close()
}

// Create implements TaskInterface.
func (s *TaskService) Create(context.Context, UserID, Task) error {
	panic("unimplemented")
}

// Delete implements TaskInterface.
func (s *TaskService) Delete(context.Context, UserID, TaskID) error {
	panic("unimplemented")
}

// Update implements TaskInterface.
func (s *TaskService) Update(context.Context, UserID, Task) error {
	panic("unimplemented")
}
