package file

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/shadowpr1est/knigapoisk-file-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-file-service/internal/domain/repository"
)

type UseCase interface {
	UploadFile(ctx context.Context, input UploadFileInput) (UploadFileOutput, error)
	GetFile(ctx context.Context, id int64) (GetFileOutput, error)
	DeleteFile(ctx context.Context, id int64) error
}

type FileUseCase struct {
	fileRepo    repository.FileRepository
	storageRepo repository.StorageRepository
}

func NewFileUseCase(
	fileRepo repository.FileRepository,
	storageRepo repository.StorageRepository,
) UseCase {
	return &FileUseCase{
		fileRepo:    fileRepo,
		storageRepo: storageRepo,
	}
}

type UploadFileInput struct {
	BookID int64
	Format entity.FileFormat
	Data   []byte
}

type UploadFileOutput struct {
	File *entity.File
}

type GetFileOutput struct {
	File *entity.File
	Data []byte
}

func (u *FileUseCase) UploadFile(ctx context.Context, input UploadFileInput) (UploadFileOutput, error) {
	now := time.Now()
	hash := sha256.Sum256(input.Data)
	checksum := hex.EncodeToString(hash[:])

	storageKey := generateStorageKey(input.BookID, now, checksum)

	if err := u.storageRepo.Upload(ctx, storageKey, input.Data); err != nil {
		return UploadFileOutput{}, err
	}

	file := &entity.File{
		BookID:     input.BookID,
		Format:     input.Format,
		StorageKey: storageKey,
		SizeBytes:  int64(len(input.Data)),
		UploadedAt: now,
		Checksum:   checksum,
	}
	if err := u.fileRepo.Create(ctx, file); err != nil {
		return UploadFileOutput{}, err
	}

	return UploadFileOutput{File: file}, nil
}

func (u *FileUseCase) GetFile(ctx context.Context, id int64) (GetFileOutput, error) {
	file, err := u.fileRepo.GetByID(ctx, id)
	if err != nil {
		return GetFileOutput{}, err
	}
	if file == nil {
		return GetFileOutput{}, nil
	}
	data, err := u.storageRepo.Download(ctx, file.StorageKey)
	if err != nil {
		return GetFileOutput{}, err
	}
	return GetFileOutput{
		File: file,
		Data: data,
	}, nil
}

func (u *FileUseCase) DeleteFile(ctx context.Context, id int64) error {
	file, err := u.fileRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}
	if err := u.storageRepo.Delete(ctx, file.StorageKey); err != nil {
		return err
	}
	return u.fileRepo.Delete(ctx, id)
}

func generateStorageKey(bookID int64, now time.Time, checksum string) string {
	return fmt.Sprintf("books/%d/%d_%s", bookID, now.UnixNano(), checksum)
}

