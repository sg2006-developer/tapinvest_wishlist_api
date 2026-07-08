package routes

import (
	"tapinvest_api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, bondHandler *handlers.BondHandler, wishlistHandler *handlers.WishlistHandler) {
	api := router.Group("/api")
	
	// Health check
	api.GET("/health", handlers.HealthCheck)

	// Bonds
	api.GET("/bonds", bondHandler.GetBonds)

	// Wishlists
	wishlists := api.Group("/wishlists")
	{
		wishlists.POST("", wishlistHandler.CreateWishlist)
		wishlists.GET("", wishlistHandler.GetWishlists)
		wishlists.GET("/:wishlistId", wishlistHandler.GetWishlistDetail)
		wishlists.PUT("/:wishlistId", wishlistHandler.RenameWishlist)
		wishlists.DELETE("/:wishlistId", wishlistHandler.DeleteWishlist)

		// Wishlist Items
		items := wishlists.Group("/:wishlistId/items")
		{
			items.POST("", wishlistHandler.AddBond)
			items.DELETE("/:bondIsin", wishlistHandler.RemoveBond)
		}
	}
}