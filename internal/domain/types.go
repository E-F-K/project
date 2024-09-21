package domain

import (
	"context"

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
