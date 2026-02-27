package entity

import "time"

type FileFormat string

const (
	FileFormatPDF  FileFormat = "pdf"
	FileFormatEPUB FileFormat = "epub"
)

type File struct {
	ID         int64
	BookID     int64
	Format     FileFormat
	StorageKey string
	SizeBytes  int64
	UploadedAt time.Time
	Checksum   string
}
