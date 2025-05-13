package repository

import (
	"context"
	"errors"

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
)

type Users struct{}

func NewUsers() *Users {
	return &Users{}
}

func (r Users) Create(ctx context.Context, connection domain.Connection, user domain.User) error {
	const query = `
insert into users
    (id, name, email, password_hash, token)
values
    ($1, $2, $3, $4, $5)`

	_, err := connection.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Token)
	if err != nil {
		return errors.Join(ErrUsersCreate, err)
	}

	return nil
}

func (r Users) Delete(ctx context.Context, connection domain.Connection, userID domain.UserID) error {
	const query = `delete from users where id = $1`

	_, err := connection.ExecContext(ctx, query, userID)
	if err != nil {
		return errors.Join(ErrUsersDelete, err)
	}

	return nil
}

func (r Users) ReadByToken(ctx context.Context, connection domain.Connection, token string) (domain.User, error) {
	const query = `select id, name, email, password_hash, token from users where token = $1`

	var user domain.User
	err := connection.GetContext(ctx, &user, query, token)
	if err != nil {
		return user, errors.Join(ErrUsersRead, err)
	}

	return user, nil
}

func (r Users) ReadByEmail(ctx context.Context, connection domain.Connection, email string) (domain.User, error) {
	const query = `select id, name, email, password_hash, token from users where email = $1`

	var user domain.User
	err := connection.GetContext(ctx, &user, query, email)
	if err != nil {
		return user, errors.Join(ErrUsersRead, err)
	}

	return user, nil
}

func (r Users) Update(ctx context.Context, connection domain.Connection, user domain.User) error {
	const query = `update users set name = $2, email = $3, password_hash = $4, token = $5, updated_at = default where id = $1`

	_, err := connection.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Token)
	if err != nil {
		return errors.Join(ErrUsersUpdate, err)
	}

	return nil
}

func (r Users) UpdateTokenByEmail(ctx context.Context, connection domain.Connection, email string, token string) error {
	const query = `update users set token = $2 where email = $1`
	updated, err := connection.ExecContext(ctx, query, email, token)
	if err != nil || updated <= 0 {
		return errors.Join(ErrUsersUpdateTokenByEmail, err)
	}

	return nil
}
