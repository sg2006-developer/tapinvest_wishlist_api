package routes

import (
	"tapinvest_api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, bondHandler *handlers.BondHandler, wishlistHandler *handlers.WishlistHandler) {
	api := router.Group("/api/v1")
	
	// Health check (not in spec, but good to keep)
	api.GET("/health", handlers.HealthCheck)

	// Bonds
	api.GET("/bond", bondHandler.GetBonds)
	api.GET("/bond/search", bondHandler.SearchBonds)

	// Wishlists
	wishlists := api.Group("/wishlist")
	{
		wishlists.POST("", wishlistHandler.CreateWishlist)
		wishlists.GET("", wishlistHandler.GetWishlists)
		wishlists.GET("/:wishlistId", wishlistHandler.GetWishlistDetail)
		wishlists.PATCH("/:wishlistId", wishlistHandler.RenameWishlist)
		wishlists.DELETE("/:wishlistId", wishlistHandler.DeleteWishlist)

		// Wishlist Items
		wishlists.POST("/:wishlistId/bond", wishlistHandler.AddBond)
		wishlists.DELETE("/:wishlistId/bond/:bondIsin", wishlistHandler.RemoveBond)
		wishlists.PATCH("/:wishlistId/bond/:bondIsin/color", wishlistHandler.SetBondColor)
		wishlists.PATCH("/:wishlistId/bond/:bondIsin/position", wishlistHandler.SetBondPosition)
		wishlists.PATCH("/:wishlistId/bond/:bondIsin/pin", wishlistHandler.SetBondPin)
		wishlists.PATCH("/:wishlistId/reorder", wishlistHandler.ReorderBonds)
	}
}