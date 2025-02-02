package domain

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
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
		ID        ListID
		UserID    UserID
		Name      string
		UpdatedAT time.Time
	}

	TaskID uuid.UUID

	Task struct {
		ID        TaskID
		ListID    ListID
		Priority  Priority
		Deadline  time.Time
		Done      bool
		Name      string
		UpdatedAT time.Time
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
		CreateList(ctx context.Context, listID, userID uuid.UUID, name string) error
		ReadAll(ctx context.Context, d uuid.UUID) error
		UpdateName(ctx context.Context, listID uuid.UUID, name string) error
		DeleteList(ctx context.Context, listID uuid.UUID) error

		io.Closer
	}

	TaskInterface interface {
		CreateTask(ctx context.Context, taskID, userID uuid.UUID, name string) error

		io.Closer
	}
)

type Priority string

const (
	Low    Priority = "low"
	Normal Priority = "normal"
	High   Priority = "high"
)

func (t *Priority) Scan(src any) error {
	v, ok := src.(string)
	if !ok {
		return errors.New("database value is not string")
	}

	*t = Priority(v)

	return nil
}

// Value implements driver.Valuer.
func (t *Priority) Value() (driver.Value, error) {
	return string(*t), nil
}

var _ driver.Valuer = (*Priority)(nil)
var _ sql.Scanner = (*Priority)(nil)
