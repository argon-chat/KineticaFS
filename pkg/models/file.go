package models

type IFileRepository interface {
	GetFileByID(id string) (*File, error)
	GetFileByName(bucketID, name string) (*File, error)
	CreateFile(file *File) error
	UpdateFile(file *File) error
	DeleteFile(id string) error
	ListFiles(bucketID string) ([]*File, error)
}

type File struct {
	ApplicationModel
	BucketID    string `json:"bucket_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Path        string `json:"path"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
	Checksum    string `json:"checksum"`
	Metadata    string `json:"metadata,omitempty"`
}

func (f File) GetID() string {
	return f.ID
}
