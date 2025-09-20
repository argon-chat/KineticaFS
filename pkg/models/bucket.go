package models

type Bucket struct {
	ApplicationModel
	Name         string `json:"name" binding:"required" gorm:"uniqueIndex"`
	Region       string `json:"region" binding:"required"`
	Endpoint     string `json:"endpoint" binding:"required"`
	AccessKey    string `json:"access_key" binding:"required"`
	SecretKey    string `json:"secret_key" binding:"required"`
	UseSSL       bool   `json:"use_ssl"`
	S3Provider   string `json:"s3_provider"`
	CustomConfig string `json:"custom_config,omitempty"`
}

func (bu Bucket) GetID() string {
	return bu.ID
}
