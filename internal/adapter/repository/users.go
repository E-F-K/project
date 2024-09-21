package repository

import (
	"context"
	"errors"
	"todo_list/internal/domain"
)

var _ domain.UsersRepository = (*Users)(nil)

var (
	errUsers       = errors.New("users repository error")
	ErrUsersCreate = errors.Join(errUsers, errors.New("create failed"))
	ErrUsersRead   = errors.Join(errUsers, errors.New("read failed"))
	ErrUsersUpdate = errors.Join(errUsers, errors.New("update failed"))
	ErrUsersDelete = errors.Join(errUsers, errors.New("delete failed"))
)

type Users struct{}

func NewUsers() *Users {
	return &Users{}
}

func (r Users) Create(ctx context.Context, connection domain.Connection, user domain.User) error {
	const query = `
insert into users
    (id, name, email, token)
values
    ($1, $2, $3, $4)`

	_, err := connection.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Token)
	if err != nil {
		err = errors.Join(ErrUsersCreate, err)
	}

	return err
}

func (r Users) Delete(ctx context.Context, connection domain.Connection, userID domain.UserID) error {
	const query = `delete from users where id = $1`

	_, err := connection.ExecContext(ctx, query, userID)
	if err != nil {
		err = errors.Join(ErrUsersDelete, err)
	}

	return err
}

func (r Users) Read(ctx context.Context, connection domain.Connection, userID domain.UserID) (domain.User, error) {
	const query = `select id, name, email, token from users where id = $1`

	var user domain.User
	err := connection.GetContext(ctx, &user, query, userID)
	if err != nil {
		err = errors.Join(ErrUsersRead, err)
	}

	return user, err
}

func (r Users) Update(ctx context.Context, connection domain.Connection, user domain.User) error {
	const query = `update users set name = $2, email = $3, token = $4 where id = $1`

	_, err := connection.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Token)
	if err != nil {
		err = errors.Join(ErrUsersUpdate, err)
	}

	return err
}
