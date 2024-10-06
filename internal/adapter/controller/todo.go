package controller

import (
	"io"

	"github.com/gin-gonic/gin"
)

var _ io.Closer = (*ToDo)(nil)

type ToDo struct{}

func NewToDo() *ToDo {
	return &ToDo{}
}

func (c *ToDo) GetUserLists(ctx *gin.Context) {

}

func (c *ToDo) Close() error {
	return nil
}
