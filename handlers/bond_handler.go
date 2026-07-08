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
	bonds, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve bonds", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Bonds retrieved successfully", bonds)
}
