package models

type File struct {
	ApplicationModel
	BucketID      string  `json:"bucket_id" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Path          string  `json:"path"`
	FileSize      int64   `json:"file_size"`
	ContentType   string  `json:"content_type"`
	Checksum      string  `json:"checksum"`
	Finalized     bool    `json:"finalized"`
	FileSizeLimit uint64  `json:"file_size_limit"`
	References    int64   `json:"references"`
	Metadata      string  `json:"metadata,omitempty"`
	UserSub       *string `json:"user_sub,omitempty"`     // GUID - user identifier from JWT
	SpaceId       *string `json:"space_id,omitempty"`     // GUID nullable - space identifier from JWT
	FileType      *string `json:"file_type,omitempty"`    // string - type/scope from JWT (avatar, profileHeader, etc.)
}

func (f File) GetID() string {
	return f.ID
}
