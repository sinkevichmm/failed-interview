package balance

import (
	"failed-interview/02/entity"
)

// Service balance usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// checkBalance check a balance
func (s *Service) checkBalance(idFrom int, idTo int, value int) error {
	b, err := s.repo.Get(idFrom, idTo)

	if err != nil {
		return err
	}

	if b == nil || len(b) == 0 {
		return entity.ErrNotFound
	}

	if len(b) == 1 {
		if b[0].ID != idFrom {
			return entity.ErrIDFromNotFound
		}

		return entity.ErrIDToNotFound
	}

	bFrom := &entity.Balance{}
	if b[0].ID == idFrom {
		bFrom.Value = b[0].Value
	} else {
		bFrom.Value = b[1].Value
	}

	return bFrom.Validate(value)
}

// ListBalances list balances
func (s *Service) ListBalances() ([]*entity.Balance, error) {
	balances, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	if len(balances) == 0 {
		return nil, entity.ErrNotFound
	}

	return balances, nil
}

// UpdateBalance Update a balance
func (s *Service) UpdateBalance(idFrom int, idTo int, value int) error {
	err := s.checkBalance(idFrom, idTo, value)
	if err != nil {
		return err
	}

	return s.repo.Update(idFrom, idTo, value)
}
