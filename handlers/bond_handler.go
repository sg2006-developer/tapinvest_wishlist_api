package handlers

import (
	"net/http"
	"tapinvest_api/repository"
	"tapinvest_api/utils"

	"github.com/gin-gonic/gin"
)

type BondHandler struct {
	repo repository.BondRepository
}

func NewBondHandler(r repository.BondRepository) *BondHandler {
	return &BondHandler{repo: r}
}

func (h *BondHandler) GetBonds(c *gin.Context) {
	sortBy := c.DefaultQuery("sortBy", "isin")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	bonds, err := h.repo.GetAll(c.Request.Context(), sortBy, sortOrder)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve bonds", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Bonds retrieved successfully", bonds)
}

func (h *BondHandler) SearchBonds(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		// return empty list if query is empty according to spec
		utils.SuccessResponse(c, http.StatusOK, "Bonds retrieved successfully", []interface{}{})
		return
	}
	
	bonds, err := h.repo.Search(c.Request.Context(), q)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve bonds", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Bonds retrieved successfully", bonds)
}
