package grpcservice

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"testing"

	"failed-interview/01/internal/models"
	pb "failed-interview/01/pkg/proto"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type mockFileServer struct {
	pb.UnimplementedFileServiceServer
}

func newClient() *FileClient {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &mockFileServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}

	return &FileClient{fc: pb.NewFileServiceClient(conn)}
}

func (m *mockFileServer) SaveFile(p pb.FileService_SaveFileServer) error {
	req, err := p.Recv()

	if err != nil {
		log.Println(err)
		return status.Error(codes.Unknown, err.Error())
	}

	f := &models.File{Name: req.GetFileInfo().Name, Extension: req.GetFileInfo().Extension, DateUpload: req.GetFileInfo().DateUpload.AsTime()}

	fileData := bytes.Buffer{}

	for {
		err := contextError(p.Context())
		if err != nil {
			return err
		}

		req, err := p.Recv()
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

	res := &pb.SaveFileResponse{Id: "id"}

	err = p.SendAndClose(res)
	if err != nil {
		log.Println(err)
		return status.Errorf(codes.Unknown, err.Error())
	}

	return nil
}

func (m *mockFileServer) GetFileInfoByID(context.Context, *pb.GetFileInfoByIDRequest) (*pb.GetFileInfoByIDResponse, error) {
	return &pb.GetFileInfoByIDResponse{FileInfo: &pb.FileInfo{}}, nil
}

func (m *mockFileServer) GetFileIDs(context.Context, *pb.GetFileIDsRequest) (*pb.GetFileIDsResponse, error) {
	return &pb.GetFileIDsResponse{Ids: []string{"file"}}, nil
}

func (m *mockFileServer) GetFreeCapacity(context.Context, *pb.GetFreeCapacityRequest) (*pb.GetFreeCapacityResponse, error) {
	return &pb.GetFreeCapacityResponse{}, nil
}

func (m *mockFileServer) DeleteFile(context.Context, *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	return &pb.DeleteFileResponse{}, nil
}

func (m *mockFileServer) GetFileByID(p *pb.GetFileByIDRequest, pp pb.FileService_GetFileByIDServer) error {
	err := pp.Send(&pb.GetFileByIDResponse{Data: &pb.GetFileByIDResponse_FileInfo{FileInfo: &pb.FileInfo{}}})
	if err != nil {
		log.Println(err)
		return status.Error(codes.Internal, err.Error())
	}

	r := bytes.NewReader(make([]byte, 100))

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

	return nil
}

func TestSaveFile(t *testing.T) {
	client := newClient()
	id, err := client.SaveFile(&models.File{File: make([]byte, 100)})
	require.NoError(t, err)
	require.Equal(t, "id", id)
}

func TestGetFileIDs(t *testing.T) {
	client := newClient()
	id := client.GetFileIDs()

	require.Equal(t, []string{"file"}, id)
}

func TestGetFileInfoByID(t *testing.T) {
	client := newClient()
	_, err := client.GetFileInfoByID("file")

	require.NoError(t, err)
}

func TestGetFreeCapacity(t *testing.T) {
	client := newClient()
	free := client.GetFreeCapacity()

	require.Equal(t, uint(0), free)
}

func TestDeleteFile(t *testing.T) {
	client := newClient()
	err := client.DeleteFile("file")

	require.NoError(t, err)
}

func TestGetFileByID(t *testing.T) {
	client := newClient()
	_, err := client.GetFileByID("file")

	require.NoError(t, err)
}
