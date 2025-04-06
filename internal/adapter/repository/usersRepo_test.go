package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"todo_list/internal/adapter/repository"
	"todo_list/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBasicUserOperations(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewUsers()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)

		user.Name = "new user name"
		require.NoError(t, repo.Update(ctx, connection, user))

		newUser, err := repo.ReadByToken(ctx, connection, user.Token)
		require.NoError(t, err)
		require.Equal(t, user.Name, newUser.Name)

		require.NoError(t, repo.Delete(ctx, connection, user.ID))

		_, err = repo.ReadByToken(ctx, connection, user.Token)
		require.ErrorIs(t, err, sql.ErrNoRows)

		return nil
	})
}

func fixtureCreateUser(t *testing.T, ctx context.Context, connection domain.Connection) domain.User {
	user := domain.User{
		ID:           domain.UserID(uuid.New()),
		Name:         "user name",
		Email:        "user@email.foo",
		PasswordHash: "some password hash",
		Token:        "some secret token",
	}
	require.NoError(t, repository.NewUsers().Create(ctx, connection, user))

	return user
}
