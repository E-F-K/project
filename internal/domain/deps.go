package domain

import "context"

type UsersRepository interface {
	Create(context.Context, Connection, User) error
	ReadByToken(context.Context, Connection, string) (User, error)
	ReadByEmail(context.Context, Connection, string) (User, error)
	Update(context.Context, Connection, User) error
	Delete(context.Context, Connection, UserID) error
	UpdateTokenByEmail(context.Context, Connection, string, string) error
}

type ListsRepository interface {
	Create(context.Context, Connection, List) error
	Update(context.Context, Connection, List) error
	Delete(context.Context, Connection, UserID, ListID) error
	ReadAll(context.Context, Connection, UserID) ([]List, error)
}

type TasksRepository interface {
	Create(context.Context, Connection, UserID, Task) error
	Read(context.Context, Connection, TaskID) (Task, error)
	Update(context.Context, Connection, UserID, Task) error
	Delete(context.Context, Connection, UserID, TaskID) error
	GetAllTasks(context.Context, Connection, []ListID) ([]Task, error)
}
