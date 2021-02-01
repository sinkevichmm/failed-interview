package http

import (
	"bytes"
	"failed-interview/01/internal/models"
	"failed-interview/01/internal/server/file/usecase"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileUpload(t *testing.T) {
	path := "testfile.txt"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", path)
	require.NoError(t, err)
	sample, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, sample)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	e := echo.New()
	uc := new(usecase.Mock)

	RegisterHTTPEndpoints(e, uc)

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewHandler(uc)

	tt, err := time.Parse("2006-01-02T15:04:05.000Z", "2021-01-30T05:55:59.888Z")
	require.NoError(t, err)

	file := &models.File{Name: "testfile.txt", Extension: "txt", DateUpload: tt}

	file.File, err = ioutil.ReadFile(path)
	require.NoError(t, err)

	uc.On("SaveFile", file).Return("uuid", nil)

	_ = h.FileUpload(c)
}

func TestFileDownload(t *testing.T) {
	e := echo.New()
	uc := new(usecase.Mock)

	RegisterHTTPEndpoints(e, uc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := NewHandler(uc)
	c := e.NewContext(req, w)
	path := "testfile.txt"

	uc.On("GetFileByID", "").Return(&models.File{Name: path}, nil)

	_ = h.FileDownload(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFileIDs(t *testing.T) {
	e := echo.New()
	uc := new(usecase.Mock)

	RegisterHTTPEndpoints(e, uc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := NewHandler(uc)
	c := e.NewContext(req, w)

	uc.On("GetFileIDs").Return([]string(nil))

	_ = h.FileIDs(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFileInfo(t *testing.T) {
	e := echo.New()
	uc := new(usecase.Mock)

	RegisterHTTPEndpoints(e, uc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := NewHandler(uc)
	c := e.NewContext(req, w)

	uc.On("GetFileInfoByID", "").Return(&models.Meta{}, nil)

	_ = h.FileInfo(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFileDelete(t *testing.T) {
	e := echo.New()
	uc := new(usecase.Mock)

	RegisterHTTPEndpoints(e, uc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := NewHandler(uc)
	c := e.NewContext(req, w)

	uc.On("DeleteFile", "").Return(nil)

	_ = h.FileDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
