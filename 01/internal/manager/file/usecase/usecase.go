package usecase

import (
	"failed-interview/01/internal/manager/file"
	"failed-interview/01/internal/models"
)

type FileUseCase struct {
	fileRepo file.Repository
}

func NewFileUseCase(fileRepo file.Repository) *FileUseCase {
	return &FileUseCase{
		fileRepo: fileRepo,
	}
}

func (f FileUseCase) SaveFile(file *models.File) (id string, err error) {
	return f.fileRepo.SaveFile(file)
}

func (f FileUseCase) GetFileInfoByID(id string) (meta *models.Meta, err error) {
	return f.fileRepo.GetFileInfoByID(id)
}

func (f FileUseCase) GetFileIDs() (ids []string) {
	return f.fileRepo.GetFileIDs()
}

func (f FileUseCase) GetFreeCapacity() uint {
	return f.fileRepo.GetFreeCapacity()
}

func (f FileUseCase) DeleteFile(id string) (err error) {
	return f.fileRepo.DeleteFile(id)
}

func (f FileUseCase) GetFileByID(id string) (file *models.File, err error) {
	return f.fileRepo.GetFileByID(id)
}
