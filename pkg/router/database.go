package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gin-gonic/gin"
)

type DatabaseExport struct {
	ServiceTokens []*models.ServiceToken `json:"service_tokens"`
	Buckets       []*models.Bucket       `json:"buckets"`
	Files         []*models.File         `json:"files"`
	FileBlobs     []*models.FileBlob     `json:"file_blobs"`
	ExportVersion string                 `json:"export_version"`
}

type DatabaseRestoreRequest struct {
	Data          DatabaseExport `json:"data" binding:"required"`
	ClearExisting bool           `json:"clear_existing"`
}

func AddDatabaseRoutes(router *router, v1 *gin.RouterGroup) {
	db := v1.Group("/database")
	db.GET("/export", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.ExportDatabaseHandler)
	db.POST("/restore", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.RestoreDatabaseHandler)
}

// ExportDatabaseHandler exports the entire database
// @Summary Export database
// @Description Export all data from the database including service tokens, buckets, files, and file blobs (admin only).
// @Tags database
// @Param x-api-token header string true "API Token"
// @Produce json
// @Success 200 {object} DatabaseExport
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 500 {object} router.ErrorResponse "Internal server error"
// @Router /api/v1/database/export [get]
// @Id ExportDatabase
func (r *router) ExportDatabaseHandler(c *gin.Context) {
	ctx := c.Request.Context()

	serviceTokens, err := r.repo.ServiceTokens.GetAllServiceTokens(ctx)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to fetch service tokens: %v", err))
		return
	}

	buckets, err := r.repo.Buckets.ListBuckets(ctx)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to fetch buckets: %v", err))
		return
	}

	var allFiles []*models.File
	for _, bucket := range buckets {
		files, err := r.repo.Files.ListFiles(ctx, bucket.ID)
		if err != nil {
			writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to fetch files for bucket %s: %v", bucket.ID, err))
			return
		}
		allFiles = append(allFiles, files...)
	}

	fileBlobs, err := r.repo.FileBlobs.GetAllFileBlobs(ctx)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to fetch file blobs: %v", err))
		return
	}

	export := DatabaseExport{
		ServiceTokens: serviceTokens,
		Buckets:       buckets,
		Files:         allFiles,
		FileBlobs:     fileBlobs,
		ExportVersion: "1.0.0",
	}

	c.JSON(http.StatusOK, export)
}

// RestoreDatabaseHandler restores database from backup
// @Summary Restore database
// @Description Restore database from a previously exported backup. Optionally clear existing data before restore (admin only).
// @Tags database
// @Param x-api-token header string true "API Token"
// @Accept json
// @Produce json
// @Param request body DatabaseRestoreRequest true "Restore Request"
// @Success 200 {object} map[string]interface{} "Restore statistics"
// @Failure 400 {object} router.ErrorResponse "Invalid request"
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 500 {object} router.ErrorResponse "Internal server error"
// @Router /api/v1/database/restore [post]
// @Id RestoreDatabase
func (r *router) RestoreDatabaseHandler(c *gin.Context) {
	var req DatabaseRestoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if req.Data.ExportVersion == "" {
		writeError(c, http.StatusBadRequest, "missing or invalid 'data' field in request")
		return
	}

	ctx := c.Request.Context()

	if req.ClearExisting {
		if err := r.repo.ClearAllData(ctx); err != nil {
			writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to clear existing data: %v", err))
			return
		}
	}

	stats := map[string]interface{}{
		"service_tokens_restored": 0,
		"buckets_restored":        0,
		"files_restored":          0,
		"file_blobs_restored":     0,
		"errors":                  []string{},
	}

	existingTokens, err := r.repo.ServiceTokens.GetAllServiceTokens(ctx)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to fetch existing tokens: %v", err))
		return
	}
	for _, existingToken := range existingTokens {
		if err := r.repo.ServiceTokens.RevokeServiceToken(ctx, existingToken.ID); err != nil {
			errorMsg := fmt.Sprintf("failed to remove existing token %s: %v", existingToken.ID, err)
			stats["errors"] = append(stats["errors"].([]string), errorMsg)
		}
	}

	for _, token := range req.Data.ServiceTokens {
		if err := r.repo.ServiceTokens.CreateServiceToken(ctx, token); err != nil {
			errorMsg := fmt.Sprintf("failed to restore service token %s: %v", token.ID, err)
			stats["errors"] = append(stats["errors"].([]string), errorMsg)
		} else {
			stats["service_tokens_restored"] = stats["service_tokens_restored"].(int) + 1
		}
	}

	for _, bucket := range req.Data.Buckets {
		if err := r.repo.Buckets.CreateBucket(ctx, bucket); err != nil {
			errorMsg := fmt.Sprintf("failed to restore bucket %s: %v", bucket.ID, err)
			stats["errors"] = append(stats["errors"].([]string), errorMsg)
		} else {
			stats["buckets_restored"] = stats["buckets_restored"].(int) + 1
		}
	}

	for _, file := range req.Data.Files {
		if err := r.repo.Files.CreateFile(ctx, file); err != nil {
			errorMsg := fmt.Sprintf("failed to restore file %s: %v", file.ID, err)
			stats["errors"] = append(stats["errors"].([]string), errorMsg)
		} else {
			stats["files_restored"] = stats["files_restored"].(int) + 1
		}
	}

	for _, blob := range req.Data.FileBlobs {
		if _, err := r.repo.FileBlobs.CreateFileBlob(ctx, blob); err != nil {
			errorMsg := fmt.Sprintf("failed to restore file blob %s: %v", blob.ID, err)
			stats["errors"] = append(stats["errors"].([]string), errorMsg)
		} else {
			stats["file_blobs_restored"] = stats["file_blobs_restored"].(int) + 1
		}
	}

	statsJSON, _ := json.Marshal(stats)
	var result map[string]interface{}
	json.Unmarshal(statsJSON, &result)

	c.JSON(http.StatusOK, result)
}
