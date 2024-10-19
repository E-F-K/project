package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var _ UserService = (*ToDoService)(nil)

var (
	errToDoService                    = errors.New("service error")
	ErrToDoServiceAuthenticate        = errors.Join(errToDoService, errors.New("authenticate user failed"))
	ErrToDoServiceRegisterUser        = errors.Join(errToDoService, errors.New("register user failed"))
	ErrToDoServiceLoginUser           = errors.Join(errToDoService, errors.New("login user failed"))
	ErrToDoServiceInvalidPasswordUser = errors.Join(ErrToDoServiceLoginUser, errors.New("invalid password email"))
	ErrToDoServiceUpdateToken         = errors.Join(errToDoService, errors.New("update token failed"))
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
	var user User
	var err error
	err = s.provider.Execute(ctx, func(ctx context.Context, connection Connection) error {
		user, err = s.userRepo.ReadByToken(ctx, connection, token)

		return err
	})
	if err != nil {
		err = errors.Join(ErrToDoServiceAuthenticate, err)
	}

	return user, err
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
	if err != nil {
		return errors.Join(ErrToDoServiceRegisterUser, err)
	}

	return nil
}

func (s *ToDoService) Login(ctx context.Context, email string, password string) error {
	var user User
	err := s.provider.Execute(ctx, func(ctx context.Context, connection Connection) error {
		var err error
		user, err = s.userRepo.ReadByEmail(ctx, connection, email)

		return err
	})
	if err != nil {
		return errors.Join(ErrToDoServiceLoginUser, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return errors.Join(ErrToDoServiceInvalidPasswordUser, err)
	}

	return nil
}

func (s *ToDoService) UpdateToken(ctx context.Context, email string, token string) error {
	err := s.provider.Execute(ctx, func(ctx context.Context, connection Connection) error {
		return s.userRepo.UpdateTokenByEmail(ctx, connection, email, token)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceUpdateToken, err)
	}

	return nil

}

func (s *ToDoService) Close() error {
	return s.provider.Close()
}
