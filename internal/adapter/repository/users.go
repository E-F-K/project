package repository

import (
	"context"
	"errors"
	"time"
	"todo_list/internal/domain"
)

var _ domain.UsersRepository = (*Users)(nil)

var (
	errUsers                   = errors.New("users repository error")
	ErrUsersCreate             = errors.Join(errUsers, errors.New("create failed"))
	ErrUsersRead               = errors.Join(errUsers, errors.New("read failed"))
	ErrUsersUpdate             = errors.Join(errUsers, errors.New("update failed"))
	ErrUsersDelete             = errors.Join(errUsers, errors.New("delete failed"))
	ErrUsersUpdateTokenByEmail = errors.Join(errUsers, errors.New("update token by email failed"))
	ErrUsersGetAllLists        = errors.Join(errUsers, errors.New("get all lists failed"))
	ErrUsersGetAllTasks        = errors.Join(errUsers, errors.New("get all tasks failed"))
)

type Users struct{}

func NewUsers() *Users {
	return &Users{}
}

func (r Users) Create(ctx context.Context, connection domain.Connection, user domain.User) error {
	const query = `
insert into users
    (id, name, email, password_hash, token, updated_at)
values
    ($1, $2, $3, $4, $5, $6)`

	_, err := connection.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Token, time.Now())
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

func (r Users) ReadByToken(ctx context.Context, connection domain.Connection, token string) (domain.User, error) {
	const query = `select id, name, email, password_hash, token from users where token = $1`

	var user domain.User
	err := connection.GetContext(ctx, &user, query, token)
	if err != nil {
		err = errors.Join(ErrUsersRead, err)
	}

	return user, err
}

func (r Users) ReadByEmail(ctx context.Context, connection domain.Connection, email string) (domain.User, error) {
	const query = `select id, name, email, password_hash, token from users where email = $1`

	var user domain.User
	err := connection.GetContext(ctx, &user, query, email)
	if err != nil {
		err = errors.Join(ErrUsersRead, err)
	}

	return user, err
}

func (r Users) Update(ctx context.Context, connection domain.Connection, user domain.User) error {
	const query = `update users set name = $2, email = $3, password_hash = $4, token = $5 where id = $1`

	_, err := connection.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Token)
	if err != nil {
		err = errors.Join(ErrUsersUpdate, err)
	}

	return err
}

func (r Users) UpdateTokenByEmail(ctx context.Context, connection domain.Connection, email string, token string) error {
	const query = `update users set token = $2 where email = $1`
	updated, err := connection.ExecContext(ctx, query, email, token)
	if err != nil || updated <= 0 {
		err = errors.Join(ErrUsersUpdateTokenByEmail, err)
	}

	return err
}

func (r Users) GetAllLists(ctx context.Context, connection domain.Connection, userID domain.UserID) ([]domain.List, error) {
	const query = `select user_id, name, email, updated_at from lists where user_id = $1`
	var lists []domain.List
	err := connection.SelectContext(ctx, &lists, query, userID)
	if err != nil {
		err = errors.Join(ErrUsersGetAllLists, err)
	}

	return lists, err
}

func (r Users) GetAllTasks(ctx context.Context, connection domain.Connection, listsID []domain.ListID) ([]domain.Task, error) {
	const query = `select id, list_id, priority, deadline, done, name, updated_at from tasks where list_id in ($1)`
	var tasks []domain.Task

	err := connection.SelectContext(ctx, &tasks, query, listsID)
	if err != nil {
		err = errors.Join(ErrUsersGetAllTasks, err)
	}

	return tasks, err
}
