package controller_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todo_list/internal/adapter/controller"
	mocks "todo_list/mocks/todo_list/src/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsersRegister(t *testing.T) {
	httpCall := func(usersMock *mocks.MockUserInterface, request *http.Request) *httptest.ResponseRecorder {
		ctl := controller.NewUsers(usersMock)
		defer func() { _ = ctl.Close() }()

		return httpPost(request, ctl.Register)
	}

	tests := []struct {
		name         string
		request      *http.Request
		prepareMocks func(*mocks.MockUserInterface)
		validation   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			prepareMocks: func(serviceMock *mocks.MockUserInterface) {
				serviceMock.EXPECT().RegisterUser(mock.Anything, "John Doe", "johh@doe.foo", mock.Anything, mock.Anything).
					Return(nil).Once()
				serviceMock.EXPECT().Close().Return(nil).Once()
			},
			request: httptest.NewRequest("POST", "/", strings.NewReader(`{
				"name": "John Doe",
				"email": "johh@doe.foo",
				"password": "secret"
			}`)),
			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, response.Code)
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Contains(t, body, `"token"`)
			},
		},
		{
			name: "Register user failed",
			prepareMocks: func(serviceMock *mocks.MockUserInterface) {
				serviceMock.EXPECT().RegisterUser(mock.Anything, "John Doe", "johh@doe.foo", mock.Anything, mock.Anything).
					Return(errors.New("some error")).Once()
				serviceMock.EXPECT().Close().Return(nil).Once()
			},
			request: httptest.NewRequest("POST", "/", strings.NewReader(`{
				"name": "John Doe",
				"email": "johh@doe.foo",
				"password": "secret"
			}`)),
			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Contains(t, body, "Register user failed")
			},
		},
		{
			name: "Empty name",
			prepareMocks: func(serviceMock *mocks.MockUserInterface) {
				serviceMock.EXPECT().Close().Return(nil).Once()
			},
			request: httptest.NewRequest("POST", "/", strings.NewReader(`{
				"name": "",
				"email": "johh@doe.foo",
				"password": "secret"
			}`)),
			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Contains(t, body, "Empty name")
			},
		},
		{
			name: "Empty password",
			prepareMocks: func(serviceMock *mocks.MockUserInterface) {
				serviceMock.EXPECT().Close().Return(nil).Once()
			},
			request: httptest.NewRequest("POST", "/", strings.NewReader(`{
				"name": "John Doe",
				"email": "johh@doe.foo",
				"password": ""
			}`)),
			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Contains(t, body, "Empty password")
			},
		},
		{
			name:    "Parse body failed",
			request: httptest.NewRequest("POST", "/", strings.NewReader(`invalid parse body`)),
			prepareMocks: func(serviceMock *mocks.MockUserInterface) {
				serviceMock.EXPECT().Close().Return(nil).Once()
			},
			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
				require.Contains(t, body, "Parse body failed")
			},
		},
		{
			name: "Parse email failed",
			request: httptest.NewRequest("POST", "/", strings.NewReader(`{
				"name": "John Doe",
				"email": "invalid_email",
				"password": "secret"
			}`)),
			prepareMocks: func(serviceMock *mocks.MockUserInterface) {
				serviceMock.EXPECT().Close().Return(nil).Once()
			},
			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Contains(t, body, "Parse email failed")
			},
		},
	}
	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceMock := mocks.NewMockUserInterface(t)
			if test.prepareMocks != nil {
				test.prepareMocks(serviceMock)
			}
			test.validation(t, httpCall(serviceMock, test.request))
		})
	}
}

func TestUsersLogin(t *testing.T) {
	httpCall := func(usersMock *mocks.MockUserInterface, request *http.Request) *httptest.ResponseRecorder {
		ctl := controller.NewUsers(usersMock)
		defer func() { _ = ctl.Close() }()

		return httpPost(request, ctl.Login)
	}

	tests := []struct {
		name         string
		request      *http.Request
		prepareMocks func(*mocks.MockUserInterface)
		validation   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			request: httptest.NewRequest("POST", "/", strings.NewReader(`{
				"email": "johh@doe.foo",
				"password": "secret"
			}`)),
			prepareMocks: func(mockService *mocks.MockUserInterface) {
				mockService.EXPECT().Login(mock.Anything, "johh@doe.foo", "secret").
					Return(nil).Once()

				mockService.EXPECT().UpdateToken(mock.Anything, "johh@doe.foo", mock.Anything).
					Return(nil).Once()

				mockService.EXPECT().Close().Return(nil).Once()
			},
			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, response.Code)
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Contains(t, body, `"token"`)
			},
		},
		{
			name:    "Parse body failed",
			request: httptest.NewRequest("POST", "/", strings.NewReader(`invalid parse body`)),
			prepareMocks: func(mockService *mocks.MockUserInterface) {
				mockService.EXPECT().Close().Return(nil).Once()
			},
			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
				require.Contains(t, body, "Parse body failed")
			},
		},
		{
			name: "Login failed",
			request: httptest.NewRequest("POST", "/", strings.NewReader(`{
				"email": "johh@doe.foo",
				"password": "wrong"
			}`)),

			prepareMocks: func(mockService *mocks.MockUserInterface) {
				mockService.EXPECT().Login(mock.Anything, "johh@doe.foo", "wrong").
					Return(errors.New("some error")).Once()
				mockService.EXPECT().Close().Return(nil).Once()
			},

			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
				require.Contains(t, body, "Login failed")
			},
		},
		{
			name: "Update token failed",
			prepareMocks: func(mockService *mocks.MockUserInterface) {
				mockService.EXPECT().Login(mock.Anything, "johh@doe.foo", "secret").
					Return(nil).Once()
				mockService.EXPECT().UpdateToken(mock.Anything, "johh@doe.foo", mock.Anything).
					Return(errors.New("some error")).Once()
				mockService.EXPECT().Close().Return(nil).Once()
			},
			request: httptest.NewRequest("POST", "/", strings.NewReader(`{
				"email": "johh@doe.foo",
				"password": "secret"
			}`)),
			validation: func(t *testing.T, response *httptest.ResponseRecorder) {
				body, err := response.Body.ReadString('\n')
				require.ErrorIs(t, err, io.EOF)
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
				require.Contains(t, body, "Update user token failed")
			},
		},
	}
	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceMock := mocks.NewMockUserInterface(t)
			if test.prepareMocks != nil {
				test.prepareMocks(serviceMock)
			}
			test.validation(t, httpCall(serviceMock, test.request))
		})
	}
}

func httpPost(request *http.Request, handler func(*gin.Context)) *httptest.ResponseRecorder {
	router, response := gin.New(), httptest.NewRecorder()
	router.POST("/", handler)
	router.ServeHTTP(response, request)

	return response
}
