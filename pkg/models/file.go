package models

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

type FileReferences struct {
	FileID   string `json:"file_id" binding:"required"`
	Metadata string `json:"metadata,omitempty"`
}

func (f File) GetID() string {
	return f.ID
}
