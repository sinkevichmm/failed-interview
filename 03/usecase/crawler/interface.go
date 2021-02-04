package crawler

import (
	"failed-interview/03/entity"
)

// Reader interface
type Reader interface {
	List(links []string, timeout uint) []*entity.Links
}

// Repository interface
type Repository interface {
	Reader
}

// UseCase interface
type UseCase interface {
	GetList(links []string, timeout uint) []*entity.Links
}
