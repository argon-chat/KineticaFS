package router

import (
	"fmt"
	"net/http"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gin-gonic/gin"
)

// AddBucketsRoutes sets up the bucket endpoints.
func AddBucketsRoutes(router *router, v1 *gin.RouterGroup) {
	bucket := v1.Group("/bucket")
	bucket.POST("/", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.CreateBucketHandler)
	bucket.GET("/", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.ListBucketsHandler)
	bucket.GET("/:id", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.GetBucketHandler)
	bucket.PATCH("/:id", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.UpdateBucketHandler)
	bucket.DELETE("/:id", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.DeleteBucketHandler)
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
// @Router /api/v1/bucket/ [post]
// @Id CreateBucket
func (r *router) CreateBucketHandler(c *gin.Context) {
	var bucket models.Bucket
	if err := c.ShouldBindJSON(&bucket); err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}
	ctx := c.Request.Context()
	existing, err := r.repo.Buckets.GetBucketByName(ctx, bucket.Name)
	if err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("failed to check existing bucket: %v", err))
		return
	}
	if existing != nil {
		writeError(c, http.StatusConflict, "Bucket with the same name already exists")
		return
	}
	err = r.repo.Buckets.CreateBucket(ctx, &bucket)
	if err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("failed to create bucket: %v", err))
		return
	}
	c.JSON(http.StatusCreated, bucket)
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
// @Router /api/v1/bucket/ [get]
// @Id ListBuckets
func (r *router) ListBucketsHandler(c *gin.Context) {
	buckets, err := r.repo.Buckets.ListBuckets(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to list buckets: %v", err))
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
// @Router /api/v1/bucket/{id} [get]
// @Id GetBucket
func (r *router) GetBucketHandler(c *gin.Context) {
	id := c.Param("id")
	bucket, err := r.repo.Buckets.GetBucketByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to get bucket: %v", err))
		return
	}
	if bucket == nil {
		writeError(c, http.StatusNotFound, "Bucket not found")
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
// @Router /api/v1/bucket/{id} [patch]
// @Id UpdateBucket
func (r *router) UpdateBucketHandler(c *gin.Context) {
	id := c.Param("id")
	var req models.Bucket
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	ctx := c.Request.Context()
	bucket, err := r.repo.Buckets.GetBucketByID(ctx, id)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to get bucket: %v", err))
		return
	}
	if bucket == nil {
		writeError(c, http.StatusNotFound, "Bucket not found")
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

	if err := r.repo.Buckets.UpdateBucket(ctx, bucket); err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("failed to update bucket: %v", err))
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
// @Router /api/v1/bucket/{id} [delete]
// @Id DeleteBucket
func (r *router) DeleteBucketHandler(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	bucket, err := r.repo.Buckets.GetBucketByID(ctx, id)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to get bucket: %v", err))
		return
	}
	if bucket == nil {
		writeError(c, http.StatusNotFound, "Bucket not found")
		return
	}
	err = r.repo.Buckets.DeleteBucket(ctx, id)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to delete bucket: %v", err))
		return
	}
	c.Status(204)
}
