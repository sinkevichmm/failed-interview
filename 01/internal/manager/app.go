package manager

import (
	g "failed-interview/01/internal/manager/file/delivery/grpc"
	"failed-interview/01/internal/manager/file/repository/filestorage"
	uc "failed-interview/01/internal/manager/file/usecase"
)

type App struct {
	GRPCFS *g.FileServer
}

func NewApp(metaFile *string, limit *uint, port *string, auth *string) *App {
	fs := filestorage.NewFileStorage(*metaFile, *limit)
	fsRepo := uc.NewFileUseCase(fs)
	gfs := g.NewFileServer(fsRepo, *port, *auth)

	return &App{
		GRPCFS: gfs,
	}
}
