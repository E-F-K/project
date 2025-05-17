package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"todo_list/internal/adapter/controller"
	"todo_list/internal/adapter/database"
	"todo_list/internal/adapter/logger"
	"todo_list/internal/adapter/repository"
	"todo_list/internal/domain"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	users, lists, tasks, authMiddleware, err := createControllers()
	if err != nil {
		slog.ErrorContext(ctx, "Create application service failed.", logger.ErrAttr(err))
		os.Exit(1)
	}
	defer func() { _ = users.Close() }()
	defer func() { _ = lists.Close() }()
	defer func() { _ = tasks.Close() }()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("CORS_ALLOWED_ORIGIN")},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           time.Minute,
	}))

	router.POST("register", users.Register)
	router.POST("login", users.Login)

	authRequired := router.Group("/v1")
	authRequired.Use(authMiddleware)
	{
		authRequired.GET("list", lists.GetUserListsAndTasks)
		authRequired.POST("list", lists.CreateList)
		authRequired.PUT("list", lists.UpdateList)
		authRequired.DELETE("list", lists.DeleteList)

		authRequired.POST("task", tasks.CreateTask)
		authRequired.PUT("task", tasks.UpdateTask)
		authRequired.DELETE("task", tasks.DeleteTask)
	}

	router.Run(os.Getenv("SERVER_ADDRESS"))
}

func createControllers() (*controller.Users, *controller.Lists, *controller.Tasks, gin.HandlerFunc, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_CONNECTION"))
	if err != nil {
		return nil, nil, nil, nil, errors.Join(errors.New("create database pool failed"), err)
	}

	provider := database.NewPostgresProvider(pool)
	userService := domain.NewUserService(provider, repository.NewUsers())
	listService := domain.NewListService(provider, repository.NewLists())
	taskService := domain.NewTaskService(provider, repository.NewTasks())
	authMiddlware := controller.NewAuthMiddleware(userService).Auth

	return controller.NewUsers(userService), controller.NewLists(listService), controller.NewTasks(taskService), authMiddlware, nil
}
