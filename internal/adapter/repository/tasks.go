package repository

import (
	"context"
	"database/sql"
	"errors"

	"todo_list/internal/domain"
)

var _ domain.TasksRepository = (*Tasks)(nil)

var (
	errTasks            = errors.New("tasks repository error")
	ErrTasksCreate      = errors.Join(errTasks, errors.New("create failed"))
	ErrTasksRead        = errors.Join(errTasks, errors.New("read failed"))
	ErrTasksUpdate      = errors.Join(errTasks, errors.New("update failed"))
	ErrTasksDelete      = errors.Join(errTasks, errors.New("delete failed"))
	ErrTasksGetAllTasks = errors.Join(errTasks, errors.New("get all failed"))
)

type Tasks struct{}

func NewTasks() *Tasks {
	return &Tasks{}
}

func (r Tasks) listExists(ctx context.Context, connection domain.Connection, userID domain.UserID, listID domain.ListID) (bool, error) {

	const query = `select 1 from lists where user_id = $1 and id = $2`

	var tmp int
	if err := connection.GetContext(ctx, &tmp, query, userID, listID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}

		return false, err
	}

	return true, nil
}

func (r Tasks) Create(ctx context.Context, connection domain.Connection, userID domain.UserID, task domain.Task) error {
	exists, err := r.listExists(ctx, connection, userID, task.ListID)
	if err != nil {
		return errors.Join(ErrTasksCreate, err)
	}
	if !exists {
		return errors.Join(ErrTasksCreate, errors.New("list not found or access denied"))
	}

	const query = `insert into tasks (id, list_id, priority, deadline, done, name) values ($1, $2, $3, $4, $5, $6)`

	_, err = connection.ExecContext(ctx, query, task.ID, task.ListID, domain.Priority(task.Priority), task.Deadline, task.Done, task.Name)
	if err != nil {
		return errors.Join(ErrTasksCreate, err)
	}

	return nil
}

func (r Tasks) Delete(ctx context.Context, connection domain.Connection, userID domain.UserID, taskID domain.TaskID) error {
	var listID domain.ListID
	if err := connection.GetContext(ctx, &listID, "select list_id from tasks where id = $1", taskID); err != nil {
		return errors.Join(ErrTasksDelete, err)
	}

	exists, err := r.listExists(ctx, connection, userID, listID)
	if err != nil {
		return errors.Join(ErrTasksDelete, err)
	}
	if !exists {
		return errors.Join(ErrTasksDelete, errors.New("list not found or access denied"))
	}

	const query = `delete from tasks where id = $1`

	_, err = connection.ExecContext(ctx, query, taskID)
	if err != nil {
		return errors.Join(ErrTasksDelete, err)
	}

	return nil
}

func (r Tasks) Read(ctx context.Context, connection domain.Connection, userID domain.UserID, taskID domain.TaskID) (domain.Task, error) {
	var task domain.Task
	var listID domain.ListID
	if err := connection.GetContext(ctx, &listID, "select list_id from tasks where id = $1", taskID); err != nil {
		return task, errors.Join(ErrTasksRead, err)
	}

	exists, err := r.listExists(ctx, connection, userID, listID)
	if err != nil {
		return task, errors.Join(ErrTasksRead, err)
	}
	if !exists {
		return task, errors.Join(ErrTasksRead, errors.New("list not found or access denied"))
	}

	const query = `select id, list_id, priority, deadline, done, name, updated_at from tasks where id = $1`

	err = connection.GetContext(ctx, &task, query, taskID)
	if err != nil {
		return task, errors.Join(ErrTasksRead, err)
	}

	task.ID = taskID

	return task, nil
}

func (r Tasks) Update(ctx context.Context, connection domain.Connection, userID domain.UserID, task domain.Task) error {
	exists, err := r.listExists(ctx, connection, userID, task.ListID)
	if err != nil {
		return errors.Join(ErrTasksUpdate, err)
	}
	if !exists {
		return errors.Join(ErrTasksUpdate, errors.New("list not found or access denied"))
	}

	const query = `update tasks set name = $2, priority = $3, deadline = $4, done = $5, updated_at = default where id = $1`

	_, err = connection.ExecContext(ctx, query, task.ID, task.Name, domain.Priority(task.Priority), task.Deadline, task.Done)
	if err != nil {
		return errors.Join(ErrTasksUpdate, err)
	}

	return nil
}

func (r Tasks) GetAllTasks(ctx context.Context, connection domain.Connection, userID domain.UserID, listsIDs []domain.ListID) ([]domain.Task, error) {
	// userID ckeck for all lists!!!
	for _, listID := range listsIDs {
		exists, err := r.listExists(ctx, connection, userID, listID)
		if err != nil {
			return nil, errors.Join(ErrTasksGetAllTasks, err)
		}
		if !exists {
			return nil, errors.Join(ErrTasksGetAllTasks, errors.New("list not found or access denied"))
		}
	}

	const query = `select id, list_id, priority, deadline, done, name, updated_at from tasks where list_id = any($1)`

	var tasks []domain.Task
	err := connection.SelectContext(ctx, &tasks, query, listsIDs)
	if err != nil {
		return nil, errors.Join(ErrTasksGetAllTasks, err)
	}

	return tasks, nil
}
