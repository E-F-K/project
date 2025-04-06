package controller

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"todo_list/internal/adapter/logger"
	"todo_list/internal/domain"

	"github.com/gin-gonic/gin"
)

const ctxAuthUser = "ctx_auth_user"

type AuthMiddleware struct {
	userService domain.UserInterface
}

func NewAuthMiddleware(userService domain.UserInterface) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

func (mw AuthMiddleware) Auth(c *gin.Context) {
	ctx := c.Request.Context()

	var token string
	{
		header := strings.Trim(c.GetHeader("Authorization"), " ")
		if split := strings.Split(header, " "); len(split) == 2 && strings.ToLower(split[0]) == "bearer" {
			token = split[1]
		} else {
			slog.WarnContext(ctx, "Parsing bearer token failed.", logger.ErrAttr(errors.New("invalid bearer token")))

			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}
	}

	curUser, err := mw.userService.Authenticate(ctx, token)
	if err != nil {
		slog.WarnContext(ctx, "Bearer token authentication failed.", logger.ErrAttr(err))

		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}
	c.Set(ctxAuthUser, curUser)

	c.Next()
}

func getCurrentUser(c *gin.Context) domain.User {
	user, _ := c.Get(ctxAuthUser)

	return user.(domain.User)
}
