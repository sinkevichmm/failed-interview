package file

import (
	"failed-interview/01/internal/models"
)

type Repository interface {
	SaveFile(file *models.File) (id string, err error)
	GetFileInfoByID(id string) (meta *models.Meta, err error)
	GetFileIDs() (ids []string)
	GetFreeCapacity() uint
	DeleteFile(id string) (err error)
	GetFileByID(id string) (file *models.File, err error)
}
