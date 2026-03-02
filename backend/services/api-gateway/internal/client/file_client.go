package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	filepb "github.com/shadowpr1est/knigapoisk-file-service/api/proto"
)

type FileClient interface {
	UploadFile(ctx context.Context, req *filepb.UploadFileRequest) (*filepb.UploadFileResponse, error)
	GetFile(ctx context.Context, req *filepb.GetFileRequest) (*filepb.GetFileResponse, error)
	DeleteFile(ctx context.Context, req *filepb.DeleteFileRequest) (*filepb.DeleteFileResponse, error)
}

type fileClient struct {
	cc filepb.FileServiceClient
}

func NewFileClient(addr string) (FileClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &fileClient{cc: filepb.NewFileServiceClient(conn)}, nil
}

func (c *fileClient) UploadFile(ctx context.Context, req *filepb.UploadFileRequest) (*filepb.UploadFileResponse, error) {
	return c.cc.UploadFile(ctx, req)
}

func (c *fileClient) GetFile(ctx context.Context, req *filepb.GetFileRequest) (*filepb.GetFileResponse, error) {
	return c.cc.GetFile(ctx, req)
}

func (c *fileClient) DeleteFile(ctx context.Context, req *filepb.DeleteFileRequest) (*filepb.DeleteFileResponse, error) {
	return c.cc.DeleteFile(ctx, req)
}

