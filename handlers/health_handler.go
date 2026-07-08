package handlers

import (
	"net/http"
	"tapinvest_api/utils"

	"github.com/gin-gonic/gin"
)

// HealthCheck verifies if the API is running
func HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "API is running properly", gin.H{"status": "ok"})
}