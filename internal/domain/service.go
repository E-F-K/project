package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var _ UserService = (*ToDoService)(nil)

var (
	errToDoService             = errors.New("service error")
	ErrToDoServiceRegisterUser = errors.Join(errToDoService, errors.New("register user failed"))
)

type ToDoService struct {
	provider ConnectionProvider
	userRepo UsersRepository
}

func NewToDoService(provider ConnectionProvider, userRepo UsersRepository) *ToDoService {
	return &ToDoService{
		provider: provider,
		userRepo: userRepo,
	}
}

func (s *ToDoService) Authenticate(ctx context.Context, token string) (User, error) {
	panic("unimplemented")
}

func (s *ToDoService) RegisterUser(ctx context.Context, name string, email string, passwordHash string, token string) error {
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		user := User{
			ID:           UserID(uuid.New()),
			Name:         name,
			Email:        email,
			PasswordHash: passwordHash,
			Token:        token,
		}

		return s.userRepo.Create(ctx, connection, user)
	})
	if err !=nil {
		return errors.Join(ErrToDoServiceRegisterUser, err)
	}

	return nil
}

func (s *ToDoService) Close() error {
	return s.provider.Close()
}
