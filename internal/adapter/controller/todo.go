package controller

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var _ io.Closer = (*ToDo)(nil)

type ToDo struct{}

func NewToDo() *ToDo {
	return &ToDo{}
}

func (ctl *ToDo) GetUserLists(c *gin.Context) {
	curUser := getCurrentUser(c)

	c.JSON(http.StatusOK, curUser)
}

func (ctl *ToDo) Close() error {
	return nil
}
