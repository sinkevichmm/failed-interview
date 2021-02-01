package filestorage

import (
	"errors"
	"failed-interview/01/internal/models"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadMetaFile(t *testing.T) {
	fs := newFileStorage()
	fs.Storage.MetaFile = "meta_test_ok.json"
	err := fs.readMetaFile()
	require.NoError(t, err)

	fs = newFileStorage()
	fs.Storage.MetaFile = "_meta_test.json"
	err = fs.readMetaFile()
	require.NoError(t, err)

	err = os.Remove(fs.Storage.MetaFile)
	require.NoError(t, err)

	fs = newFileStorage()
	fs.Storage.MetaFile = "meta_test_double.json"
	err = fs.readMetaFile()
	require.ErrorIs(t, ErrMetaDouble, errors.Unwrap(err))

	fs = newFileStorage()
	fs.Storage.MetaFile = "meta_test_corrupt.json"
	err = fs.readMetaFile()
	require.ErrorIs(t, ErrMetaInvalid, errors.Unwrap(err))
}

func TestOpenFileStorage(t *testing.T) {
	_, err := OpenFileStorage("meta_test_ok.json", 5)
	require.NoError(t, err)

	_, err = OpenFileStorage("_meta_test.json", 5)
	require.NoError(t, err)

	err = os.Remove("_meta_test.json")
	require.NoError(t, err)

	_, err = OpenFileStorage("meta_test_double.json", 5)
	require.ErrorIs(t, ErrMetaDouble, errors.Unwrap(err))

	_, err = OpenFileStorage("meta_test_corrupt.json", 5)
	require.ErrorIs(t, ErrMetaInvalid, errors.Unwrap(err))
}

func TestSaveMetaFile(t *testing.T) {
	fs := newFileStorage()
	fs.Storage.MetaFile = "meta_test_ok.json"
	err := fs.readMetaFile()
	require.NoError(t, err)

	fs.Storage.MetaFile = "meta_test_write.json"
	err = fs.writeMetaFile()
	require.NoError(t, err)

	fs2 := newFileStorage()
	fs2.Storage.MetaFile = "meta_test_write.json"
	err = fs2.readMetaFile()
	require.NoError(t, err)

	for k := range fs.Storage.FilesMap {
		if _, ok := fs2.Storage.FilesMap[k]; !ok {
			t.Errorf("not equal")
		}
	}
}

func TestSaveFile(t *testing.T) {
	var err error

	fs := newFileStorage()
	fs.Storage.MetaFile = "meta_test_save.json"

	file := &models.File{}

	fs.Storage.SetLimit(1)

	file.Name = "file for save.pdf"

	file.File, err = ioutil.ReadFile(file.Name)
	require.NoError(t, err)

	_, err = fs.SaveFile(file)
	require.NoError(t, err)

	fs.Storage.SetLimit(0)

	_, err = fs.SaveFile(file)
	require.Error(t, ErrStorageFull, err)
}

func TestGetFileIDs(t *testing.T) {
	fs := newFileStorage()
	fs.Storage.MetaFile = "meta_test_ok.json"
	err := fs.readMetaFile()
	require.NoError(t, err)

	ids := fs.GetFileIDs()

	for _, v := range ids {
		if _, ok := fs.Storage.FilesMap[v]; !ok {
			t.Errorf("not equal")
		}
	}
}

func TestGetFileInfoByID(t *testing.T) {
	fs := newFileStorage()
	fs.Storage.MetaFile = "meta_test_ok.json"
	err := fs.readMetaFile()
	require.NoError(t, err)

	id := "3"
	_, err = fs.GetFileInfoByID(id)

	require.NoError(t, err)

	id = "-3"
	_, err = fs.GetFileInfoByID(id)

	require.ErrorIs(t, ErrIDNotFound, errors.Unwrap(err))
}

func TestGetFreeCapacity(t *testing.T) {
	fs := newFileStorage()
	fs.Storage.SetLimit(6)
	fs.Storage.MetaFile = "meta_test_ok.json"
	err := fs.readMetaFile()
	require.NoError(t, err)

	c := fs.GetFreeCapacity()

	require.Equal(t, uint(3), c)
}

func TestDeleteFile(t *testing.T) {
	var err error

	fs := newFileStorage()
	fs.Storage.SetLimit(3)
	fs.Storage.MetaFile = "meta_test_ok.json"
	err = fs.readMetaFile()
	require.NoError(t, err)

	fs.Storage.MetaFile = "meta_test_tmp.json"

	id := "7"

	err = fs.DeleteFile(id)
	require.ErrorIs(t, ErrIDNotFound, errors.Unwrap(err))

	id = "3"

	os.Args[0] = ""

	err = fs.DeleteFile(id)
	require.NoError(t, err)

	input, err := ioutil.ReadFile("2")
	require.NoError(t, err)

	err = ioutil.WriteFile("3", input, 0600)
	require.NoError(t, err)
}

func TestGetFileByID(t *testing.T) {
	var err error

	tm, err := time.Parse(time.RFC3339, "2020-09-01T21:46:43Z")
	require.NoError(t, err)

	ff := &models.File{ID: "3", Name: "test3", Extension: "gif", DateUpload: tm}

	fff, err := ioutil.ReadFile("3")
	if err != nil {
		t.Errorf(ErrFileRead.Error())
	}

	ff.File = fff

	fs := newFileStorage()
	fs.Storage.SetLimit(3)
	fs.Storage.MetaFile = "meta_test_ok.json"
	err = fs.readMetaFile()
	require.NoError(t, err)

	fs.Storage.MetaFile = "meta_test_tmp.json"

	id := "7"

	_, err = fs.GetFileByID(id)
	require.ErrorIs(t, ErrIDNotFound, errors.Unwrap(err))

	id = "3"

	f, err := fs.GetFileByID(id)
	require.NoError(t, err)

	if !reflect.DeepEqual(f, ff) {
		t.Errorf("not equal")
	}
}
