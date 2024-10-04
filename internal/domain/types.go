package domain

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"

	"github.com/google/uuid"
)

type (
	UserID uuid.UUID

	User struct {
		ID    UserID
		Name  string
		Email string
		Token string
	}

	ListID    uuid.UUID
	Timestamp time.Time

	List struct {
		ID        ListID
		UserID    UserID
		Name      string
		Email     string
		UpdatedAT Timestamp
	}

	TaskID uuid.UUID

	Task struct {
		ID        TaskID
		ListID    ListID
		Priority  Priority
		Deadline  Timestamp
		Done      bool
		Name      string
		UpdatedAT Timestamp
	}

	Connection interface {
		GetContext(context.Context, any, string, ...any) error
		SelectContext(context.Context, any, string, ...any) error
		ExecContext(context.Context, string, ...any) (int64, error)
	}
	ConnectionProvider interface {
		Execute(context.Context, func(context.Context, Connection) error) error
		ExecuteTx(context.Context, func(context.Context, Connection) error) error
	}
)

// Scan implements sql.Scanner.
func (t *Timestamp) Scan(src any) error {
	v, ok := src.(time.Time)
	if !ok {
		return errors.New("database value is not time.Time")
	}

	*t = Timestamp(v)

	return nil
}

// Value implements driver.Valuer.
func (t *Timestamp) Value() (driver.Value, error) {
	return time.Time(*t), nil
}

var _ driver.Valuer = (*Timestamp)(nil)
var _ sql.Scanner = (*Timestamp)(nil)

type Priority int

const (
	Low Priority = iota + 1
	Medium
	High
)

func (t *Priority) Scan(src any) error {
	v, ok := src.(int)
	if !ok {
		return errors.New("database value is not int")
	}

	*t = Priority(v)

	return nil
}

// Value implements driver.Valuer.
func (t *Priority) Value() (driver.Value, error) {
	return int(*t), nil
}

var _ driver.Valuer = (*Priority)(nil)
var _ sql.Scanner = (*Priority)(nil)
