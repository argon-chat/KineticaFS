package router

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// Helper function to generate a valid JWT token for testing
func generateTestJWT(sub, spaceId, fileType, secret string) (string, error) {
	claims := JWTClaims{
		Sub:  sub,
		Type: fileType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	if spaceId != "" {
		claims.SpaceId = &spaceId
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	backends := []string{"postgres"}
	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			app := setupTestApp(t, backend)
			adminToken := getOrCreateAdminToken(t, app)

			// Set JWT secret
			jwtSecret := "test-secret-key-for-jwt-validation"
			viper.Set("jwt_secret", jwtSecret)

			// Generate valid JWT token
			testSub := "550e8400-e29b-41d4-a716-446655440000"
			testSpaceId := "660e8400-e29b-41d4-a716-446655440001"
			testType := "avatar"

			jwtToken, err := generateTestJWT(testSub, testSpaceId, testType, jwtSecret)
			require.NoError(t, err)

			// Setup regions and bucket for file upload
			regionsPath, bucket1, _, _, cleanup := setupRegionsConfig(t, app, adminToken)
			defer cleanup()
			viper.Set("region", regionsPath)

			// Initiate file upload
			initiateBody := map[string]interface{}{
				"regionId":   "ru-1",
				"bucketCode": bucket1.ID,
			}
			initiateBodyBytes, err := json.Marshal(initiateBody)
			require.NoError(t, err)

			resp := makeAuthRequest(t, app, "POST", "/api/v1/file/", adminToken, bytes.NewReader(initiateBodyBytes))
			require.Equal(t, http.StatusCreated, resp.StatusCode)

			var initiateResp InitiateFileUploadResponse
			err = json.NewDecoder(resp.Body).Decode(&initiateResp)
			require.NoError(t, err)
			resp.Body.Close()
			blobID := initiateResp.URL

			// Upload file with JWT token
			fileContent := []byte("test file content for JWT upload")
			req, err := http.NewRequest("PATCH", app.server.URL+"/api/v1/upload/"+blobID, bytes.NewReader(fileContent))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+jwtToken)
			req.Header.Set("Content-Type", "application/octet-stream")

			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent {
				body, _ := io.ReadAll(resp.Body)
				t.Logf("Upload failed with status %d: %s", resp.StatusCode, string(body))
			}
			require.Equal(t, http.StatusNoContent, resp.StatusCode)

			// Verify JWT claims were stored in file record
			resp = makeAuthRequest(t, app, "POST", "/api/v1/file/"+blobID+"/finalize", adminToken, nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			var file map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&file)
			require.NoError(t, err)
			resp.Body.Close()

			// Check that JWT claims were stored
			require.Equal(t, testSub, file["user_sub"])
			require.Equal(t, testSpaceId, file["space_id"])
			require.Equal(t, testType, file["file_type"])
		})
	}
}

func TestJWTAuthMiddleware_MissingToken(t *testing.T) {
	backends := []string{"postgres"}
	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			app := setupTestApp(t, backend)
			adminToken := getOrCreateAdminToken(t, app)

			// Setup regions and bucket
			regionsPath, bucket1, _, _, cleanup := setupRegionsConfig(t, app, adminToken)
			defer cleanup()
			viper.Set("region", regionsPath)

			// Initiate file upload
			initiateBody := map[string]interface{}{
				"regionId":   "ru-1",
				"bucketCode": bucket1.ID,
			}
			initiateBodyBytes, err := json.Marshal(initiateBody)
			require.NoError(t, err)

			resp := makeAuthRequest(t, app, "POST", "/api/v1/file/", adminToken, bytes.NewReader(initiateBodyBytes))
			require.Equal(t, http.StatusCreated, resp.StatusCode)

			var initiateResp InitiateFileUploadResponse
			err = json.NewDecoder(resp.Body).Decode(&initiateResp)
			require.NoError(t, err)
			resp.Body.Close()
			blobID := initiateResp.URL

			// Try to upload without authentication
			fileContent := []byte("test file content")
			req, err := http.NewRequest("PATCH", app.server.URL+"/api/v1/upload/"+blobID, bytes.NewReader(fileContent))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/octet-stream")

			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	backends := []string{"postgres"}
	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			app := setupTestApp(t, backend)
			adminToken := getOrCreateAdminToken(t, app)

			// Set JWT secret
			jwtSecret := "test-secret-key-for-jwt-validation"
			viper.Set("jwt_secret", jwtSecret)

			// Setup regions and bucket
			regionsPath, bucket1, _, _, cleanup := setupRegionsConfig(t, app, adminToken)
			defer cleanup()
			viper.Set("region", regionsPath)

			// Initiate file upload
			initiateBody := map[string]interface{}{
				"regionId":   "ru-1",
				"bucketCode": bucket1.ID,
			}
			initiateBodyBytes, err := json.Marshal(initiateBody)
			require.NoError(t, err)

			resp := makeAuthRequest(t, app, "POST", "/api/v1/file/", adminToken, bytes.NewReader(initiateBodyBytes))
			require.Equal(t, http.StatusCreated, resp.StatusCode)

			var initiateResp InitiateFileUploadResponse
			err = json.NewDecoder(resp.Body).Decode(&initiateResp)
			require.NoError(t, err)
			resp.Body.Close()
			blobID := initiateResp.URL

			// Try to upload with invalid JWT token
			fileContent := []byte("test file content")
			req, err := http.NewRequest("PATCH", app.server.URL+"/api/v1/upload/"+blobID, bytes.NewReader(fileContent))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer invalid-token-here")
			req.Header.Set("Content-Type", "application/octet-stream")

			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func TestJWTAuthMiddleware_MissingRequiredClaims(t *testing.T) {
	backends := []string{"postgres"}
	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			app := setupTestApp(t, backend)
			adminToken := getOrCreateAdminToken(t, app)

			// Set JWT secret
			jwtSecret := "test-secret-key-for-jwt-validation"
			viper.Set("jwt_secret", jwtSecret)

			// Setup regions and bucket
			regionsPath, bucket1, _, _, cleanup := setupRegionsConfig(t, app, adminToken)
			defer cleanup()
			viper.Set("region", regionsPath)

			// Initiate file upload
			initiateBody := map[string]interface{}{
				"regionId":   "ru-1",
				"bucketCode": bucket1.ID,
			}
			initiateBodyBytes, err := json.Marshal(initiateBody)
			require.NoError(t, err)

			resp := makeAuthRequest(t, app, "POST", "/api/v1/file/", adminToken, bytes.NewReader(initiateBodyBytes))
			require.Equal(t, http.StatusCreated, resp.StatusCode)

			var initiateResp InitiateFileUploadResponse
			err = json.NewDecoder(resp.Body).Decode(&initiateResp)
			require.NoError(t, err)
			resp.Body.Close()
			blobID := initiateResp.URL

			// Generate JWT token missing required 'type' claim
			claims := jwt.MapClaims{
				"sub": "550e8400-e29b-41d4-a716-446655440000",
				"exp": time.Now().Add(1 * time.Hour).Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			jwtToken, err := token.SignedString([]byte(jwtSecret))
			require.NoError(t, err)

			// Try to upload with token missing 'type' claim
			fileContent := []byte("test file content")
			req, err := http.NewRequest("PATCH", app.server.URL+"/api/v1/upload/"+blobID, bytes.NewReader(fileContent))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+jwtToken)
			req.Header.Set("Content-Type", "application/octet-stream")

			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func TestJWTAuthMiddleware_WithoutSpaceId(t *testing.T) {
	backends := []string{"postgres"}
	for _, backend := range backends {
		t.Run(backend, func(t *testing.T) {
			app := setupTestApp(t, backend)
			adminToken := getOrCreateAdminToken(t, app)

			// Set JWT secret
			jwtSecret := "test-secret-key-for-jwt-validation"
			viper.Set("jwt_secret", jwtSecret)

			// Generate JWT token without spaceId (should be valid)
			testSub := "550e8400-e29b-41d4-a716-446655440000"
			testType := "avatar"

			jwtToken, err := generateTestJWT(testSub, "", testType, jwtSecret)
			require.NoError(t, err)

			// Setup regions and bucket
			regionsPath, bucket1, _, _, cleanup := setupRegionsConfig(t, app, adminToken)
			defer cleanup()
			viper.Set("region", regionsPath)

			// Initiate file upload
			initiateBody := map[string]interface{}{
				"regionId":   "ru-1",
				"bucketCode": bucket1.ID,
			}
			initiateBodyBytes, err := json.Marshal(initiateBody)
			require.NoError(t, err)

			resp := makeAuthRequest(t, app, "POST", "/api/v1/file/", adminToken, bytes.NewReader(initiateBodyBytes))
			require.Equal(t, http.StatusCreated, resp.StatusCode)

			var initiateResp InitiateFileUploadResponse
			err = json.NewDecoder(resp.Body).Decode(&initiateResp)
			require.NoError(t, err)
			resp.Body.Close()
			blobID := initiateResp.URL

			// Upload file with JWT token (no spaceId)
			fileContent := []byte("test file content without space")
			req, err := http.NewRequest("PATCH", app.server.URL+"/api/v1/upload/"+blobID, bytes.NewReader(fileContent))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+jwtToken)
			req.Header.Set("Content-Type", "application/octet-stream")

			resp, err = http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent {
				body, _ := io.ReadAll(resp.Body)
				t.Logf("Upload failed with status %d: %s", resp.StatusCode, string(body))
			}
			require.Equal(t, http.StatusNoContent, resp.StatusCode)

			// Verify JWT claims were stored (spaceId should be nil)
			resp = makeAuthRequest(t, app, "POST", "/api/v1/file/"+blobID+"/finalize", adminToken, nil)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			var file map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&file)
			require.NoError(t, err)
			resp.Body.Close()

			// Check that JWT claims were stored
			require.Equal(t, testSub, file["user_sub"])
			require.Nil(t, file["space_id"])
			require.Equal(t, testType, file["file_type"])
		})
	}
}
