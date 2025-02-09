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

var _ io.Closer = (*Lists)(nil)

type Lists struct {
	service domain.ListInterface
}

func NewLists(service domain.ListInterface) *Lists {
	return &Lists{service: service}
}

func (ctl *Lists) GetUserListsAndTasks(c *gin.Context) {
	ctx, curUser := c.Request.Context(), getCurrentUser(c)

	listAndTasks, err := ctl.service.ReadAll(ctx, curUser.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Read all failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read all failed."))

		return
	}

	c.JSON(http.StatusOK, listAndTasks)
}

func (ctl *Lists) CreateList(c *gin.Context) {
	ctx, curUser := c.Request.Context(), getCurrentUser(c)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	var list domain.List
	if err = json.Unmarshal(body, &list); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	list.UserID = curUser.ID

	if err = ctl.service.Create(ctx, list); err != nil {
		slog.ErrorContext(ctx, "Create list failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Create list failed."))

		return
	}

	c.Status(http.StatusCreated)
}

func (ctl *Lists) UpdateList(c *gin.Context) {
	ctx, curUser := c.Request.Context(), getCurrentUser(c)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	var list domain.List
	if err = json.Unmarshal(body, &list); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	list.UserID = curUser.ID

	if err = ctl.service.Update(ctx, list); err != nil {
		slog.ErrorContext(ctx, "Update name failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Update name failed."))

		return
	}

	c.Status(http.StatusNoContent)
}

func (ctl *Lists) DeleteList(c *gin.Context) {
	ctx, curUser := c.Request.Context(), getCurrentUser(c)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	var message struct {
		ListID domain.ListID `json:"id"`
	}
	if err = json.Unmarshal(body, &message); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	if err = ctl.service.Delete(ctx, curUser.ID, message.ListID); err != nil {
		slog.ErrorContext(ctx, "Delete list failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Delete list failed."))

		return
	}

	c.Status(http.StatusNoContent)
}

func (ctl *Lists) Close() error {
	return ctl.service.Close()
}
