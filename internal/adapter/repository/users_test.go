package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"todo_list/internal/adapter/repository"
	"todo_list/internal/domain"
	dbMocks "todo_list/mocks/todo_list/src/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsersIntegration(t *testing.T) {
	ctx := context.Background()

	repo := repository.NewUsers()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)

		user.Name = "new user name"
		require.NoError(t, repo.Update(ctx, connection, user))

		updatedUser, err := repo.ReadByToken(ctx, connection, user.Token)
		require.NoError(t, err)
		require.Equal(t, user.Name, updatedUser.Name)

		updatedUser, err = repo.ReadByEmail(ctx, connection, user.Email)
		require.NoError(t, err)
		require.Equal(t, user.Name, updatedUser.Name)

		newToken := "some new token"
		require.NoError(t, repo.UpdateTokenByEmail(ctx, connection, user.Email, newToken))
		updatedUser, err = repo.ReadByEmail(ctx, connection, user.Email)
		require.NoError(t, err)
		require.Equal(t, newToken, updatedUser.Token)

		require.NoError(t, repo.Delete(ctx, connection, user.ID))

		_, err = repo.ReadByToken(ctx, connection, user.Token)
		require.ErrorIs(t, err, sql.ErrNoRows)

		return nil
	})
}

func TestUsersUnit(t *testing.T) {
	validEmptyUser := domain.User{
		ID:           domain.UserID(uuid.New()),
		Name:         "Some user name",
		Email:        "some@email.foo",
		PasswordHash: "password hash",
		Token:        "some token",
	}
	ctx := context.Background()

	tests := []struct {
		name  string
		check func(*testing.T, *repository.Users, *dbMocks.MockConnection)
	}{
		{
			name: "Create DB Error",
			check: func(t *testing.T, repo *repository.Users, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, validEmptyUser.ID, validEmptyUser.Name, validEmptyUser.Email, validEmptyUser.PasswordHash, validEmptyUser.Token).
					Return(0, errors.New("some error")).
					Once()

				err := repo.Create(ctx, connection, validEmptyUser)

				require.ErrorIs(t, err, repository.ErrUsersCreate)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Delete DB Error",
			check: func(t *testing.T, repo *repository.Users, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, validEmptyUser.ID).
					Return(0, errors.New("some error")).
					Once()

				err := repo.Delete(ctx, connection, validEmptyUser.ID)

				require.ErrorIs(t, err, repository.ErrUsersDelete)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Read by Token DB Error",
			check: func(t *testing.T, repo *repository.Users, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, mock.Anything, validEmptyUser.Token).
					Return(errors.New("some error")).
					Once()

				_, err := repo.ReadByToken(ctx, connection, validEmptyUser.Token)

				require.ErrorIs(t, err, repository.ErrUsersRead)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Read by Email DB Error",
			check: func(t *testing.T, repo *repository.Users, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, mock.Anything, validEmptyUser.Email).
					Return(errors.New("some error")).
					Once()

				_, err := repo.ReadByEmail(ctx, connection, validEmptyUser.Email)

				require.ErrorIs(t, err, repository.ErrUsersRead)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Update Token by Email DB Error",
			check: func(t *testing.T, repo *repository.Users, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, validEmptyUser.ID, validEmptyUser.Name, validEmptyUser.Email, validEmptyUser.PasswordHash, validEmptyUser.Token).
					Return(0, errors.New("some error")).
					Once()

				err := repo.Update(ctx, connection, validEmptyUser)

				require.ErrorIs(t, err, repository.ErrUsersUpdate)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Update Token by Email DB Error",
			check: func(t *testing.T, repo *repository.Users, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, validEmptyUser.Email, validEmptyUser.Token).
					Return(0, errors.New("some error")).
					Once()

				err := repo.UpdateTokenByEmail(ctx, connection, validEmptyUser.Email, validEmptyUser.Token)

				require.ErrorIs(t, err, repository.ErrUsersUpdateTokenByEmail)
				require.ErrorContains(t, err, "some error")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(t, repository.NewUsers(), dbMocks.NewMockConnection(t))
		})
	}

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
