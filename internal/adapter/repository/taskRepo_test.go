package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"todo_list/internal/domain"

	"todo_list/internal/adapter/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUsers_GetAllTasks(t *testing.T) {
	repoTask := repository.NewTasks()

	ctx := context.Background()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)

		list := fixtureCreateList(t, ctx, connection, user.ID)

		_ = fixtureCreateTask(t, ctx, connection, list.ID, "firstTask")
		_ = fixtureCreateTask(t, ctx, connection, list.ID, "secondTask")

		tasks, err := repoTask.GetAllTasks(ctx, connection, []domain.ListID{list.ID})
		require.NoError(t, err)
		require.Equal(t, 2, len(tasks))

		task := fixtureCreateTask(t, ctx, connection, list.ID, "thirdTask")

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

func fixtureCreateTask(t *testing.T, ctx context.Context, connection domain.Connection, listID domain.ListID, name string) domain.Task {
	task := domain.Task{
		ID:        domain.TaskID(uuid.New()),
		ListID:    listID,
		Priority:  "low",
		Deadline:  time.Now(),
		Done:      false,
		Name:      name,
		UpdatedAT: time.Now(),
	}
	require.NoError(t, repository.NewTasks().Create(ctx, connection, task))

	return task
}
