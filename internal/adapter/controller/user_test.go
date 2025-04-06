package controller_test

import (
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
	httpCall := func(serviceMock *mocks.MockUserInterface, request *http.Request) *httptest.ResponseRecorder {
		router, ctl, response := gin.New(), controller.NewUsers(serviceMock), httptest.NewRecorder()
		defer func() { _ = ctl.Close() }()
		router.POST("/", ctl.Register)
		router.ServeHTTP(response, request)

		return response
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
