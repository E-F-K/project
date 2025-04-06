package controller

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"todo_list/internal/adapter/logger"
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
	ctx, curUser := c.Request.Context(), getCurrentUser(c)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	var task domain.Task
	if err = json.Unmarshal(body, &task); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	if err = ctl.service.Create(ctx, curUser.ID, task); err != nil {
		slog.ErrorContext(ctx, "Create task failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Create task failed."))

		return
	}

	c.Status(http.StatusCreated)
}

func (ctl *Tasks) UpdateTask(c *gin.Context) {
	ctx, curUser := c.Request.Context(), getCurrentUser(c)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	var task domain.Task
	if err = json.Unmarshal(body, &task); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	if err = ctl.service.Update(ctx, curUser.ID, task); err != nil {
		slog.ErrorContext(ctx, "Update task failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Update task failed."))

		return
	}

	c.Status(http.StatusNoContent)
}

/*func (ctl *Tasks) UpdateTaskStatus(c *gin.Context) {
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
}*/

func (ctl *Tasks) DeleteTask(c *gin.Context) {
	ctx, curUser := c.Request.Context(), getCurrentUser(c)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	var message struct {
		TaskID domain.TaskID `json:"id"`
	}
	if err = json.Unmarshal(body, &message); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	if err = ctl.service.Delete(ctx, curUser.ID, message.TaskID); err != nil {
		slog.ErrorContext(ctx, "Delete task failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Delete task failed."))

		return
	}

	c.Status(http.StatusNoContent)
}

func (ctl *Tasks) Close() error {
	return ctl.service.Close()
}
