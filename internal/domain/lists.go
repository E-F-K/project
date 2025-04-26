package domain

import (
	"context"
	"errors"
)

var (
	_ ListInterface = (*ListService)(nil)
)

var (
	errListService                 = errors.New("list service error")
	ErrListServiceCreate           = errors.Join(errListService, errors.New("create failed"))
	ErrToDoServiceDeleteList       = errors.Join(errListService, errors.New("delete list failed"))
	ErrToDoServiceListAccessDenied = errors.Join(errListService, errors.New("access to list failed"))
	ErrToDoServiceReadAllLits      = errors.Join(errListService, errors.New("read all lists failed"))
	ErrToDoServiceUpdateList       = errors.Join(errListService, errors.New("update lists failed"))
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
		return errors.Join(ErrListServiceCreate, err)
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
		lists, getError = s.listRepo.ReadAll(ctx, connection, userID)
		return getError
	})
	if err != nil {
		return nil, errors.Join(ErrToDoServiceReadAllLits, err)
	}

	return lists, nil
}

func (s *ListService) Update(ctx context.Context, list List) error {
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {
		return s.listRepo.Update(ctx, connection, list)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceUpdateList, err)
	}

	return nil
}
