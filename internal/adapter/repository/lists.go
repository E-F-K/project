package repository

import (
	"context"
	"errors"

	"todo_list/internal/domain"
)

var _ domain.ListsRepository = (*Lists)(nil)

var (
	errLists        = errors.New("lists repository error")
	ErrListsCreate  = errors.Join(errLists, errors.New("create failed"))
	ErrListsRead    = errors.Join(errLists, errors.New("read failed"))
	ErrListsUpdate  = errors.Join(errLists, errors.New("update failed"))
	ErrListsDelete  = errors.Join(errLists, errors.New("delete failed"))
	ErrListsReadAll = errors.Join(errLists, errors.New("read all failed"))
)

type Lists struct{}

func NewLists() *Lists {
	return &Lists{}
}

func (r Lists) Create(ctx context.Context, connection domain.Connection, list domain.List) error {
	const query = `insert into lists
    (id, user_id, name, updated_at)
	values
    ($1, $2, $3, default)`

	_, err := connection.ExecContext(ctx, query, list.ID, list.UserID, list.Name)
	if err != nil {
		return errors.Join(ErrListsCreate, err)
	}

	return nil
}

func (r Lists) Delete(ctx context.Context, connection domain.Connection, userID domain.UserID, listID domain.ListID) error {
	const query = `delete from lists where user_id = $1 and id = $2`

	if _, err := connection.ExecContext(ctx, query, userID, listID); err != nil {
		return errors.Join(ErrListsDelete, err)
	}

	return nil
}

func (r Lists) Read(ctx context.Context, connection domain.Connection, userID domain.UserID, listID domain.ListID) (domain.List, error) {
	const query = `select id, user_id, name, updated_at from lists where user_id = $1 and id = $2`

	var list domain.List
	if err := connection.GetContext(ctx, &list, query, userID, listID); err != nil {
		return list, errors.Join(ErrListsRead, err)
	}

	return list, nil
}

func (r Lists) Update(ctx context.Context, connection domain.Connection, list domain.List) error {
	const query = `update lists set name = $3, updated_at = default where user_id = $1 and id = $2`

	if list.Tasks != nil {
		return errors.Join(ErrListsUpdate, errors.New("task updates are not supported"))
	}

	if _, err := connection.ExecContext(ctx, query, list.UserID, list.ID, list.Name); err != nil {
		return errors.Join(ErrListsUpdate, err)
	}

	return nil
}

func (r Lists) ReadAll(ctx context.Context, connection domain.Connection, userID domain.UserID) ([]domain.List, error) {
	const query = `select id, user_id, name, updated_at from lists where user_id = $1`

	var lists []domain.List
	if err := connection.SelectContext(ctx, &lists, query, userID); err != nil {
		return nil, errors.Join(ErrListsReadAll, err)
	}

	return lists, nil
}
