package router

import (
	"fmt"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gin-gonic/gin"
)

// AddBucketsRoutes sets up the bucket endpoints.
func AddBucketsRoutes(v1 *gin.RouterGroup) {
	bucket := v1.Group("/bucket")
	bucket.POST("/", AuthMiddleware, AdminOnlyMiddleware, CreateBucketHandler)
	bucket.GET("/", AuthMiddleware, AdminOnlyMiddleware, ListBucketsHandler)
	bucket.GET("/:id", AuthMiddleware, AdminOnlyMiddleware, GetBucketHandler)
	bucket.PATCH("/:id", AuthMiddleware, AdminOnlyMiddleware, UpdateBucketHandler)
	bucket.DELETE("/:id", AuthMiddleware, AdminOnlyMiddleware, DeleteBucketHandler)
}

type BucketInsertDTO struct {
	Name         string             `json:"name" binding:"required" gorm:"uniqueIndex"`
	Region       string             `json:"region" binding:"required"`
	Endpoint     string             `json:"endpoint" binding:"required"`
	AccessKey    string             `json:"access_key" binding:"required"`
	SecretKey    string             `json:"secret_key" binding:"required"`
	UseSSL       bool               `json:"use_ssl"`
	S3Provider   string             `json:"s3_provider"`
	CustomConfig string             `json:"custom_config,omitempty"`
	StorageType  models.StorageType `json:"storage_type" gorm:"default:0"`
}

// CreateBucketHandler creates a new bucket
// @Summary Create bucket
// @Description Create a new S3 bucket. Only admin users can create buckets.
// @Tags buckets
// @Param x-api-token header string true "API Token"
// @Accept json
// @Produce json
// @Param bucket body BucketInsertDTO true "Bucket"
// @Success 201 {object} models.Bucket
// @Failure 400 {object} router.ErrorResponse
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Router /v1/bucket/ [post]
func CreateBucketHandler(c *gin.Context) {
	var bucket models.Bucket
	if err := c.ShouldBindJSON(&bucket); err != nil {
		c.JSON(400, ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("invalid request body: %v", err),
		})
		return
	}
	existing, err := applicationRepository.Buckets.GetBucketByName(bucket.Name)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("failed to check existing bucket: %v", err),
		})
		return
	}
	if existing != nil {
		c.JSON(409, ErrorResponse{
			Code:    409,
			Message: "Bucket with the same name already exists",
		})
		return
	}
	err = applicationRepository.Buckets.CreateBucket(&bucket)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("failed to create bucket: %v", err),
		})
		return
	}
	c.JSON(201, bucket)
}

// ListBucketsHandler lists all buckets
// @Summary List buckets
// @Description List all S3 buckets
// @Tags buckets
// @Param x-api-token header string true "API Token"
// @Produce json
// @Success 200 {array} models.Bucket
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Router /v1/bucket/ [get]
func ListBucketsHandler(c *gin.Context) {
	buckets, err := applicationRepository.Buckets.ListBuckets()
	if err != nil {
		c.JSON(500, ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to list buckets: %v", err),
		})
		return
	}
	c.JSON(200, buckets)
}

// GetBucketHandler gets a bucket by ID
// @Summary Get bucket
// @Description Get a bucket by ID
// @Tags buckets
// @Param x-api-token header string true "API Token"
// @Produce json
// @Param id path string true "Bucket ID"
// @Success 200 {object} models.Bucket
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/bucket/{id} [get]
func GetBucketHandler(c *gin.Context) {
	id := c.Param("id")
	bucket, err := applicationRepository.Buckets.GetBucketByID(id)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to get bucket: %v", err),
		})
		return
	}
	if bucket == nil {
		c.JSON(404, ErrorResponse{
			Code:    404,
			Message: "Bucket not found",
		})
		return
	}
	c.JSON(200, bucket)
}

// UpdateBucketHandler updates a bucket by ID
// @Summary Update bucket
// @Description Update a bucket by ID. Only admin users can update buckets.
// @Tags buckets
// @Param x-api-token header string true "API Token"
// @Accept json
// @Produce json
// @Param id path string true "Bucket ID"
// @Param bucket body BucketInsertDTO true "Bucket"
// @Success 200 {object} models.Bucket
// @Failure 400 {object} router.ErrorResponse
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/bucket/{id} [patch]
func UpdateBucketHandler(c *gin.Context) {
	id := c.Param("id")
	var req models.Bucket
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("invalid request body: %v", err),
		})
		return
	}
	bucket, err := applicationRepository.Buckets.GetBucketByID(id)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to get bucket: %v", err),
		})
		return
	}
	if bucket == nil {
		c.JSON(404, ErrorResponse{
			Code:    404,
			Message: "Bucket not found",
		})
		return
	}
	bucket.Name = req.Name
	bucket.Region = req.Region
	bucket.Endpoint = req.Endpoint
	bucket.AccessKey = req.AccessKey
	bucket.SecretKey = req.SecretKey
	bucket.UseSSL = req.UseSSL
	bucket.S3Provider = req.S3Provider
	bucket.CustomConfig = req.CustomConfig
	bucket.StorageType = req.StorageType
	err = applicationRepository.Buckets.UpdateBucket(bucket)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("failed to update bucket: %v", err),
		})
		return
	}
	c.JSON(200, bucket)
}

// DeleteBucketHandler deletes a bucket by ID
// @Summary Delete bucket
// @Description Delete a bucket by ID. Only admin users can delete buckets.
// @Tags buckets
// @Param x-api-token header string true "API Token"
// @Param id path string true "Bucket ID"
// @Success 204 {object} nil
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/bucket/{id} [delete]
func DeleteBucketHandler(c *gin.Context) {
	id := c.Param("id")
	bucket, err := applicationRepository.Buckets.GetBucketByID(id)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to get bucket: %v", err),
		})
		return
	}
	if bucket == nil {
		c.JSON(404, ErrorResponse{
			Code:    404,
			Message: "Bucket not found",
		})
		return
	}
	err = applicationRepository.Buckets.DeleteBucket(id)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to delete bucket: %v", err),
		})
		return
	}
	c.Status(204)
}
