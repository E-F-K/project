package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"todo_list/internal/adapter/repository"
	"todo_list/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBasicListsOperations(t *testing.T) {
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

		require.NoError(t, repoList.Delete(ctx, connection, list.ID))

		_, err = repoList.Read(ctx, connection, list.ID)
		require.ErrorIs(t, err, sql.ErrNoRows)

		return nil
	})
}

func fixtureCreateList(t *testing.T, ctx context.Context, connection domain.Connection, userID domain.UserID) domain.List {
	list := domain.List{
		ID:        domain.ListID(uuid.New()),
		UserID:    userID,
		Name:      "list name",
		Email:     "user@email.foo",
		UpdatedAT: time.Now(),
	}
	require.NoError(t, repository.NewLists().Create(ctx, connection, list))

	return list
}
