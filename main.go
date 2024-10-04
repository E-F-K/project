package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func provideConnection(ctx context.Context, receiver func(context.Context, *pgx.Conn) error) error {
	conn, err := pgx.Connect(ctx, "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	return receiver(ctx, conn)
}

func printUserName(ctx context.Context, userID string) {
	var name string
	err := provideConnection(ctx, func(ctx context.Context, conn *pgx.Conn) error {
		return conn.QueryRow(context.Background(), "select name from users where id=$1", userID).Scan(&name)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(name)
}

func printTaskByID(ctx context.Context, taskID string) {
	var task string
	err := provideConnection(ctx, func(ctx context.Context, conn *pgx.Conn) error {
		return conn.QueryRow(context.Background(), "select name from tasks where id=$1", taskID).Scan(&task)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(task)
}

func printALLTasksFromList(ctx context.Context, listID string) {
	var tasks string
	err := provideConnection(ctx, func(ctx context.Context, conn *pgx.Conn) error {
		return conn.QueryRow(context.Background(), "select name from tasks where list_id=$1", listID).Scan(&tasks)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(tasks)
}

func main() {
	// urlExample := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var inputUserID string
	inputUserID = "00000000-0000-0000-0000-000000000001"
	printUserName(context.Background(), inputUserID)

	var inputTaskID string
	inputTaskID = "20000000-0000-0000-0000-000000000002"
	printTaskByID(context.Background(), inputTaskID)

	var inputListID string
	inputListID = "10000000-0000-0000-0000-000000000001"
	printALLTasksFromList(context.Background(), inputListID)

}

/*
const (
	value1 = 1
	value2 = 2
	value3 = 3
)

func main() {
	fmt.Println(value1)
	fmt.Println(value2)
	fmt.Println(value3)
}
*/
