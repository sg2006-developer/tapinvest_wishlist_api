package utils

import "github.com/gin-gonic/gin"

// SuccessResponse formats a standard JSON response for success
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse formats a standard JSON response for errors
func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"error":   err,
	})
}
