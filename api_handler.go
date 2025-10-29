package odtimagereplacer

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleReplaceImages is the Gin handler for replacing images in ODT
func HandleReplaceImages(c *gin.Context) {
	var req ReplaceRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ReplaceResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid JSON: %v", err),
		})
		return
	}

	// Process the request
	response, outputData, err := ProcessReplaceRequest(req)
	if err != nil {
		// Determine status code based on error
		statusCode := http.StatusInternalServerError
		if response != nil && !response.Success {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, response)
		return
	}

	// Encode output to base64
	response.OutputBase64 = base64.StdEncoding.EncodeToString(outputData)

	c.JSON(http.StatusOK, response)
}

// HandleReplaceImagesDownload handles image replacement and returns the ODT file directly
func HandleReplaceImagesDownload(c *gin.Context) {
	var req ReplaceRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ReplaceResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid JSON: %v", err),
		})
		return
	}

	// Process the request
	response, outputData, err := ProcessReplaceRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Return the ODT file directly
	c.Header("Content-Type", "application/vnd.oasis.opendocument.text")
	c.Header("Content-Disposition", "attachment; filename=output.odt")
	c.Data(http.StatusOK, "application/vnd.oasis.opendocument.text", outputData)
}

// HandleHealthCheck is a simple health check endpoint
func HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "odt-image-replacer",
	})
}

// HandleInfo returns information about the service
func HandleInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":     "ODT Image Replacer API",
		"version":     "2.0.0",
		"description": "Replace images in ODT documents via JSON API",
		"endpoints": map[string]string{
			"POST /api/replace":          "Replace images and return JSON with base64 output",
			"POST /api/replace/download": "Replace images and download ODT file directly",
			"GET  /health":               "Health check endpoint",
			"GET  /info":                 "Service information",
		},
	})
}

// SetupRouter creates and configures the Gin router
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Health and info endpoints
	router.GET("/health", HandleHealthCheck)
	router.GET("/info", HandleInfo)

	// API endpoints
	api := router.Group("/api")
	{
		api.POST("/replace", HandleReplaceImages)
		api.POST("/replace/download", HandleReplaceImagesDownload)
	}

	return router
}
