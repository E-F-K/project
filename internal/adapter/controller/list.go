package controller

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"todo_list/internal/adapter/logger"
	"todo_list/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var _ io.Closer = (*Lists)(nil)

type Lists struct {
	service domain.ListInterface
}

func NewLists(service domain.ListInterface) *Lists {
	return &Lists{service: service}
}

func (ctl *Lists) GetUserListsAndTasks(c *gin.Context) {
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	type messageType struct {
		UserID uuid.UUID
	}

	var message messageType
	if err = json.Unmarshal(body, &message); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	// UUID check?
	_, err = uuid.Parse(message.UserID.String())
	if err != nil {
		slog.ErrorContext(ctx, "wrong UserID, noy UUID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("wrong UserID, noy UUID."))

		return
	}
	if len(message.UserID) <= 0 {
		slog.ErrorContext(ctx, "Empty userID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty userID."))

		return
	}

	err = ctl.service.ReadAll(ctx, message.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "Read all failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read all failed."))

		return
	}

	// return userID?
	c.JSON(http.StatusOK, struct {
		UserID uuid.UUID `userID:"UserID"`
	}{
		UserID: message.UserID,
	})
}

func (ctl *Lists) CreateList(c *gin.Context) {
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	type messageType struct {
		ListID uuid.UUID
		UserID uuid.UUID
		Name   string
	}

	var message messageType
	if err = json.Unmarshal(body, &message); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	// UUID check?
	_, err = uuid.Parse(message.ListID.String())
	if err != nil {
		slog.ErrorContext(ctx, "wrong listID, noy UUID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("wrong listID, noy UUID."))

		return
	}

	if len(message.ListID) <= 0 {
		slog.ErrorContext(ctx, "Empty listID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty listID."))

		return
	}

	// UUID check?
	_, err = uuid.Parse(message.UserID.String())
	if err != nil {
		slog.ErrorContext(ctx, "wrong UserID, noy UUID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("wrong UserID, noy UUID."))

		return
	}
	if len(message.UserID) <= 0 {
		slog.ErrorContext(ctx, "Empty userID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty userID."))

		return
	}

	if len(message.Name) <= 0 {
		slog.ErrorContext(ctx, "Empty name.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty name."))

		return
	}

	err = ctl.service.CreateList(ctx, message.ListID, message.UserID, message.Name)
	if err != nil {
		slog.ErrorContext(ctx, "Create list failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Create list failed."))

		return
	}

	// return listID?
	c.JSON(http.StatusOK, struct {
		ListID uuid.UUID `listID:"ListID"`
	}{
		ListID: message.ListID,
	})
}

func (ctl *Lists) UpdateList(c *gin.Context) {
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	type messageType struct {
		ListID  uuid.UUID
		NewName string
	}

	var message messageType
	if err = json.Unmarshal(body, &message); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	// UUID check?
	_, err = uuid.Parse(message.ListID.String())
	if err != nil {
		slog.ErrorContext(ctx, "wrong listID, noy UUID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("wrong listID, noy UUID."))

		return
	}

	if len(message.ListID) <= 0 {
		slog.ErrorContext(ctx, "Empty listID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty listID."))

		return
	}

	if len(message.NewName) <= 0 {
		slog.ErrorContext(ctx, "Empty Name.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty Name."))

		return
	}

	err = ctl.service.UpdateName(ctx, message.ListID, message.NewName)
	if err != nil {
		slog.ErrorContext(ctx, "Update name failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Update name failed."))

		return
	}

	// return listID?
	c.JSON(http.StatusOK, struct {
		ListID uuid.UUID `listID:"ListID"`
	}{
		ListID: message.ListID,
	})
}

func (ctl *Lists) DeleteList(c *gin.Context) {
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	type messageType struct {
		ListID uuid.UUID
	}

	var message messageType
	if err = json.Unmarshal(body, &message); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	// UUID check?
	_, err = uuid.Parse(message.ListID.String())
	if err != nil {
		slog.ErrorContext(ctx, "wrong listID, noy UUID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("wrong listID, noy UUID."))

		return
	}

	if len(message.ListID) <= 0 {
		slog.ErrorContext(ctx, "Empty listID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty listID."))

		return
	}

	err = ctl.service.DeleteList(ctx, message.ListID)
	if err != nil {
		slog.ErrorContext(ctx, "Delete list failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Delete list failed."))

		return
	}

	// return listID?
	c.JSON(http.StatusOK, struct {
		ListID uuid.UUID `listID:"ListID"`
	}{
		ListID: message.ListID,
	})
}

func (ctl *Lists) Close() error {
	return ctl.service.Close()
}
