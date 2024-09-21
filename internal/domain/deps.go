package domain

import "context"

type UsersRepository interface {
	Create(context.Context, Connection, User) error
	Read(context.Context, Connection, UserID) (User, error)
	Update(context.Context, Connection, User) error
	Delete(context.Context, Connection, UserID) error
}
