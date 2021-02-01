package http

import (
	"failed-interview/01/internal/models"
	"failed-interview/01/internal/server/file"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase file.UseCase
}

func NewHandler(useCase file.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) FileUpload(c echo.Context) error {
	f, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	src, err := f.Open()
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	p, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	// Destination
	dst, err := os.Create(p + string(os.PathSeparator) + f.Filename)
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	defer os.Remove(p + string(os.PathSeparator) + f.Filename)

	ext := ""
	i := strings.Index(f.Filename, ".")

	if i > 0 {
		ext = f.Filename[i+1:]
	}

	file := &models.File{Name: f.Filename, Extension: ext, DateUpload: time.Now()}

	file.File, err = ioutil.ReadFile(p + string(os.PathSeparator) + file.Name)
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	id, err := h.useCase.SaveFile(file)
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with id: %s.</p>", f.Filename, id))
}

func (h *Handler) FileDownload(c echo.Context) error {
	id := c.Param("id")

	file, err := h.useCase.GetFileByID(id)

	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	p, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	err = ioutil.WriteFile(p+string(os.PathSeparator)+file.Name, file.File, 0600)
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	defer os.Remove(p + string(os.PathSeparator) + file.Name)

	return c.Inline(p+string(os.PathSeparator)+file.Name, file.Name)
}

func (h *Handler) FileIDs(c echo.Context) error {
	ids := h.useCase.GetFileIDs()

	s := ""

	for _, id := range ids {
		s += fmt.Sprintf("<p>id: %s</p>", id)
	}

	return c.HTML(http.StatusOK, s)
}

func (h *Handler) FileInfo(c echo.Context) error {
	id := c.Param("id")

	meta, err := h.useCase.GetFileInfoByID(id)
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	var s string

	s += fmt.Sprintf(`<p>DateUpload: %s</p>`, meta.DateUpload)
	s += fmt.Sprintf(`<p>Name: %s</p>`, meta.Name)
	s += fmt.Sprintf(`<p>Extension: %s</p>`, meta.Extension)

	return c.HTML(http.StatusOK, s)
}

func (h *Handler) FileDelete(c echo.Context) error {
	id := c.FormValue("id")

	err := h.useCase.DeleteFile(id)
	if err != nil {
		log.Println(err)
		return c.HTML(http.StatusInternalServerError, err.Error())
	}

	return c.HTML(http.StatusOK, fmt.Sprintf(`<p>File: %s is deleted</p>`, id))
}
