package repository

import (
	"context"
	"errors"
	"time"

	"todo_list/internal/domain"
)

var _ domain.ListsRepository = (*Lists)(nil)

var (
	errLists       = errors.New("lists repository error")
	ErrListsCreate = errors.Join(errLists, errors.New("create failed"))
	ErrListsRead   = errors.Join(errLists, errors.New("read failed"))
	ErrListsUpdate = errors.Join(errLists, errors.New("update failed"))
	ErrListsDelete = errors.Join(errLists, errors.New("delete failed"))
)

type Lists struct{}

func NewLists() *Lists {
	return &Lists{}
}

func (r Lists) Create(ctx context.Context, connection domain.Connection, list domain.List) error {
	const query = `
insert into lists
    (id, user_id, name, email, updated_at)
values
    ($1, $2, $3, $4, $5)`

	_, err := connection.ExecContext(ctx, query, list.ID, list.UserID, list.Name, list.Email, time.Now())
	if err != nil {
		err = errors.Join(ErrListsCreate, err)
	}

	return err
}

func (r Lists) Delete(ctx context.Context, connection domain.Connection, ListID domain.ListID) error {
	const query = `delete from lists where id = $1`

	_, err := connection.ExecContext(ctx, query, ListID)
	if err != nil {
		err = errors.Join(ErrListsDelete, err)
	}

	return err
}

func (r Lists) Read(ctx context.Context, connection domain.Connection, listID domain.ListID) (domain.List, error) {
	const query = `select user_id, name, email, updated_at from lists where id = $1`

	var list domain.List
	err := connection.GetContext(ctx, &list, query, listID)
	if err != nil {
		err = errors.Join(ErrListsRead, err)
	}

	list.ID = listID

	return list, err
}

func (r Lists) Update(ctx context.Context, connection domain.Connection, list domain.List) error {
	const query = `update lists set name = $2, email = $3, updated_at = $4 where id = $1`

	_, err := connection.ExecContext(ctx, query, list.ID, list.Name, list.Email, time.Time(list.UpdatedAT))
	if err != nil {
		err = errors.Join(ErrListsUpdate, err)
	}

	return err
}
