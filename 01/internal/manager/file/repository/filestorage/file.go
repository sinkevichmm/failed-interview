package filestorage

import (
	"failed-interview/01/internal/models"
	f "failed-interview/01/pkg/filestorage"
	"log"
)

type FileStorage struct {
	fs *f.FileStorage
}

func NewFileStorage(metaFile string, limit uint) *FileStorage {
	ff, err := f.OpenFileStorage(metaFile, limit)
	if err != nil {
		log.Fatalln(err)
	}

	return &FileStorage{fs: ff}
}

func (fs *FileStorage) SaveFile(file *models.File) (id string, err error) {
	return fs.fs.SaveFile(file)
}

func (fs *FileStorage) GetFileIDs() (ids []string) {
	return fs.fs.GetFileIDs()
}

func (fs *FileStorage) GetFileInfoByID(id string) (meta *models.Meta, err error) {
	return fs.fs.GetFileInfoByID(id)
}

func (fs *FileStorage) GetFreeCapacity() uint {
	return fs.fs.GetFreeCapacity()
}

func (fs *FileStorage) DeleteFile(id string) (err error) {
	return fs.fs.DeleteFile(id)
}

func (fs *FileStorage) GetFileByID(id string) (file *models.File, err error) {
	return fs.fs.GetFileByID(id)
}
