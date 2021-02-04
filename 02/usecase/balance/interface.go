package balance

import (
	"failed-interview/02/entity"
)

// Reader interface
type Reader interface {
	Get(idFrom int, idTo int) ([]*entity.Balance, error)
	List() ([]*entity.Balance, error)
}

// Writer balance writer
type Writer interface {
	Update(idFrom int, idTo int, value int) error
}

// Repository interface
type Repository interface {
	Reader
	Writer
}

// UseCase interface
type UseCase interface {
	ListBalances() ([]*entity.Balance, error)
	UpdateBalance(idFrom int, idTo int, value int) error
}
