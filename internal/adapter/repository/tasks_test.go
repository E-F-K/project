package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"todo_list/internal/domain"
	dbMocks "todo_list/mocks/todo_list/src/domain"

	"todo_list/internal/adapter/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTasksIntegration(t *testing.T) {
	repoTask := repository.NewTasks()

	ctx := context.Background()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)

		list := fixtureCreateList(t, ctx, connection, user.ID)

		_ = fixtureCreateTask(t, ctx, connection, user.ID, list.ID, "firstTask")
		_ = fixtureCreateTask(t, ctx, connection, user.ID, list.ID, "secondTask")

		tasks, err := repoTask.GetAllTasks(ctx, connection, user.ID, []domain.ListID{list.ID})
		require.NoError(t, err)
		require.Equal(t, 2, len(tasks))

		task := fixtureCreateTask(t, ctx, connection, user.ID, list.ID, "thirdTask")

		task.Name = "new task name"
		require.NoError(t, repoTask.Update(ctx, connection, user.ID, task))

		newTask, err := repoTask.Read(ctx, connection, user.ID, task.ID)
		require.NoError(t, err)
		require.Equal(t, task.Name, newTask.Name)

		require.NoError(t, repoTask.Delete(ctx, connection, user.ID, task.ID))

		_, err = repoTask.Read(ctx, connection, user.ID, task.ID)
		require.ErrorIs(t, err, sql.ErrNoRows)

		return nil
	})
}

func TestTasksIntegrationInvalidUserIDCreate(t *testing.T) {
	repoTask := repository.NewTasks()

	ctx := context.Background()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)
		task := domain.Task{
			ID:     uuid.New(),
			ListID: uuid.New(), // <- non existing list ID
			Name:   "Some name",
		}

		err := repoTask.Create(ctx, connection, user.ID, task)
		require.Error(t, err)
		require.ErrorContains(t, err, "not found or access denied")

		return nil
	})
}

func TestTasksIntegrationInvalidUserIDDelete(t *testing.T) {
	repoTask := repository.NewTasks()

	ctx := context.Background()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)
		list := fixtureCreateList(t, ctx, connection, user.ID)
		task := fixtureCreateTask(t, ctx, connection, user.ID, list.ID, "firstTask")

		wrongUserID := uuid.New()

		err := repoTask.Delete(ctx, connection, wrongUserID, task.ID)
		require.Error(t, err)
		require.ErrorContains(t, err, "not found or access denied")

		return nil
	})
}

func TestTasksIntegrationInvalidUserIDRead(t *testing.T) {
	repoTask := repository.NewTasks()

	ctx := context.Background()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)
		list := fixtureCreateList(t, ctx, connection, user.ID)
		task := fixtureCreateTask(t, ctx, connection, user.ID, list.ID, "firstTask")

		wrongUserID := uuid.New()

		_, err := repoTask.Read(ctx, connection, wrongUserID, task.ID)
		require.Error(t, err)
		require.ErrorContains(t, err, "not found or access denied")

		return nil
	})
}

func TestTasksIntegrationInvalidUserIDUpdate(t *testing.T) {
	repoTask := repository.NewTasks()

	ctx := context.Background()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)
		list := fixtureCreateList(t, ctx, connection, user.ID)
		task := fixtureCreateTask(t, ctx, connection, user.ID, list.ID, "firstTask")

		wrongUserID := uuid.New()

		err := repoTask.Update(ctx, connection, wrongUserID, task)
		require.Error(t, err)
		require.ErrorContains(t, err, "not found or access denied")

		return nil
	})
}

func TestTasksIntegrationInvalidUserIDGetAllTasks(t *testing.T) {
	repoTask := repository.NewTasks()

	ctx := context.Background()
	provider := cleanTablesAndCreateProvider(ctx, t)
	defer func() { _ = provider.Close() }()

	provider.ExecuteTx(ctx, func(ctx context.Context, connection domain.Connection) error {
		user := fixtureCreateUser(t, ctx, connection)
		list1 := fixtureCreateList(t, ctx, connection, user.ID)
		list2 := fixtureCreateList(t, ctx, connection, user.ID)

		_ = fixtureCreateTask(t, ctx, connection, user.ID, list1.ID, "firstTask")
		_ = fixtureCreateTask(t, ctx, connection, user.ID, list1.ID, "secondTask")

		wrongUserID := uuid.New()

		listIDs := []domain.ListID{list1.ID, list2.ID}

		_, err := repoTask.GetAllTasks(ctx, connection, wrongUserID, listIDs)
		require.Error(t, err)
		require.ErrorContains(t, err, "not found or access denied")

		return nil
	})
}

func TestTasksUnit(t *testing.T) {
	now := time.Now().Add(+time.Hour)
	validEmptyTask := domain.Task{
		ID:       domain.TaskID(uuid.New()),
		ListID:   domain.ListID(uuid.New()),
		Name:     "Some task name",
		Priority: "low",
		Deadline: &now,
		Done:     false,
	}
	userID := domain.UserID(uuid.New())
	ctx := context.Background()

	mockListExistsCall := func(connection *dbMocks.MockConnection, userID domain.UserID, listID domain.ListID, err error) {
		connection.EXPECT().
			GetContext(mock.Anything, mock.Anything, mock.Anything, userID, listID).
			Return(err).
			Once()
	}
	mockListExists := func(connection *dbMocks.MockConnection, userID domain.UserID, listID domain.ListID) {
		mockListExistsCall(connection, userID, listID, nil)
	}

	tests := []struct {
		name  string
		check func(*testing.T, *repository.Tasks, *dbMocks.MockConnection)
	}{
		{
			name: "Create Task DB Error on List Exists",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				mockListExistsCall(connection, userID, validEmptyTask.ListID, errors.New("some db error"))

				err := repo.Create(ctx, connection, userID, validEmptyTask)

				require.ErrorIs(t, err, repository.ErrTasksCreate)
				require.ErrorContains(t, err, "some db error")
			},
		},
		{
			name: "Create DB Error",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				mockListExists(connection, userID, validEmptyTask.ListID)
				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(0, errors.New("some error")).
					Once()

				err := repo.Create(ctx, connection, userID, validEmptyTask)

				require.ErrorIs(t, err, repository.ErrTasksCreate)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Delete Task DB Error on List Exists 1",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, "select list_id from tasks where id = $1", validEmptyTask.ID).
					Return(errors.New("empty rows, list not found")).
					Once()

				err := repo.Delete(ctx, connection, userID, validEmptyTask.ID)

				require.ErrorIs(t, err, repository.ErrTasksDelete)
				require.ErrorContains(t, err, "empty rows, list not found")
			},
		},
		{
			name: "Delete Task DB Error on List Exists 2",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, "select list_id from tasks where id = $1", validEmptyTask.ID).
					Run(func(_ context.Context, shouldBeListID any, _ string, _ ...any) {
						p := shouldBeListID.(*uuid.UUID)
						*p = validEmptyTask.ListID
					}).
					Return(nil).
					Once()

				mockListExistsCall(connection, userID, validEmptyTask.ListID, errors.New("some db error"))

				err := repo.Delete(ctx, connection, userID, validEmptyTask.ID)

				require.ErrorIs(t, err, repository.ErrTasksDelete)
				require.ErrorContains(t, err, "some db error")
			},
		},
		{
			name: "Delete DB Error",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, "select list_id from tasks where id = $1", validEmptyTask.ID).
					Run(func(_ context.Context, shouldBeListID any, _ string, _ ...any) {
						p := shouldBeListID.(*uuid.UUID)
						*p = validEmptyTask.ListID
					}).
					Return(nil).
					Once()

				mockListExists(connection, userID, validEmptyTask.ListID)

				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, validEmptyTask.ID).
					Return(0, errors.New("some error")).
					Once()

				err := repo.Delete(ctx, connection, userID, validEmptyTask.ID)

				require.ErrorIs(t, err, repository.ErrTasksDelete)
				require.ErrorContains(t, err, "some error")
				require.ErrorIs(t, err, repository.ErrTasksDelete)
			},
		},
		{
			name: "Read Task DB Error on List Exists 1",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, "select list_id from tasks where id = $1", validEmptyTask.ID).
					Return(errors.New("empty rows, list not found")).
					Once()

				_, err := repo.Read(ctx, connection, userID, validEmptyTask.ID)

				require.ErrorIs(t, err, repository.ErrTasksRead)
				require.ErrorContains(t, err, "empty rows, list not found")
			},
		},
		{
			name: "Read Task DB Error on List Exists 2",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, "select list_id from tasks where id = $1", validEmptyTask.ID).
					Run(func(_ context.Context, shouldBeListID any, _ string, _ ...any) {
						p := shouldBeListID.(*uuid.UUID)
						*p = validEmptyTask.ListID
					}).
					Return(nil).
					Once()

				mockListExistsCall(connection, userID, validEmptyTask.ListID, errors.New("some db error"))

				_, err := repo.Read(ctx, connection, userID, validEmptyTask.ID)

				require.ErrorIs(t, err, repository.ErrTasksRead)
				require.ErrorContains(t, err, "some db error")
			},
		},
		{
			name: "Read DB Error",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, "select list_id from tasks where id = $1", validEmptyTask.ID).
					Run(func(_ context.Context, dest any, _ string, _ ...any) {
						p := dest.(*uuid.UUID)
						*p = uuid.UUID(validEmptyTask.ListID)
					}).
					Return(nil).
					Once()

				mockListExists(connection, userID, validEmptyTask.ListID)

				connection.EXPECT().
					GetContext(mock.Anything, mock.Anything, `select id, list_id, priority, deadline, done, name, updated_at from tasks where id = $1`, validEmptyTask.ID).
					Return(errors.New("some error")).
					Once()

				_, err := repo.Read(ctx, connection, userID, validEmptyTask.ID)

				require.ErrorIs(t, err, repository.ErrTasksRead)
				require.ErrorContains(t, err, "some error")
			},
		},
		{
			name: "Update Task DB Error on List Exists",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				mockListExistsCall(connection, userID, validEmptyTask.ListID, errors.New("some db error"))

				err := repo.Update(ctx, connection, userID, validEmptyTask)

				require.ErrorIs(t, err, repository.ErrTasksUpdate)
				require.ErrorContains(t, err, "some db error")
			},
		},
		{
			name: "Update DB Error",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				mockListExists(connection, userID, validEmptyTask.ListID)

				connection.EXPECT().
					ExecContext(mock.Anything, mock.Anything, validEmptyTask.ID, validEmptyTask.Name, domain.Priority(validEmptyTask.Priority), validEmptyTask.Deadline, validEmptyTask.Done).
					Return(0, errors.New("update error")).
					Once()

				err := repo.Update(ctx, connection, userID, validEmptyTask)

				require.ErrorIs(t, err, repository.ErrTasksUpdate)
				require.ErrorContains(t, err, "update error")
			},
		},
		{
			name: "GetAllTasks listExists error",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				ids := []domain.ListID{validEmptyTask.ListID, domain.ListID(uuid.New())}

				mockListExists(connection, userID, ids[0])
				mockListExistsCall(connection, userID, ids[1], errors.New("access check failed"))

				_, err := repo.GetAllTasks(ctx, connection, userID, ids)

				require.ErrorIs(t, err, repository.ErrTasksGetAllTasks)
				require.ErrorContains(t, err, "access check failed")
			},
		},
		{
			name: "GetAllTasks DB error",
			check: func(t *testing.T, repo *repository.Tasks, connection *dbMocks.MockConnection) {
				ids := []domain.ListID{validEmptyTask.ListID}
				mockListExists(connection, userID, validEmptyTask.ListID)

				connection.EXPECT().
					SelectContext(mock.Anything, mock.Anything, mock.Anything, ids).
					Return(errors.New("select error")).
					Once()

				_, err := repo.GetAllTasks(ctx, connection, userID, ids)

				require.ErrorIs(t, err, repository.ErrTasksGetAllTasks)
				require.ErrorContains(t, err, "select error")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(t, repository.NewTasks(), dbMocks.NewMockConnection(t))
		})
	}
}

func fixtureCreateTask(t *testing.T, ctx context.Context, connection domain.Connection, userID domain.UserID, listID domain.ListID, name string) domain.Task {
	now := time.Now()
	task := domain.Task{
		ID:       domain.TaskID(uuid.New()),
		ListID:   listID,
		Priority: "low",
		Deadline: &now,
		Done:     false,
		Name:     name,
	}
	require.NoError(t, repository.NewTasks().Create(ctx, connection, userID, task))

	return task
}
