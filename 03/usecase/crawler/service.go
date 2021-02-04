package crawler

import (
	"failed-interview/03/entity"
)

// Service crawler usecase
type Service struct {
	repo Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// ListBalances list crawlers
func (s *Service) GetList(links []string, timeout uint) []*entity.Links {
	return s.repo.List(links, timeout)
}
