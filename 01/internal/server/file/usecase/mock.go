package usecase

import (
	"failed-interview/01/internal/models"
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) SaveFile(file *models.File) (id string, err error) {
	f := file

	t, err := time.Parse("2006-01-02T15:04:05.000Z", "2021-01-30T05:55:59.888Z")

	if err != nil {
		fmt.Println(err)
	}

	f.DateUpload = t
	args := m.Called(f)

	return args.Get(0).(string), args.Error(1)
}
func (m *Mock) GetFileInfoByID(id string) (meta *models.Meta, err error) {
	args := m.Called(id)
	return args.Get(0).(*models.Meta), args.Error(1)
}

func (m *Mock) GetFileIDs() (ids []string) {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *Mock) DeleteFile(id string) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Mock) GetFileByID(id string) (file *models.File, err error) {
	args := m.Called(id)
	return args.Get(0).(*models.File), args.Error(1)
}
