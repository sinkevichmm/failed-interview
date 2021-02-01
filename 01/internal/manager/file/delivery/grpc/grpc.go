package grpc

import (
	"bufio"
	"bytes"
	"context"
	"failed-interview/01/internal/manager/file"
	"failed-interview/01/internal/models"
	pb "failed-interview/01/pkg/proto"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	//"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type FileServer struct {
	file.UseCase
	pb.UnimplementedFileServiceServer
	auth string
	port string
}

// NOTE: замокировать протестировать
func (g *FileServer) authInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	pref := "/fs.FileService/"
	ctx := stream.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	// пардия на проверку доступа
	switch info.FullMethod {
	case pref + "SaveFile", pref + "GetFileByID", pref + "DeleteFile":
		if md.Get("auth")[0] != g.auth {
			return status.Error(codes.PermissionDenied, "bad access-token")
		}
	}

	return handler(srv, stream)
}

func NewFileServer(uc file.UseCase, port string, auth string) *FileServer {
	return &FileServer{UseCase: uc, port: port, auth: auth}
}

func (g *FileServer) Start() error {
	lis, err := net.Listen("tcp", "0.0.0.0:"+g.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(g.authInterceptor),
	)
	pb.RegisterFileServiceServer(s, g)
	//reflection.Register(s)

	return s.Serve(lis)
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		log.Println(context.Canceled.Error())
		return status.Error(codes.Canceled, context.Canceled.Error())
	case context.DeadlineExceeded:
		log.Println(context.DeadlineExceeded.Error())
		return status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error())
	default:
		return nil
	}
}

func (g *FileServer) SaveFile(s pb.FileService_SaveFileServer) error {
	req, err := s.Recv()

	if err != nil {
		log.Println(err)
		return status.Error(codes.Unknown, err.Error())
	}

	f := &models.File{Name: req.GetFileInfo().Name, Extension: req.GetFileInfo().Extension, DateUpload: req.GetFileInfo().DateUpload.AsTime()}

	fileData := bytes.Buffer{}

	for {
		err := contextError(s.Context())
		if err != nil {
			return err
		}

		req, err := s.Recv()
		if err == io.EOF {
			// приняли файл
			break
		}

		if err != nil {
			log.Println(err)
			return status.Errorf(codes.Unknown, err.Error())
		}

		chunk := req.GetChunkData()

		_, err = fileData.Write(chunk)
		if err != nil {
			log.Println(err)
			return status.Errorf(codes.Internal, err.Error())
		}
	}

	f.File = fileData.Bytes()

	id, err := g.UseCase.SaveFile(f)
	if err != nil {
		log.Println(err)
		return status.Errorf(codes.Internal, err.Error())
	}

	res := &pb.SaveFileResponse{Id: id}

	err = s.SendAndClose(res)
	if err != nil {
		log.Println(err)
		return status.Errorf(codes.Unknown, err.Error())
	}

	return nil
}

func (g *FileServer) GetFileInfoByID(ctx context.Context, p *pb.GetFileInfoByIDRequest) (*pb.GetFileInfoByIDResponse, error) {
	meta, err := g.UseCase.GetFileInfoByID(p.Id)

	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	pbtime, err := ptypes.TimestampProto(meta.DateUpload)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	pb := &pb.GetFileInfoByIDResponse{FileInfo: &pb.FileInfo{}}

	pb.FileInfo.DateUpload = pbtime
	pb.FileInfo.Extension = meta.Extension
	pb.FileInfo.Name = meta.Name

	return pb, nil
}

func (g *FileServer) GetFileIDs(ctx context.Context, p *pb.GetFileIDsRequest) (*pb.GetFileIDsResponse, error) {
	ids := g.UseCase.GetFileIDs()
	pb := &pb.GetFileIDsResponse{Ids: ids}

	return pb, nil
}

func (g *FileServer) GetFreeCapacity(ctx context.Context, p *pb.GetFreeCapacityRequest) (*pb.GetFreeCapacityResponse, error) {
	free := g.UseCase.GetFreeCapacity()
	pb := &pb.GetFreeCapacityResponse{FreeCapacity: uint32(free)}

	return pb, nil
}

func (g *FileServer) DeleteFile(ctx context.Context, p *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	err := g.UseCase.DeleteFile(p.Id)

	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.NotFound, err.Error())
	}

	pb := &pb.DeleteFileResponse{}

	return pb, nil
}

func (g *FileServer) GetFileByID(p *pb.GetFileByIDRequest, pp pb.FileService_GetFileByIDServer) error {
	meta, err := g.UseCase.GetFileInfoByID(p.Id)
	if err != nil {
		log.Println(err)
		return status.Error(codes.NotFound, err.Error())
	}

	file, err := g.UseCase.GetFileByID(p.Id)
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, err.Error())
	}

	pbtime, err := ptypes.TimestampProto(file.DateUpload)
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, err.Error())
	}

	fi := &pb.FileInfo{Name: meta.Name, Extension: meta.Extension, DateUpload: pbtime}

	err = pp.Send(&pb.GetFileByIDResponse{Data: &pb.GetFileByIDResponse_FileInfo{FileInfo: fi}})
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, err.Error())
	}

	r := bytes.NewReader(file.File)

	reader := bufio.NewReader(r)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			return status.Error(codes.Internal, err.Error())
		}

		req := &pb.GetFileByIDResponse{Data: &pb.GetFileByIDResponse_ChunkData{ChunkData: buffer[:n]}}

		err = pp.Send(req)
		if err != nil {
			log.Printf("%s %s\n", err, pp.RecvMsg(nil).Error())
			return status.Errorf(codes.Internal, "%s %s", err, pp.RecvMsg(nil).Error())
		}
	}

	return err
}
