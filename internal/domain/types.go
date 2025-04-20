package domain

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

type (
	UserID = uuid.UUID

	User struct {
		ID           UserID
		Name         string
		Email        string
		PasswordHash string
		Token        string
	}

	ListID = uuid.UUID

	List struct {
		ID        ListID    `json:"id"`
		UserID    UserID    `json:"user_id,omitempty"`
		Name      string    `json:"name"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
		Tasks     []Task    `json:"tasks,omitempty"`
	}

	TaskID = uuid.UUID

	Task struct {
		ID        TaskID     `json:"id"`
		ListID    ListID     `json:"list_id"`
		Priority  Priority   `json:"priority"`
		Deadline  *time.Time `json:"deadline"`
		Done      bool       `json:"done,omitempty"`
		Name      string     `json:"name"`
		UpdatedAT time.Time  `json:"updated_at,omitempty"`
	}

	Connection interface {
		GetContext(context.Context, any, string, ...any) error
		SelectContext(context.Context, any, string, ...any) error
		ExecContext(context.Context, string, ...any) (int64, error)
	}

	ConnectionProvider interface {
		Execute(context.Context, func(context.Context, Connection) error) error
		ExecuteTx(context.Context, func(context.Context, Connection) error) error
		io.Closer
	}

	UserInterface interface {
		RegisterUser(ctx context.Context, name, email, passwordHash, token string) error
		Authenticate(ctx context.Context, token string) (User, error)
		Login(ctx context.Context, email, password string) error
		UpdateToken(ctx context.Context, email, token string) error

		io.Closer
	}

	ListInterface interface {
		Create(context.Context, List) error
		GetAll(context.Context, UserID) ([]List, error)
		Update(context.Context, List) error
		Delete(context.Context, UserID, ListID) error

		io.Closer
	}

	TaskInterface interface {
		Create(context.Context, UserID, Task) error
		Update(context.Context, UserID, Task) error
		Delete(context.Context, UserID, TaskID) error

		io.Closer
	}
)

type Priority = string

const (
	Low    Priority = "low"
	Normal Priority = "normal"
	High   Priority = "high"
)
