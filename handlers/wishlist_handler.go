package handlers

import (
	"net/http"
	"tapinvest_api/repository"
	"tapinvest_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
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
	id := c.Param("wishlistId")
	if !isValidUUID(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID.", nil)
		return
	}

	sortBy := c.DefaultQuery("sortBy", "manual")

	wishlist, err := h.repo.GetByID(c.Request.Context(), id, sortBy)
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
	id := c.Param("wishlistId")
	if !isValidUUID(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID.", nil)
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
	id := c.Param("wishlistId")
	if !isValidUUID(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID.", nil)
		return
	}

	err := h.repo.Delete(c.Request.Context(), id)
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
	BondIsin string `json:"bondIsin" binding:"required"`
}

func (h *WishlistHandler) AddBond(c *gin.Context) {
	id := c.Param("wishlistId")
	if !isValidUUID(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID.", nil)
		return
	}

	var req AddBondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err := h.repo.AddBond(c.Request.Context(), id, req.BondIsin)
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
	id := c.Param("wishlistId")
	if !isValidUUID(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID.", nil)
		return
	}

	isin := c.Param("bondIsin")
	if isin == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid bond ISIN", nil)
		return
	}

	err := h.repo.RemoveBond(c.Request.Context(), id, isin)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrWishlistNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, "Failed to remove bond from wishlist", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Bond removed from wishlist successfully", nil)
}

type ColorRequest struct {
	Color *string `json:"color"`
}

func (h *WishlistHandler) SetBondColor(c *gin.Context) {
	id := c.Param("wishlistId")
	if !isValidUUID(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID.", nil)
		return
	}

	isin := c.Param("bondIsin")
	var req ColorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body.", err.Error())
		return
	}

	err := h.repo.SetBondColor(c.Request.Context(), id, isin, req.Color)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrWishlistNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, "Failed to set bond color", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Bond color updated successfully", nil)
}

type PositionRequest struct {
	Position *int `json:"position" binding:"required"`
}

func (h *WishlistHandler) SetBondPosition(c *gin.Context) {
	id := c.Param("wishlistId")
	if !isValidUUID(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID.", nil)
		return
	}

	isin := c.Param("bondIsin")
	var req PositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "position must be a non-negative integer.", err.Error())
		return
	}

	if *req.Position < 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "position must be a non-negative integer.", nil)
		return
	}

	err := h.repo.SetBondPosition(c.Request.Context(), id, isin, *req.Position)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrWishlistNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, "Failed to set bond position", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Bond position updated successfully", nil)
}

type PinRequest struct {
	IsPinned *bool `json:"isPinned"`
}

func (h *WishlistHandler) SetBondPin(c *gin.Context) {
	id := c.Param("wishlistId")
	if !isValidUUID(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID.", nil)
		return
	}

	isin := c.Param("bondIsin")
	var req PinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "isPinned must be a boolean.", err.Error())
		return
	}

	pinned := false
	if req.IsPinned != nil {
		pinned = *req.IsPinned
	}

	err := h.repo.SetBondPin(c.Request.Context(), id, isin, pinned)
	if err != nil {
		status := http.StatusInternalServerError
		if err == repository.ErrWishlistNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, "Failed to pin/unpin bond", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Bond pin status updated successfully", nil)
}

type ReorderRequest struct {
	BondIsins []string `json:"bondIsins" binding:"required"`
}

func (h *WishlistHandler) ReorderBonds(c *gin.Context) {
	id := c.Param("wishlistId")
	if !isValidUUID(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid UUID.", nil)
		return
	}

	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.BondIsins) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "bondIsins is required and must be a non-empty array.", err.Error())
		return
	}

	err := h.repo.ReorderBonds(c.Request.Context(), id, req.BondIsins)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bondIsins must contain all N bonds in the wishlist" || err.Error() == "one or more bondIsins do not exist in this wishlist" {
			status = http.StatusBadRequest
		} else if err == repository.ErrWishlistNotFound {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, "Failed to reorder bonds", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Bonds reordered successfully", nil)
}
