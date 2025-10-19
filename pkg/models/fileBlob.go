package models

type FileBlob struct {
	ApplicationModel
	FileID string `json:"file_id" binding:"required"`
}

func (fb FileBlob) GetID() string {
	return fb.ID
}
