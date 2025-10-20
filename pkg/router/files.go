package router

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/argon-chat/KineticaFS/pkg/guid"
	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/argon-chat/KineticaFS/pkg/timestamp"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// AddFileRoutes sets up the server-side file management endpoints.
func AddFileRoutes(router *router, v1 *gin.RouterGroup) {
	files := v1.Group("/file")
	files.POST("/", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.InitiateFileUploadHandler)
	files.POST("/:blob/finalize", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.FinalizeFileUploadHandler)
	files.DELETE("/:id", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.DeleteFileHandler)
}

// AddFileBlobRoutes sets up the client-side upload endpoint.
func AddFileBlobRoutes(router *router, v1 *gin.RouterGroup) {
	upload := v1.Group("/upload")
	upload.PATCH("/:blob", AuthMiddleware(router.repo), router.UploadFileBlobHandler)
}

type InitiateFileUploadDTO struct {
	RegionID   string `json:"regionId" binding:"required"`
	BucketCode string `json:"bucketCode"`
}

type InitiateFileUploadResponse struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"` // seconds
}

type RegionBucket struct {
	ID       uint16 `json:"id"`
	BucketID string `json:"bucketId"`
}

type RegionInfo struct {
	ID      uint8          `json:"id"`
	Buckets []RegionBucket `json:"buckets"`
}

type Regions map[string]RegionInfo

func loadRegionsConfig(regions *Regions) error {
	path := viper.GetString("region")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, regions)
	if err != nil {
		return err
	}
	return nil
}

func generateRandomEntropy() uint64 {
	var entropy uint64
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return uint64(timestamp.CurrentTimestamp())
	}
	entropy = binary.BigEndian.Uint64(bytes)
	return entropy
}

// Initiate a new file upload (admin only)
// @Summary Initiate file upload
// @Description Initiate a new file upload. Receives regionId and bucketCode, returns a pre-signed upload URL and TTL (seconds). Admin access required.
// @Tags files
// @Accept json
// @Produce json
// @Param x-api-token header string true "API Token"
// @Param data body InitiateFileUploadDTO true "Upload initiation data"
// @Success 201 {object} InitiateFileUploadResponse
// @Failure 400 {object} router.ErrorResponse
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Router /v1/file/ [post]
func (r *router) InitiateFileUploadHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var dto InitiateFileUploadDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, ErrorResponse{Message: "Invalid request body: " + err.Error()})
		return
	}
	regions := Regions{}
	err := loadRegionsConfig(&regions)
	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Failed to load regions configuration: " + err.Error()})
		return
	}
	region, ok := regions[dto.RegionID]
	if !ok {
		c.JSON(400, ErrorResponse{Message: "Invalid region ID"})
		return
	}
	var bucketID uint16
	if dto.BucketCode == "" {
		if len(region.Buckets) == 0 {
			c.JSON(400, ErrorResponse{Message: "No buckets defined for the specified region"})
			return
		}
		randIndexBytes := make([]byte, 2)
		_, err := rand.Read(randIndexBytes)
		if err != nil {
			c.JSON(500, ErrorResponse{Message: "Failed to generate random bucket selection: " + err.Error()})
			return
		}
		randIndex := binary.BigEndian.Uint16(randIndexBytes) % uint16(len(region.Buckets))
		bucketID = region.Buckets[randIndex].ID
		dto.BucketCode = region.Buckets[randIndex].BucketID
	}
	found := false
	for _, bucket := range region.Buckets {
		if bucket.BucketID == dto.BucketCode {
			bucketID = bucket.ID
			found = true
			break
		}
	}
	if !found {
		c.JSON(400, ErrorResponse{Message: "Invalid bucket code for the specified region"})
		return
	}
	entropy := generateRandomEntropy()
	guid := guid.NewGuid(timestamp.CurrentTimestamp(), region.ID, bucketID, entropy, 0x0A)
	guidString, err := guid.Pack()
	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Failed to generate file GUID: " + err.Error()})
		return
	}

	model := &models.File{BucketID: dto.BucketCode, Name: guidString}
	blob := &models.FileBlob{FileID: guidString}

	err = r.repo.Files.CreateFile(ctx, model)
	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Failed to create file record: " + err.Error()})
		return
	}
	blob, err = r.repo.FileBlobs.CreateFileBlob(ctx, blob)
	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Failed to create file blob: " + err.Error()})
		return
	}

	response := InitiateFileUploadResponse{
		URL: blob.GetID(),
		TTL: 600,
	}
	c.JSON(201, response)
}

// Upload file data (client)
// @Summary Upload file data
// @Description Upload file data using the blob ID provided by the server. Supports stream, form-data, and multipart uploads. No admin access required.
// @Tags upload
// @Accept octet-stream
// @Accept multipart/form-data
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param x-api-token header string true "API Token"
// @Param blob path string true "Blob ID"
// @Param file formData file false "File data (multipart or form-data, required if not using raw stream)"
// @Param file body []byte false "File data (raw stream, required if not using multipart/form-data)"
// @Success 204 {object} nil
// @Failure 400 {object} router.ErrorResponse
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/upload/{blob} [patch]
func (r *router) UploadFileBlobHandler(c *gin.Context) {
	blobId := c.Param("blob")
	ctx := c.Request.Context()

	blob, err := r.repo.FileBlobs.GetFileBlobByID(ctx, blobId)
	if err != nil {
		c.JSON(404, ErrorResponse{Message: "File blob not found: " + err.Error()})
		return
	}

	file, err := r.repo.Files.GetFileByID(ctx, blob.FileID)
	if err != nil {
		c.JSON(404, ErrorResponse{Message: "File not found: " + err.Error()})
		return
	}

	bucket, err := r.repo.Buckets.GetBucketByID(ctx, file.BucketID)
	if err != nil {
		c.JSON(404, ErrorResponse{Message: "Bucket not found: " + err.Error()})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, ErrorResponse{Message: "Failed to read request body: " + err.Error()})
		return
	}

	if len(body) == 0 {
		c.JSON(400, ErrorResponse{Message: "Empty file data"})
		return
	}

	s3Client, err := createS3Client(bucket)
	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Failed to create S3 client: " + err.Error()})
		return
	}

	objectKey := file.Name

	contentType := c.GetHeader("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(bucket.Name),
		Key:           aws.String(objectKey),
		Body:          strings.NewReader(string(body)),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(body))),
	})

	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Failed to upload file to S3: " + err.Error()})
		return
	}

	file.Path = fmt.Sprintf("s3://%s/%s", bucket.Name, objectKey)
	file.FileSize = int64(len(body))
	file.ContentType = contentType
	file.Finalized = true

	err = r.repo.Files.UpdateFile(ctx, file)
	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Failed to update file record: " + err.Error()})
		return
	}
	c.Status(204)
}

func createS3Client(bucket *models.Bucket) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:           bucket.Endpoint,
				SigningRegion: bucket.Region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			bucket.AccessKey,
			bucket.SecretKey,
			"",
		)),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(bucket.Region),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return client, nil
}

// Finalize file upload (admin only)
// @Summary Finalize file upload
// @Description Finalize a file upload after client notifies server. Admin access required.
// @Tags files
// @Produce json
// @Param x-api-token header string true "API Token"
// @Param blob path string true "Blob ID"
// @Success 200 {object} models.File
// @Failure 400 {object} router.ErrorResponse
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/file/{blob}/finalize [post]
func (r *router) FinalizeFileUploadHandler(c *gin.Context) {
	blobId := c.Param("blob")
	ctx := c.Request.Context()
	log.Println(blobId)
	blob, err := r.repo.FileBlobs.GetFileBlobByID(ctx, blobId)
	if err != nil {
		c.JSON(404, ErrorResponse{Message: "File blob not found: " + err.Error()})
		return
	}

	file, err := r.repo.Files.GetFileByID(ctx, blob.FileID)
	if err != nil {
		c.JSON(404, ErrorResponse{Message: "File not found: " + err.Error()})
		return
	}
	err = r.repo.FileBlobs.DeleteFileBlobByID(ctx, blob.ID)
	if err != nil {
		c.JSON(500, ErrorResponse{Message: "Failed to delete file blob record: " + err.Error()})
		return
	}
	c.JSON(200, file)
}

// Delete file (admin only)
// @Summary Delete file
// @Description Delete a file by ID. Admin access required.
// @Tags files
// @Param x-api-token header string true "API Token"
// @Param id path string true "File ID"
// @Success 204 {object} nil
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/file/{id} [delete]
func (r *router) DeleteFileHandler(c *gin.Context) {
	// Implementation goes here
}
