package domain

import (
	"context"
	"errors"
)

var (
	_ ListInterface = (*ListService)(nil)
)

var (
	ErrToDoServiceCreateList = errors.Join(errToDoService, errors.New("create list failed"))
	ErrToDoServiceDeleteList = errors.Join(errToDoService, errors.New("delete list failed"))
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

// Close implements ListInterface.
func (s *ListService) Close() error {
	return s.provider.Close()
}

// Create implements ListInterface.
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

// Delete implements ListInterface.
func (s *ListService) Delete(ctx context.Context, userID UserID, listID ListID) error {
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, connection Connection) error {

		// listID ckeck for connection to the user ? s.listRepo.Ckeck ?

		return s.listRepo.Delete(ctx, connection, userID, listID)
	})
	if err != nil {
		return errors.Join(ErrToDoServiceDeleteList, err)
	}

	return nil
}

// ReadAll implements ListInterface.
func (s *ListService) ReadAll(context.Context, UserID) ([]List, error) {
	panic("unimplemented")
}

// Update implements ListInterface.
func (s *ListService) Update(context.Context, List) error {
	panic("unimplemented")
}
