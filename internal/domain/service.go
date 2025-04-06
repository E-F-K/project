package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	_ UserInterface = (*UserService)(nil)
	_ ListInterface = (*ListService)(nil)
	_ TaskInterface = (*TaskService)(nil)
)

var (
	errToDoService                    = errors.New("service error")
	ErrToDoServiceAuthenticate        = errors.Join(errToDoService, errors.New("authenticate user failed"))
	ErrToDoServiceRegisterUser        = errors.Join(errToDoService, errors.New("register user failed"))
	ErrToDoServiceLoginUser           = errors.Join(errToDoService, errors.New("login user failed"))
	ErrToDoServiceInvalidPasswordUser = errors.Join(ErrToDoServiceLoginUser, errors.New("invalid password email"))
	ErrToDoServiceUpdateToken         = errors.Join(errToDoService, errors.New("update token failed"))
	ErrToDoServiceCreateList          = errors.Join(errToDoService, errors.New("create list failed"))
)

type UserService struct {
	provider ConnectionProvider
	userRepo UsersRepository
}

func NewUserService(provider ConnectionProvider, userRepo UsersRepository) *UserService {
	return &UserService{
		provider: provider,
		userRepo: userRepo,
	}
}

func (s *UserService) Authenticate(ctx context.Context, token string) (User, error) {
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

func (s *UserService) RegisterUser(ctx context.Context, name string, email string, passwordHash string, token string) error {
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

func (s *UserService) Login(ctx context.Context, email string, password string) error {
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

func (s *UserService) UpdateToken(ctx context.Context, email string, token string) error {
	err := s.provider.Execute(ctx, func(ctx context.Context, connection Connection) error {
		return s.userRepo.UpdateTokenByEmail(ctx, connection, email, token)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceUpdateToken, err)
	}

	return nil

}

/*
	func (s *ToDoService) CreateList(ctx context.Context, userID UserID, name string) error {
		err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
			list := List{
				ID:     ListID(uuid.New()),
				UserID: userID,
				Name:   name,
			}

			return s.listRepo.Create(ctx, connection, list)
		})
		if err != nil {
			return errors.Join(ErrToDoServiceCreateList, err)
		}

		return nil
	}
*/
func (s *UserService) Close() error {
	return s.provider.Close()
}

type ListService struct {
	provider ConnectionProvider
	listRepo ListsRepository
}

func NewListService(provider ConnectionProvider, listRepo ListsRepository) *ListService {
	return &ListService{
		provider: provider,
		listRepo: listRepo,
	}
}

// Close implements ListInterface.
func (s *ListService) Close() error {
	return s.provider.Close()
}

// Create implements ListInterface.
func (s *ListService) Create(context.Context, List) error {
	panic("unimplemented")
}

// Delete implements ListInterface.
func (s *ListService) Delete(context.Context, UserID, ListID) error {
	panic("unimplemented")
}

// ReadAll implements ListInterface.
func (s *ListService) ReadAll(context.Context, UserID) ([]List, error) {
	panic("unimplemented")
}

// Update implements ListInterface.
func (s *ListService) Update(context.Context, List) error {
	panic("unimplemented")
}

type TaskService struct {
	provider ConnectionProvider
	taskRepo TasksRepository
}

func NewTaskService(provider ConnectionProvider, taskRepo TasksRepository) *TaskService {
	return &TaskService{
		provider: provider,
		taskRepo: taskRepo,
	}
}

// Close implements TaskInterface.
func (s *TaskService) Close() error {
	return s.provider.Close()
}

// Create implements TaskInterface.
func (s *TaskService) Create(context.Context, Task) error {
	panic("unimplemented")
}

// Delete implements TaskInterface.
func (s *TaskService) Delete(context.Context, TaskID) error {
	panic("unimplemented")
}

// Update implements TaskInterface.
func (s *TaskService) Update(context.Context, Task) error {
	panic("unimplemented")
}
