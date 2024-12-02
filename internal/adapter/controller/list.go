package controller

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"todo_list/internal/adapter/logger"
	"todo_list/internal/domain"

	"github.com/gin-gonic/gin"
)

var _ io.Closer = (*Lists)(nil)

type Lists struct {
	service domain.ListService
}

// Close implements io.Closer.
func (ctl *Lists) Close() error {
	panic("unimplemented")
}

func NewLists(service domain.ListService) *Lists {
	return &Lists{service: service}
}

func (ctl *Lists) Create(c *gin.Context) {
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	type messageType struct {
		UserID string
		Name   string
		Email  string
	}
	var message messageType
	if err = json.Unmarshal(body, &message); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	//UserID check??
	if len(message.UserID) <= 0 {
		slog.ErrorContext(ctx, "Empty userID.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty userID."))

		return
	}

	parsedEmail, err := mail.ParseAddress(message.Email)
	if err != nil {
		slog.ErrorContext(ctx, "Parse email failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse email failed."))

		return
	}
	if len(message.Name) <= 0 {
		slog.ErrorContext(ctx, "Empty name.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty name."))

		return
	}

	err = ctl.service.CreateList(ctx, message.UserID, message.Name, parsedEmail.Address)
	if err != nil {
		slog.ErrorContext(ctx, "Create list failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Create list failed."))

		return
	}

	c.JSON(http.StatusOK, struct {
		ListID string `listID:"ListID"`
	}{
		ListID: "aokay",
	})
}
