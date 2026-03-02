package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	filepb "github.com/shadowpr1est/knigapoisk-file-service/api/proto"
	"github.com/shadowpr1est/knigapoisk-file-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-file-service/internal/usecase/file"
)

type FileServer struct {
	filepb.UnimplementedFileServiceServer
	useCase file.UseCase
	logger  *zap.Logger
}

func NewFileServer(uc file.UseCase, logger *zap.Logger) *FileServer {
	return &FileServer{
		useCase: uc,
		logger:  logger,
	}
}

func (s *FileServer) UploadFile(ctx context.Context, req *filepb.UploadFileRequest) (*filepb.UploadFileResponse, error) {
	input := file.UploadFileInput{
		BookID: req.BookId,
		Format: fromProtoFormat(req.Format),
		Data:   req.Data,
	}
	out, err := s.useCase.UploadFile(ctx, input)
	if err != nil {
		s.logger.Error("UploadFile error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	f := out.File
	return &filepb.UploadFileResponse{
		File: &filepb.File{
			Id:         f.ID,
			BookId:     f.BookID,
			Format:     toProtoFormat(f.Format),
			StorageKey: f.StorageKey,
			SizeBytes:  f.SizeBytes,
			UploadedAt: f.UploadedAt.Format(time.RFC3339),
			Checksum:   f.Checksum,
		},
	}, nil
}

func (s *FileServer) GetFile(ctx context.Context, req *filepb.GetFileRequest) (*filepb.GetFileResponse, error) {
	out, err := s.useCase.GetFile(ctx, req.Id)
	if err != nil {
		s.logger.Error("GetFile error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	if out.File == nil {
		return nil, status.Error(codes.NotFound, "file not found")
	}
	f := out.File
	return &filepb.GetFileResponse{
		File: &filepb.File{
			Id:         f.ID,
			BookId:     f.BookID,
			Format:     toProtoFormat(f.Format),
			StorageKey: f.StorageKey,
			SizeBytes:  f.SizeBytes,
			UploadedAt: f.UploadedAt.Format(time.RFC3339),
			Checksum:   f.Checksum,
		},
		Data: out.Data,
	}, nil
}

func (s *FileServer) DeleteFile(ctx context.Context, req *filepb.DeleteFileRequest) (*filepb.DeleteFileResponse, error) {
	if err := s.useCase.DeleteFile(ctx, req.Id); err != nil {
		s.logger.Error("DeleteFile error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &filepb.DeleteFileResponse{Success: true}, nil
}

func (s *FileServer) Health(ctx context.Context, _ *filepb.HealthRequest) (*filepb.HealthResponse, error) {
	return &filepb.HealthResponse{Status: "ok"}, nil
}

func fromProtoFormat(f filepb.FileFormat) entity.FileFormat {
	switch f {
	case filepb.FileFormat_FILE_FORMAT_PDF:
		return entity.FileFormatPDF
	case filepb.FileFormat_FILE_FORMAT_EPUB:
		return entity.FileFormatEPUB
	default:
		return ""
	}
}

func toProtoFormat(f entity.FileFormat) filepb.FileFormat {
	switch f {
	case entity.FileFormatPDF:
		return filepb.FileFormat_FILE_FORMAT_PDF
	case entity.FileFormatEPUB:
		return filepb.FileFormat_FILE_FORMAT_EPUB
	default:
		return filepb.FileFormat_FILE_FORMAT_UNSPECIFIED
	}
}

