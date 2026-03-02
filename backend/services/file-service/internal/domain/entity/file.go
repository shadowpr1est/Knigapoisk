package entity

import "time"

type FileFormat string

const (
	FileFormatPDF  FileFormat = "pdf"
	FileFormatEPUB FileFormat = "epub"
)

type File struct {
	ID         int64      `db:"id"`
	BookID     int64      `db:"book_id"`
	Format     FileFormat `db:"format"`
	StorageKey string     `db:"storage_key"`
	SizeBytes  int64      `db:"size_bytes"`
	UploadedAt time.Time  `db:"uploaded_at"`
	Checksum   string     `db:"checksum"`
}

