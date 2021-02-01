package grpcservice

import (
	"bufio"
	"bytes"
	"context"
	"failed-interview/01/internal/models"
	pb "failed-interview/01/pkg/proto"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileClient struct {
	fc pb.FileServiceClient
}

type Auth struct {
	Token string
}

func (t *Auth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"auth": t.Token,
	}, nil
}

func (c *Auth) RequireTransportSecurity() bool {
	return false
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

func NewFileClient(grpcaddress string, auth string) *FileClient {
	c, err := grpc.Dial(grpcaddress,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&Auth{auth}))
	if err != nil {
		log.Printf("cannot dial server %s err: %s\n", grpcaddress, err)
		return nil
	}

	fileClients := pb.NewFileServiceClient(c)

	return &FileClient{fc: fileClients}
}

func (f *FileClient) SaveFile(file *models.File) (id string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := f.fc.SaveFile(ctx)
	if err != nil {
		log.Println(err)
		return id, err
	}

	pbtime, err := ptypes.TimestampProto(file.DateUpload)
	if err != nil {
		log.Println(err)
		return id, err
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

	err = p.Send(req)
	if err != nil {
		log.Printf("%s %s\n", err, p.RecvMsg(nil).Error())

		return id, fmt.Errorf("%s %s", err, p.RecvMsg(nil).Error())
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
			return id, err
		}

		req := &pb.SaveFileRequest{
			Data: &pb.SaveFileRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = p.Send(req)
		if err != nil {
			log.Printf("%s %s\n", err.Error(), p.RecvMsg(nil).Error())
			return id, fmt.Errorf("%s %s", err.Error(), p.RecvMsg(nil).Error())
		}
	}

	res, err := p.CloseAndRecv()
	if err != nil {
		log.Println(err)
		return id, err
	}

	id = res.GetId()

	return id, err
}

func (f *FileClient) GetFileIDs() (ids []string) {
	p := &pb.GetFileIDsRequest{}
	ctx := context.Background()
	id, err := f.fc.GetFileIDs(ctx, p)

	if err != nil {
		log.Println(err)
		return ids
	}

	return id.Ids
}

func (f *FileClient) GetFileInfoByID(id string) (meta *models.Meta, err error) {
	p := &pb.GetFileInfoByIDRequest{Id: id}
	ctx := context.Background()
	m, err := f.fc.GetFileInfoByID(ctx, p)

	if err != nil {
		log.Println(err)
		return meta, err
	}

	meta = &models.Meta{}
	meta.DateUpload = m.FileInfo.DateUpload.AsTime()
	meta.Extension = m.FileInfo.Extension
	meta.Name = m.FileInfo.Name

	return meta, err
}

func (f *FileClient) GetFreeCapacity() uint {
	p := &pb.GetFreeCapacityRequest{}
	ctx := context.Background()
	free, err := f.fc.GetFreeCapacity(ctx, p)

	if err != nil {
		log.Println(err)
		return uint(0)
	}

	return uint(free.FreeCapacity)
}

func (f *FileClient) DeleteFile(id string) (err error) {
	p := &pb.DeleteFileRequest{Id: id}
	ctx := context.Background()
	_, err = f.fc.DeleteFile(ctx, p)

	return err
}

func (f *FileClient) GetFileByID(id string) (file *models.File, err error) {
	file = &models.File{}

	p := &pb.GetFileByIDRequest{Id: id}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := f.fc.GetFileByID(ctx, p)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	r, err := req.Recv()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	file.DateUpload = r.GetFileInfo().DateUpload.AsTime()
	file.Extension = r.GetFileInfo().Extension
	file.Name = r.GetFileInfo().Name

	fileData := bytes.Buffer{}

	for {
		err := contextError(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		req, err := req.Recv()
		if err == io.EOF {
			// приняли файл
			break
		}

		if err != nil {
			log.Println(err)
			return nil, err
		}

		chunk := req.GetChunkData()

		_, err = fileData.Write(chunk)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	file.File = fileData.Bytes()

	return file, err
}
