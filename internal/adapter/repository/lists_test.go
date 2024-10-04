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
	repo := repository.NewUsers()
	cleanTablesAndCreateProvider(ctx, t).ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		myuuid := uuid.New()
		user := domain.User{
			ID:    domain.UserID(myuuid),
			Name:  "user name",
			Email: "user@email.foo",
			Token: "some secret token",
		}
		require.NoError(t, repo.Create(ctx, connection, user))

		list := domain.List{
			ID:        domain.ListID(uuid.New()),
			UserID:    domain.UserID(myuuid),
			Name:      "list name",
			Email:     "user@email.foo",
			UpdatedAT: domain.Timestamp(time.Now()),
		}
		require.NoError(t, repoList.Create(ctx, connection, list))

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
