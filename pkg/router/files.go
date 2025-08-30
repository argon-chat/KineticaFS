package router

import "github.com/gin-gonic/gin"

// AddFileRoutes sets up the server-side file management endpoints.
func AddFileRoutes(v1 *gin.RouterGroup) {
	files := v1.Group("/file")
	files.POST("/", InitiateFileUploadHandler)
	files.POST("/:id/finalize", FinalizeFileUploadHandler)
	files.DELETE("/:id", DeleteFileHandler)
}

// AddFileBlobRoutes sets up the client-side upload endpoint.
func AddFileBlobRoutes(v1 *gin.RouterGroup) {
	upload := v1.Group("/upload")
	upload.PATCH("/:blob", UploadFileBlobHandler)
}

// Initiate a new file upload (admin only)
// @Summary Initiate file upload
// @Description Initiate a new file upload. Returns a blob ID for the client to upload data. Admin access required.
// @Tags files
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string "{blob: blob_id}"
// @Failure 400 {object} router.ErrorResponse
// @Failure 403 {object} router.ErrorResponse
// @Router /v1/file/ [post]
func InitiateFileUploadHandler(c *gin.Context) {
	// Implementation goes here
}

// Upload file data (client)
// @Summary Upload file data
// @Description Upload file data using the blob ID provided by the server. Supports stream, form-data, and multipart uploads. No admin access required.
// @Tags upload
// @Accept octet-stream
// @Accept multipart/form-data
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param blob path string true "Blob ID"
// @Param file formData file false "File data (multipart or form-data, required if not using raw stream)"
// @Param file body []byte false "File data (raw stream, required if not using multipart/form-data)"
// @Success 204 {object} nil
// @Failure 400 {object} router.ErrorResponse
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
// @Param id path string true "File ID"
// @Success 200 {object} models.File
// @Failure 400 {object} router.ErrorResponse
// @Failure 403 {object} router.ErrorResponse
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/file/{id}/finalize [post]
func FinalizeFileUploadHandler(c *gin.Context) {
	// Implementation goes here
}

// Delete file (admin only)
// @Summary Delete file
// @Description Delete a file by ID. Admin access required.
// @Tags files
// @Param id path string true "File ID"
// @Success 204 {object} nil
// @Failure 403 {object} router.ErrorResponse
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/file/{id} [delete]
func DeleteFileHandler(c *gin.Context) {
	// Implementation goes here
}
