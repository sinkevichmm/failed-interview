package filestorage

import (
	"encoding/json"
	"failed-interview/01/internal/models"
	"failed-interview/01/pkg/filestorage/idgenerator"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type IDGenerator interface {
	GenerateID() (id string)
}

type FileStorage struct {
	Storage *models.Storage
	IDGenerator
}

func newFileStorage() *FileStorage {
	return &FileStorage{
		Storage: &models.Storage{FilesMap: make(map[string]*models.Meta),
			Mutex: new(sync.Mutex),
		},

		IDGenerator: idgenerator.NewIDGenerator(),
	}
}

func (fs *FileStorage) readMetaFile() (err error) {
	var f []byte

	if _, err := os.Stat(fs.Storage.MetaFile); os.IsNotExist(err) {
		file, err := os.Create(fs.Storage.MetaFile)
		if err != nil {
			return fmt.Errorf("%w %s", ErrMetaRead, fs.Storage.MetaFile)
		}
		defer file.Close()

		fs.Storage.FilesMap = make(map[string]*models.Meta)

		return err
	}

	f, err = ioutil.ReadFile(fs.Storage.MetaFile)
	if err != nil {
		return fmt.Errorf("%w %s", ErrMetaRead, fs.Storage.MetaFile)
	}

	if len(f) == 0 {
		fs.Storage.FilesMap = make(map[string]*models.Meta)

		return err
	}

	m := make([]models.Meta, 0)

	err = json.Unmarshal(f, &m)
	if err != nil {
		return fmt.Errorf("%w %s", ErrMetaInvalid, fs.Storage.MetaFile)
	}

	fs.Storage.Mutex.Lock()
	defer fs.Storage.Mutex.Unlock()

	fs.Storage.FilesMap = make(map[string]*models.Meta)
	// для удобного поиска
	for i := 0; i < len(m); i++ {
		if _, ok := fs.Storage.FilesMap[m[i].ID]; ok {
			return fmt.Errorf("%w %s", ErrMetaDouble, fs.Storage.MetaFile)
		}

		fs.Storage.FilesMap[m[i].ID] = &m[i]
	}

	return err
}

func (fs *FileStorage) writeMetaFile() (err error) {
	m := make([]models.Meta, 0, len(fs.Storage.FilesMap))

	fs.Storage.Mutex.Lock()
	defer fs.Storage.Mutex.Unlock()

	for _, v := range fs.Storage.FilesMap {
		m = append(m, *v)
	}

	file, _ := json.MarshalIndent(m, "", " ")

	err = ioutil.WriteFile(fs.Storage.MetaFile, file, 0600)

	if err != nil {
		return fmt.Errorf("%w %s", ErrMetaWrite, fs.Storage.MetaFile)
	}

	return err
}

func OpenFileStorage(metaFile string, limit uint) (fs *FileStorage, err error) {
	//NOTE: сделать проверку на существование файлов
	fs = newFileStorage()
	fs.Storage.MetaFile = metaFile

	//
	if limit == uint(0) {
		fmt.Println("Warning limit size file storage is 0")
	}

	fs.Storage.SetLimit(limit)

	err = fs.readMetaFile()

	if err != nil {
		return nil, err
	}

	return fs, err
}

func (fs *FileStorage) SaveFile(file *models.File) (id string, err error) {
	id = fs.GenerateID()

	if fs.Storage.GetFreeCapacity() == 0 {
		return "", ErrStorageFull
	}

	p, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", fmt.Errorf("%w %s.%s", ErrFileSave, file.Name, file.Extension)
	}

	err = ioutil.WriteFile(p+string(os.PathSeparator)+id, file.File, 0600)

	if err != nil {
		return "", fmt.Errorf("%w %s.%s", ErrFileSave, file.Name, file.Extension)
	}

	fs.Storage.Mutex.Lock()

	if _, ok := fs.Storage.FilesMap[id]; ok {
		return "", fmt.Errorf(" %s %w", id, ErrIDExist)
	}

	file.ID = id

	fs.Storage.FilesMap[id] = fileToMeta(file)

	fs.Storage.Mutex.Unlock()

	err = fs.writeMetaFile()
	if err != nil {
		return "", err
	}

	err = fs.readMetaFile()
	if err != nil {
		return "", err
	}

	return id, err
}

func (fs *FileStorage) GetFileIDs() (ids []string) {
	ids = make([]string, 0, len(fs.Storage.FilesMap))

	fs.Storage.Mutex.Lock()
	defer fs.Storage.Mutex.Unlock()

	for k := range fs.Storage.FilesMap {
		ids = append(ids, k)
	}

	return ids
}

func (fs *FileStorage) GetFileInfoByID(id string) (meta *models.Meta, err error) {
	fs.Storage.Mutex.Lock()
	defer fs.Storage.Mutex.Unlock()

	meta = fs.Storage.FilesMap[id]

	if meta == nil {
		return nil, fmt.Errorf(" %s %w", id, ErrIDNotFound)
	}

	return meta, err
}

func (fs *FileStorage) GetFreeCapacity() uint {
	return fs.Storage.GetFreeCapacity()
}

func (fs *FileStorage) DeleteFile(id string) (err error) {
	fs.Storage.Mutex.Lock()

	if _, ok := fs.Storage.FilesMap[id]; !ok {
		fs.Storage.Mutex.Unlock()
		return fmt.Errorf(" %s %w", id, ErrIDNotFound)
	}

	fs.Storage.Mutex.Unlock()

	p, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("%w %s", ErrFileDelete, id)
	}

	err = os.Remove(p + string(os.PathSeparator) + id)
	if err != nil {
		return fmt.Errorf("%w %s", ErrFileDelete, id)
	}

	fs.Storage.Mutex.Lock()

	delete(fs.Storage.FilesMap, id)

	fs.Storage.Mutex.Unlock()

	err = fs.writeMetaFile()
	if err != nil {
		return err
	}

	return err
}

func (fs *FileStorage) GetFileByID(id string) (file *models.File, err error) {
	fs.Storage.Mutex.Lock()
	if _, ok := fs.Storage.FilesMap[id]; !ok {
		fs.Storage.Mutex.Unlock()
		return file, fmt.Errorf(" %s %w", id, ErrIDNotFound)
	}

	m := fs.Storage.FilesMap[id]

	fs.Storage.Mutex.Unlock()

	p, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return file, fmt.Errorf("%w %s", ErrFileRead, id)
	}

	f, err := ioutil.ReadFile(p + string(os.PathSeparator) + id)
	if err != nil {
		return file, fmt.Errorf("%w %s", ErrFileRead, id)
	}

	file = metaToFile(m)
	file.File = f

	return file, err
}
