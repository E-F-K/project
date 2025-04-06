package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/mail"

	"todo_list/internal/adapter/logger"
	"todo_list/internal/domain"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordBCryptoCost  = 12
	passwordHashBytesLen = 16
)

var _ io.Closer = (*Users)(nil)

type Users struct {
	service domain.UserInterface
}

func NewUsers(service domain.UserInterface) *Users {
	return &Users{service: service}
}

func (ctl *Users) Register(c *gin.Context) {
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	type messageType struct {
		Name     string
		Email    string
		Password string
	}
	var message messageType
	if err = json.Unmarshal(body, &message); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

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
	if len(message.Password) <= 0 {
		slog.ErrorContext(ctx, "Empty password.")
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Empty password."))

		return
	}

	var passwordHash string
	{
		passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(message.Password), passwordBCryptoCost)
		if err != nil {
			slog.ErrorContext(ctx, "Hashing password failed.", logger.ErrAttr(err))
			c.JSON(http.StatusUnprocessableEntity, errorResponse("Hashing password failed."))

			return
		}
		passwordHash = string(passwordHashBytes)
	}

	var token string
	{
		token, err = ctl.generateToken()
		if err != nil {
			slog.ErrorContext(ctx, "Create token failed.", logger.ErrAttr(err))
			c.JSON(http.StatusUnprocessableEntity, errorResponse("Create token failed."))

			return
		}
	}

	err = ctl.service.RegisterUser(ctx, message.Name, parsedEmail.Address, passwordHash, token)
	if err != nil {
		slog.ErrorContext(ctx, "Register user failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Register user failed."))

		return
	}

	c.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (ctl *Users) Login(c *gin.Context) {
	ctx := c.Request.Context()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Read request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Read body failed."))

		return
	}

	type messageType struct {
		Email    string
		Password string
	}
	var message messageType
	if err = json.Unmarshal(body, &message); err != nil {
		slog.ErrorContext(ctx, "Parse request body failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Parse body failed."))

		return
	}

	err = ctl.service.Login(ctx, message.Email, message.Password)
	if err != nil {
		slog.ErrorContext(ctx, "Login failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Login failed."))

		return
	}

	var token string
	{
		token, err = ctl.generateToken()
		if err != nil {
			slog.ErrorContext(ctx, "Create token failed.", logger.ErrAttr(err))
			c.JSON(http.StatusUnprocessableEntity, errorResponse("Create token failed."))

			return
		}
	}

	err = ctl.service.UpdateToken(ctx, message.Email, token)
	if err != nil {
		slog.ErrorContext(ctx, "Update user token failed.", logger.ErrAttr(err))
		c.JSON(http.StatusUnprocessableEntity, errorResponse("Update user token failed."))

		return
	}

	c.JSON(http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (ctl *Users) Close() error {
	return ctl.service.Close()
}

func (ctl *Users) generateToken() (string, error) {
	tokenBytes := make([]byte, passwordHashBytesLen)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(tokenBytes), nil
}
