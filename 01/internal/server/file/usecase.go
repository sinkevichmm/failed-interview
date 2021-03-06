package file

import (
	"failed-interview/01/internal/models"
)

type UseCase interface {
	SaveFile(file *models.File) (id string, err error)
	GetFileInfoByID(id string) (meta *models.Meta, err error)
	GetFileIDs() (ids []string)
	DeleteFile(id string) (err error)
	GetFileByID(id string) (file *models.File, err error)
}
