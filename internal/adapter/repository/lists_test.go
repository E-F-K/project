package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"todo_list/internal/adapter/repository"
	"todo_list/internal/domain"
	dbMocks "todo_list/mocks/todo_list/src/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListsIntegration(t *testing.T) {
	ctx := context.Background()
	repoList := repository.NewLists()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)

		list := fixtureCreateList(t, ctx, connection, user.ID)

		list.Name = "new list name"
		require.NoError(t, repoList.Update(ctx, connection, list))

		newList, err := repoList.Read(ctx, connection, list.ID)
		require.NoError(t, err)
		require.Equal(t, list.Name, newList.Name)

		allListsBeforeDelete, err := repoList.ReadAll(ctx, connection, user.ID)
		require.NoError(t, err)

		require.NoError(t, repoList.Delete(ctx, connection, user.ID, list.ID))

		allListsAfterDelete, err := repoList.ReadAll(ctx, connection, user.ID)
		require.NoError(t, err)
		require.Equal(t, len(allListsBeforeDelete)-1, len(allListsAfterDelete))

		_, err = repoList.Read(ctx, connection, list.ID)
		require.ErrorIs(t, err, sql.ErrNoRows)

		return nil
	})
}

func TestListsUnit(t *testing.T) {
	validEmptyList := domain.List{
		ID:        domain.ListID(uuid.New()),
		UserID:    domain.UserID(uuid.New()),
		Name:      "Some list name",
		UpdatedAt: time.Now().Add(-time.Hour),
		Tasks:     nil,
	}
	ctx := context.Background()

	tests := []struct {
		name  string
		check func(*testing.T, *repository.Lists, *dbMocks.MockConnection)
	}{
		{
			name: "Create DB Error",
			check: func(t *testing.T, repo *repository.Lists, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, validEmptyList.ID, validEmptyList.UserID, validEmptyList.Name).
					Return(0, errors.New("some error")).
					Once()

				err := repo.Create(ctx, connection, validEmptyList)

				require.ErrorIs(t, err, repository.ErrListsCreate)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Delete DB Error",
			check: func(t *testing.T, repo *repository.Lists, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, validEmptyList.UserID, validEmptyList.ID).
					Return(0, errors.New("some error")).
					Once()

				err := repo.Delete(ctx, connection, validEmptyList.UserID, validEmptyList.ID)

				require.ErrorIs(t, err, repository.ErrListsDelete)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Read DB Error",
			check: func(t *testing.T, repo *repository.Lists, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, mock.Anything, validEmptyList.ID).
					Return(errors.New("some error")).
					Once()

				_, err := repo.Read(ctx, connection, validEmptyList.ID)

				require.ErrorIs(t, err, repository.ErrListsRead)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Invalid Update Input",
			check: func(t *testing.T, repo *repository.Lists, connection *dbMocks.MockConnection) {
				listWithTasks := validEmptyList
				listWithTasks.Tasks = append(listWithTasks.Tasks, domain.Task{})

				err := repo.Update(ctx, connection, listWithTasks)

				require.ErrorIs(t, err, repository.ErrListsUpdate)
				require.ErrorContains(t, err, "task updates are not supported")
			},
		},
		{
			name: "Update DB Error",
			check: func(t *testing.T, repo *repository.Lists, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, validEmptyList.UserID, validEmptyList.ID, validEmptyList.Name).
					Return(0, errors.New("some error")).
					Once()

				err := repo.Update(ctx, connection, validEmptyList)

				require.ErrorIs(t, err, repository.ErrListsUpdate)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Read All DB Error",
			check: func(t *testing.T, repo *repository.Lists, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					SelectContext(mock.Anything, mock.Anything, mock.Anything, validEmptyList.UserID).
					Return(errors.New("some error")).
					Once()

				_, err := repo.ReadAll(ctx, connection, validEmptyList.UserID)

				require.ErrorIs(t, err, repository.ErrListsReadAll)
				require.ErrorContains(t, err, "some error")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(t, repository.NewLists(), dbMocks.NewMockConnection(t))
		})
	}
}

func fixtureCreateList(t *testing.T, ctx context.Context, connection domain.Connection, userID domain.UserID) domain.List {
	list := domain.List{
		ID:        domain.ListID(uuid.New()),
		UserID:    userID,
		Name:      "list name",
		UpdatedAt: time.Now(),
	}
	require.NoError(t, repository.NewLists().Create(ctx, connection, list))

	return list
}
