package http

import (
	"failed-interview/01/internal/server/file"

	"github.com/labstack/echo/v4"
)

func RegisterHTTPEndpoints(r *echo.Echo, uc file.UseCase) {
	h := NewHandler(uc)
	r.POST("/file/upload", h.FileUpload)
	r.GET("/file/download/:id", h.FileDownload)
	r.GET("/file/ids", h.FileIDs)
	r.GET("/file/info/:id", h.FileInfo)
	r.POST("/file/delete", h.FileDelete)
}
