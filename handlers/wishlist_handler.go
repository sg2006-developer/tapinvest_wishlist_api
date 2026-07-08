package handlers

import (
	"net/http"
	"strconv"
	"tapinvest_api/repository"
	"tapinvest_api/utils"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	repo repository.WishlistRepository
}

func NewWishlistHandler(r repository.WishlistRepository) *WishlistHandler {
	return &WishlistHandler{repo: r}
}

type CreateWishlistRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *WishlistHandler) CreateWishlist(c *gin.Context) {
	var req CreateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	validName, valid := utils.ValidateWishlistName(req.Name)
	if !valid {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid wishlist name", "Name must be non-empty and max 25 characters")
		return
	}

	wishlist, err := h.repo.Create(c.Request.Context(), validName)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrMaxWishlistsReached || err == repository.ErrWishlistNameExists {
			status = http.StatusBadRequest
		}
		utils.ErrorResponse(c, status, "Failed to create wishlist", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Wishlist created successfully", wishlist)
}

func (h *WishlistHandler) GetWishlists(c *gin.Context) {
	wishlists, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve wishlists", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Wishlists retrieved successfully", wishlists)
}

func (h *WishlistHandler) GetWishlistDetail(c *gin.Context) {
	idStr := c.Param("wishlistId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid wishlist ID", nil)
		return
	}

	wishlist, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrWishlistNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, "Failed to retrieve wishlist", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Wishlist retrieved successfully", wishlist)
}

type UpdateWishlistRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *WishlistHandler) RenameWishlist(c *gin.Context) {
	idStr := c.Param("wishlistId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid wishlist ID", nil)
		return
	}

	var req UpdateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	validName, valid := utils.ValidateWishlistName(req.Name)
	if !valid {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid wishlist name", "Name must be non-empty and max 25 characters")
		return
	}

	wishlist, err := h.repo.Update(c.Request.Context(), id, validName)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrWishlistNotFound {
			status = http.StatusNotFound
		} else if err == repository.ErrWishlistNameExists {
			status = http.StatusBadRequest
		}
		utils.ErrorResponse(c, status, "Failed to rename wishlist", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Wishlist renamed successfully", wishlist)
}

func (h *WishlistHandler) DeleteWishlist(c *gin.Context) {
	idStr := c.Param("wishlistId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid wishlist ID", nil)
		return
	}

	err = h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrWishlistNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, "Failed to delete wishlist", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Wishlist deleted successfully", nil)
}

type AddBondRequest struct {
	Isin string `json:"isin" binding:"required"`
}

func (h *WishlistHandler) AddBond(c *gin.Context) {
	idStr := c.Param("wishlistId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid wishlist ID", nil)
		return
	}

	var req AddBondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err = h.repo.AddBond(c.Request.Context(), id, req.Isin)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case repository.ErrWishlistNotFound, repository.ErrBondNotFound:
			status = http.StatusNotFound
		case repository.ErrMaxBondsReached, repository.ErrBondDuplicate:
			status = http.StatusBadRequest
		}
		utils.ErrorResponse(c, status, "Failed to add bond to wishlist", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Bond added to wishlist successfully", nil)
}

func (h *WishlistHandler) RemoveBond(c *gin.Context) {
	idStr := c.Param("wishlistId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid wishlist ID", nil)
		return
	}

	isin := c.Param("bondIsin")
	if isin == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid bond ISIN", nil)
		return
	}

	err = h.repo.RemoveBond(c.Request.Context(), id, isin)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove bond from wishlist", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Bond removed from wishlist successfully", nil)
}
