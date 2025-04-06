package controller

import (
	"io"
	"net/http"

	"todo_list/internal/domain"

	"github.com/gin-gonic/gin"
)

var _ io.Closer = (*Tasks)(nil)

type Tasks struct {
	service domain.TaskInterface
}

func NewTasks(service domain.TaskInterface) *Tasks {
	return &Tasks{service: service}
}

func (ctl *Tasks) CreateTask(c *gin.Context) {
	curUser := getCurrentUser(c)

	c.JSON(http.StatusOK, curUser)
}

func (ctl *Tasks) UpdateTaskStatus(c *gin.Context) {
	curUser := getCurrentUser(c)

	c.JSON(http.StatusOK, curUser)
}

func (ctl *Tasks) UpdateTaskPriority(c *gin.Context) {
	curUser := getCurrentUser(c)

	c.JSON(http.StatusOK, curUser)
}

func (ctl *Tasks) UpdateTaskName(c *gin.Context) {
	curUser := getCurrentUser(c)

	c.JSON(http.StatusOK, curUser)
}

func (ctl *Tasks) UpdateTaskDeadline(c *gin.Context) {
	curUser := getCurrentUser(c)

	c.JSON(http.StatusOK, curUser)
}

func (ctl *Tasks) DeleteTask(c *gin.Context) {
	curUser := getCurrentUser(c)

	c.JSON(http.StatusOK, curUser)
}

func (ctl *Tasks) Close() error {
	return ctl.service.Close()
}
