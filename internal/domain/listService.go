package domain

import (
	"context"
	"errors"
)

var (
	_ ListInterface = (*ListService)(nil)
)

var (
	ErrToDoServiceCreateList       = errors.Join(errToDoService, errors.New("create list failed"))
	ErrToDoServiceDeleteList       = errors.Join(errToDoService, errors.New("delete list failed"))
	ErrToDoServiceListAccessDenied = errors.Join(errToDoService, errors.New("access to list failed"))
	ErrToDoServiceReadAllLits      = errors.Join(errToDoService, errors.New("read all lists failed"))
	ErrToDoServiceUpdateList       = errors.Join(errToDoService, errors.New("update lists failed"))
)

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

func (s *ListService) Close() error {
	return s.provider.Close()
}

func (s *ListService) Create(ctx context.Context, list List) error {
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		list := List{
			ID:        list.ID,
			UserID:    list.UserID,
			Name:      list.Name,
			UpdatedAt: list.UpdatedAt,
			Tasks:     list.Tasks,
		}

		return s.listRepo.Create(ctx, connection, list)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceCreateList, err)
	}

	return nil
}

func (s *ListService) Delete(ctx context.Context, userID UserID, listID ListID) error {
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		return s.listRepo.Delete(ctx, connection, userID, listID)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceDeleteList, err)
	}

	return nil
}

func (s *ListService) GetAll(ctx context.Context, userID UserID) ([]List, error) {
	var lists []List
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		var getError error
		lists, getError = s.listRepo.GetAllLists(ctx, connection, userID)
		return getError
	})
	if err != nil {
		// first argument?
		return lists, errors.Join(ErrToDoServiceReadAllLits, err)
	}

	return lists, nil
}

func (s *ListService) Update(ctx context.Context, list List) error {
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		// Мы присваиваем юзер ID текущему юзеру в listController
		return s.listRepo.Update(ctx, connection, list)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceUpdateList, err)
	}

	return nil
}
