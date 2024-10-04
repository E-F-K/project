package domain

import "context"

type UsersRepository interface {
	Create(context.Context, Connection, User) error
	Read(context.Context, Connection, UserID) (User, error)
	Update(context.Context, Connection, User) error
	Delete(context.Context, Connection, UserID) error
}

type ListsRepository interface {
	Create(context.Context, Connection, List) error
	Read(context.Context, Connection, ListID) (List, error)
	Update(context.Context, Connection, List) error
	Delete(context.Context, Connection, ListID) error
}

type TasksRepository interface {
	Create(context.Context, Connection, Task) error
	Read(context.Context, Connection, TaskID) (Task, error)
	Update(context.Context, Connection, Task) error
	Delete(context.Context, Connection, TaskID) error
}
