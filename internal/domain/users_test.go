package domain_test

import (
	"context"
	"errors"
	"testing"

	"todo_list/internal/domain"
	dbMocks "todo_list/mocks/todo_list/src/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUsersUnit(t *testing.T) {
	validEmail := "some@email.ru"
	validPassword := "correct password"
	invalidEmail := "emailNotExist@mail.ru"
	invalidPassword := "wrong password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(validPassword), bcrypt.DefaultCost)
	validUser := domain.User{
		ID:           domain.UserID(uuid.New()),
		Name:         "name",
		Email:        validEmail,
		PasswordHash: string(hashedPassword),
	}

	tests := []struct {
		name            string
		email, password string
		prepareMocks    func(*dbMocks.MockUsersRepository)
		check           func(*testing.T, error)
	}{
		{
			name:     "Success",
			email:    validEmail,
			password: validPassword,
			prepareMocks: func(repo *dbMocks.MockUsersRepository) {
				repo.EXPECT().ReadByEmail(mock.Anything, mock.Anything, validEmail).
					Return(validUser, nil).
					Once()
			},
			check: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:     "Failed - user not found",
			email:    invalidEmail,
			password: validPassword,
			prepareMocks: func(repo *dbMocks.MockUsersRepository) {
				repo.EXPECT().ReadByEmail(mock.Anything, mock.Anything, invalidEmail).
					Return(domain.User{}, errors.New("some error")).
					Once()
			},
			check: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "some error")
				require.ErrorIs(t, err, domain.ErrToDoServiceLoginUser)
			},
		},
		{
			name:     "Failed - wrong password",
			email:    validEmail,
			password: invalidPassword,
			prepareMocks: func(repo *dbMocks.MockUsersRepository) {
				repo.EXPECT().ReadByEmail(mock.Anything, mock.Anything, validEmail).
					Return(validUser, nil).
					Once()
			},
			check: func(t *testing.T, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, domain.ErrToDoServiceInvalidPasswordUser)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := newFakeProvider(dbMocks.NewMockConnection(t))
			repository := dbMocks.NewMockUsersRepository(t)

			if test.prepareMocks != nil {
				test.prepareMocks(repository)
			}

			err := domain.NewUserService(provider, repository).Login(context.Background(), test.email, test.password)

			test.check(t, err)
		})
	}

}

type fakeProvider struct {
	connection domain.Connection
}

func newFakeProvider(connection domain.Connection) *fakeProvider {
	return &fakeProvider{
		connection: connection,
	}
}

func (f *fakeProvider) Close() error { return nil }

func (f *fakeProvider) Execute(ctx context.Context, receiver func(context.Context, domain.Connection) error) error {
	return receiver(ctx, f.connection)
}

func (f *fakeProvider) ExecuteTx(ctx context.Context, receiver func(context.Context, domain.Connection) error) error {
	return receiver(ctx, f.connection)
}

var _ domain.ConnectionProvider = (*fakeProvider)(nil)
