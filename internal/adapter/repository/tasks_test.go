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

func TestBasicTasksOperations(t *testing.T) {
	ctx := context.Background()

	repoTask := repository.NewTasks()
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

		listuuid := uuid.New()
		list := domain.List{
			ID:        domain.ListID(listuuid),
			UserID:    domain.UserID(myuuid),
			Name:      "task name",
			Email:     "user@email.foo",
			UpdatedAT: domain.Timestamp(time.Now()),
		}
		require.NoError(t, repoList.Create(ctx, connection, list))

		task := domain.Task{
			ID:        domain.TaskID(uuid.New()),
			ListID:    domain.ListID(listuuid),
			Priority:  domain.Low,
			Deadline:  domain.Timestamp(time.Now()),
			Done:      false,
			Name:      "task name",
			UpdatedAT: domain.Timestamp(time.Now()),
		}
		require.NoError(t, repoTask.Create(ctx, connection, task))

		task.Name = "new task name"
		require.NoError(t, repoTask.Update(ctx, connection, task))

		newTask, err := repoTask.Read(ctx, connection, task.ID)
		require.NoError(t, err)
		require.Equal(t, task.Name, newTask.Name)

		require.NoError(t, repoTask.Delete(ctx, connection, task.ID))

		_, err = repoTask.Read(ctx, connection, task.ID)
		require.ErrorIs(t, err, sql.ErrNoRows)

		return nil
	})
}
