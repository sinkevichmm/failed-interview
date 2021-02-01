package models

import (
	"sync"
	"time"
)

//NOTE: перенести сюда read/write meta
type Storage struct {
	limit    uint
	MetaFile string
	Mutex    *sync.Mutex
	FilesMap map[string]*Meta
}

func (s *Storage) SetLimit(limit uint) {
	s.limit = limit
}

func (s *Storage) GetFreeCapacity() (free uint) {
	return s.limit - uint(len(s.FilesMap))
}

type Meta struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Extension  string    `json:"extension"`
	DateUpload time.Time `json:"dateUpload"`
}
