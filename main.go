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

	user, toDo, authMiddleware, err := createControllers()
	if err != nil {
		slog.ErrorContext(ctx, "Create application service failed.", logger.ErrAttr(err))
		os.Exit(1)
	}
	defer func() { _ = user.Close() }()
	defer func() { _ = toDo.Close() }()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("CORS_ALLOWED_ORIGIN")},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           time.Minute,
	}))

	router.POST("register", user.Register)
	router.POST("login", user.Login)

	authRequired := router.Group("/v1")
	authRequired.Use(authMiddleware)
	{
		authRequired.GET("lists", toDo.GetUserLists)
	}

	router.Run(os.Getenv("SERVER_ADDRESS"))
}

func createControllers() (*controller.Users, *controller.ToDo, gin.HandlerFunc, error) {
	pool, err := pgxpool.New(context.Background(), domain.ConnectionString)
	if err != nil {
		return nil, nil, nil, errors.Join(errors.New("create database pool failed"), err)
	}

	appService := domain.NewToDoService(database.NewPostgresProvider(pool), repository.NewUsers())

	authMiddlware := controller.NewAuthMiddleware(appService).Auth

	return controller.NewUsers(appService), controller.NewToDo(), authMiddlware, nil
}
