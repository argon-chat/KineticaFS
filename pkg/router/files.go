package router

import "github.com/gin-gonic/gin"

// AddFileRoutes sets up the server-side file management endpoints.
func AddFileRoutes(router *router, v1 *gin.RouterGroup) {
	files := v1.Group("/file")
	files.POST("/", AuthMiddleware(router.repo), AdminOnlyMiddleware, InitiateFileUploadHandler)
	files.POST("/:id/finalize", AuthMiddleware(router.repo), AdminOnlyMiddleware, FinalizeFileUploadHandler)
	files.DELETE("/:id", AuthMiddleware(router.repo), AdminOnlyMiddleware, DeleteFileHandler)
}

// AddFileBlobRoutes sets up the client-side upload endpoint.
func AddFileBlobRoutes(router *router, v1 *gin.RouterGroup) {
	upload := v1.Group("/upload")
	upload.PATCH("/:blob", AuthMiddleware(router.repo), UploadFileBlobHandler)
}

type InitiateFileUploadDTO struct {
	RegionID   string `json:"regionId" binding:"required"`
	BucketCode string `json:"bucketCode" binding:"required"`
}

type InitiateFileUploadResponse struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"` // seconds
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
func InitiateFileUploadHandler(c *gin.Context) {

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
func UploadFileBlobHandler(c *gin.Context) {
	// Implementation goes here
}

// Finalize file upload (admin only)
// @Summary Finalize file upload
// @Description Finalize a file upload after client notifies server. Admin access required.
// @Tags files
// @Produce json
// @Param x-api-token header string true "API Token"
// @Param id path string true "File ID"
// @Success 200 {object} models.File
// @Failure 400 {object} router.ErrorResponse
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/file/{id}/finalize [post]
func FinalizeFileUploadHandler(c *gin.Context) {
	// Implementation goes here
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
func DeleteFileHandler(c *gin.Context) {
	// Implementation goes here
}
