package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"todo_list/internal/adapter/database"
	"todo_list/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func cleanTablesAndCreateProvider(ctx context.Context, t *testing.T) domain.ConnectionProvider {
	godotenv.Load("../../../.env")

	tablesToClean := []string{"users", "lists", "tasks"}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_CONNECTION"))
	require.NoError(t, err)

	provider := database.NewPostgresProvider(pool)

	err = provider.Execute(ctx, func(ctx context.Context, connection domain.Connection) error {
		for _, table := range tablesToClean {
			_, err := connection.ExecContext(ctx, fmt.Sprintf("truncate %s cascade", table))
			require.NoError(t, err)
		}

		return nil
	})
	require.NoError(t, err)

	return provider
}
