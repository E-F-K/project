package repository

import (
	"context"
	"errors"
	"time"

	"todo_list/internal/domain"
)

var _ domain.TasksRepository = (*Tasks)(nil)

var (
	errTasks       = errors.New("tasks repository error")
	ErrTasksCreate = errors.Join(errTasks, errors.New("create failed"))
	ErrTasksRead   = errors.Join(errTasks, errors.New("read failed"))
	ErrTasksUpdate = errors.Join(errTasks, errors.New("update failed"))
	ErrTasksDelete = errors.Join(errTasks, errors.New("delete failed"))
)

type Tasks struct{}

func NewTasks() *Tasks {
	return &Tasks{}
}

func (r Tasks) Create(ctx context.Context, connection domain.Connection, task domain.Task) error {
	const query = `
insert into tasks
    (id, list_id, priority, deadline, done, name, updated_at)
values
    ($1, $2, $3, $4, $5, $6, $7)`

	_, err := connection.ExecContext(ctx, query, task.ID, task.ListID, domain.Priority(task.Priority), task.Deadline, task.Done, task.Name, time.Time(task.UpdatedAT))
	if err != nil {
		err = errors.Join(ErrTasksCreate, err)
	}

	return err
}

func (r Tasks) Delete(ctx context.Context, connection domain.Connection, taskID domain.TaskID) error {
	const query = `delete from tasks where id = $1`

	_, err := connection.ExecContext(ctx, query, taskID)
	if err != nil {
		err = errors.Join(ErrTasksDelete, err)
	}

	return err
}

func (r Tasks) Read(ctx context.Context, connection domain.Connection, taskID domain.TaskID) (domain.Task, error) {
	const query = `select id, list_id, priority, deadline, done, name, updated_at from tasks where id = $1`

	var task domain.Task
	err := connection.GetContext(ctx, &task, query, taskID)
	if err != nil {
		err = errors.Join(ErrTasksRead, err)
	}

	task.ID = taskID

	return task, err
}

func (r Tasks) Update(ctx context.Context, connection domain.Connection, task domain.Task) error {
	const query = `update tasks set name = $2, priority = $3, deadline = $4, done = $5, updated_at = $6 where id = $1`

	_, err := connection.ExecContext(ctx, query, task.ID, task.Name, domain.Priority(task.Priority), task.Deadline, task.Done, time.Time(task.UpdatedAT))
	if err != nil {
		err = errors.Join(ErrTasksUpdate, err)
	}

	return err
}

func (r Tasks) GetAllTasks(ctx context.Context, connection domain.Connection, listsIDs []domain.ListID) ([]domain.Task, error) {
	const query = `select id, list_id, priority, deadline, done, name, updated_at from tasks where list_id = any($1)`

	var tasks []domain.Task
	err := connection.SelectContext(ctx, &tasks, query, listsIDs)
	if err != nil {
		err = errors.Join(ErrUsersGetAllTasks, err)
	}

	return tasks, err
}
