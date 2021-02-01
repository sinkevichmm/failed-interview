package server

import (
	"failed-interview/01/internal/server/file"
	h "failed-interview/01/internal/server/file/delivery/http"
	g "failed-interview/01/internal/server/file/repository/grpcfile"
	"failed-interview/01/internal/server/file/usecase"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

type App struct {
	httpServer *http.Server
	fileUC     file.UseCase
	httpport   string
}

func NewApp(httpport string, grpcaddress string, auth string) *App {
	fileRepo := g.NewGRPCService(grpcaddress, auth)

	return &App{
		fileUC:   usecase.NewFileUseCase(fileRepo),
		httpport: httpport,
	}
}

func (a *App) Start() error {
	router := echo.New()

	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	router.Static("/", p+"/public")

	h.RegisterHTTPEndpoints(router, a.fileUC)

	a.httpServer = &http.Server{
		Addr:           ":" + a.httpport,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return a.httpServer.ListenAndServe()
}
