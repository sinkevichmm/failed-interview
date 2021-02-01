package grpc

import (
	"bufio"
	"bytes"
	"context"
	m "failed-interview/01/internal/manager/file/repository/mockfilestorage"
	u "failed-interview/01/internal/manager/file/usecase"
	"failed-interview/01/internal/models"
	pb "failed-interview/01/pkg/proto"
	"io"
	"log"
	"net"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func newClient() (pb.FileServiceClient, context.Context) {
	lis = bufconn.Listen(bufSize)

	uc := &m.Mock{}
	fsRepo := u.NewFileUseCase(uc)
	////////////////////////////////////////////////////////
	uc.On("GetFileIDs").Return([]string(nil))

	b := make([]byte, 1000)
	file := &models.File{File: b}
	uc.On("SaveFile", file).Return("uuid", nil)

	uc.On("GetFileByID", "file").Return(file, nil)
	uc.On("GetFileInfoByID", "file").Return(&models.Meta{Name: "file"}, nil)

	uc.On("GetFreeCapacity").Return(uint(7))

	uc.On("DeleteFile", "file").Return(nil)

	////////////////////////////////////////////////////////
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &FileServer{UseCase: fsRepo})

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
	//defer conn.Close()

	return pb.NewFileServiceClient(conn), ctx
}

func TestGetFileIDs(t *testing.T) {
	client, ctx := newClient()

	resp, err := client.GetFileIDs(ctx, &pb.GetFileIDsRequest{})
	if err != nil {
		t.Fatalf("GetFileIDs failed: %v", err)
	}

	require.Equal(t, []string(nil), resp.Ids)
}

func TestSaveFile(t *testing.T) {
	client, ctx := newClient()

	resp, err := client.SaveFile(ctx)
	if err != nil {
		t.Fatalf("GetFileIDs failed: %v", err)
	}

	b := make([]byte, 1000)
	file := &models.File{File: b}

	pbtime, err := ptypes.TimestampProto(file.DateUpload)

	if err != nil {
		t.Fatalf("SaveFile err: %v", err)
	}

	req := &pb.SaveFileRequest{

		Data: &pb.SaveFileRequest_FileInfo{
			FileInfo: &pb.FileInfo{
				Name:       file.Name,
				Extension:  file.Extension,
				DateUpload: pbtime,
			},
		},
	}
	err = resp.Send(req)

	if err != nil {
		t.Fatalf("SaveFile err: %v", err)
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
			t.Fatalf("cannot read chunk to buffer: %s", err.Error())
		}

		req := &pb.SaveFileRequest{
			Data: &pb.SaveFileRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = resp.Send(req)
		if err != nil {
			t.Fatalf("cannot send chunk to server: %s %s", err.Error(), resp.RecvMsg(nil).Error())
		}
	}

	res, err := resp.CloseAndRecv()
	if err != nil {
		t.Fatalf("SaveFile err: %v", err)
	}

	require.Equal(t, "uuid", res.Id)
}

func TestGetFileByID(t *testing.T) {
	client, ctx := newClient()

	p := &pb.GetFileByIDRequest{Id: "file"}
	resp, err := client.GetFileByID(ctx, p)

	if err != nil {
		t.Fatalf("GetFileIDs failed: %v", err)
	}

	r, err := resp.Recv()

	if err != nil {
		t.Fatalf("GetFileIDs failed: %v", err)
	}

	file := &models.File{}
	file.DateUpload = r.GetFileInfo().DateUpload.AsTime()
	file.Extension = r.GetFileInfo().Extension
	file.Name = r.GetFileInfo().Name

	fileData := bytes.Buffer{}

	for {
		err := contextError(ctx)
		if err != nil {
			t.Fatalf("GetFileIDs failed: %v", err)
		}

		req, err := resp.Recv()
		if err == io.EOF {
			// приняли файл
			break
		}

		if err != nil {
			t.Fatalf("GetFileIDs failed: %v", err)
		}

		chunk := req.GetChunkData()

		_, err = fileData.Write(chunk)
		if err != nil {
			t.Fatalf("GetFileIDs failed: %v", err)
		}
	}

	fileData.Bytes()

	require.Equal(t, make([]byte, 1000), fileData.Bytes())
}

func TestGetFileInfoByID(t *testing.T) {
	client, ctx := newClient()

	resp, err := client.GetFileInfoByID(ctx, &pb.GetFileInfoByIDRequest{Id: "file"})
	if err != nil {
		t.Fatalf("GetFileIDs failed: %v", err)
	}

	require.Equal(t, "file", resp.FileInfo.Name)
}

func TestGetFreeCapacity(t *testing.T) {
	client, ctx := newClient()

	resp, err := client.GetFreeCapacity(ctx, &pb.GetFreeCapacityRequest{})
	if err != nil {
		t.Fatalf("GetFileIDs failed: %v", err)
	}

	require.Equal(t, uint32(7), resp.FreeCapacity)
}

func TestDeleteFile(t *testing.T) {
	client, ctx := newClient()

	_, err := client.DeleteFile(ctx, &pb.DeleteFileRequest{Id: "file"})
	require.NoError(t, err)
}
