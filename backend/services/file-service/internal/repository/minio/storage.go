package minio

import (
	"bytes"
	"context"

	minioapi "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/shadowpr1est/knigapoisk-file-service/internal/domain/repository"
)

type Storage struct {
	client *minioapi.Client
	bucket string
}

func NewStorage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (repository.StorageRepository, error) {
	client, err := minioapi.New(endpoint, &minioapi.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &Storage{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *Storage) Upload(ctx context.Context, key string, data []byte) error {
	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(ctx, s.bucket, key, reader, int64(len(data)), minioapi.PutObjectOptions{})
	return err
}

func (s *Storage) Download(ctx context.Context, key string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minioapi.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucket, key, minioapi.RemoveObjectOptions{})
}

