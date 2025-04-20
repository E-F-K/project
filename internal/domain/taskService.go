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
	ErrToDoServiceDeleteTask = errors.Join(errToDoService, errors.New("delete task failed"))
	ErrToDoServiceUpdateTask = errors.Join(errToDoService, errors.New("update task failed"))
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
func (s *TaskService) Create(ctx context.Context, userID UserID, task Task) error {
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		task := Task{
			ID:        task.ID,
			ListID:    task.ListID,
			Priority:  task.Priority,
			Deadline:  task.Deadline,
			Done:      task.Done,
			Name:      task.Name,
			UpdatedAT: task.UpdatedAT,
		}

		return s.taskRepo.Create(ctx, connection, userID, task)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceCreateTask, err)
	}

	return nil
}

// Delete implements TaskInterface.
func (s *TaskService) Delete(ctx context.Context, userID UserID, taskID TaskID) error {
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		// listID ckeck for connection to the user? norm?
		return s.taskRepo.Delete(ctx, connection, userID, taskID)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceDeleteTask, err)
	}

	return nil
}

// Update implements TaskInterface.
func (s *TaskService) Update(ctx context.Context, userID UserID, task Task) error {
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		// listID ckeck for connection to the user? norm?
		return s.taskRepo.Update(ctx, connection, userID, task)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceUpdateTask, err)
	}

	return nil
}
