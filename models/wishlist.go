package models

import "time"

type Wishlist struct {
	WishListID   string    `json:"id"`
	WishListName string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type WishlistResponse struct {
	WishListID   string    `json:"id"`
	WishListName string    `json:"name"`
	BondCount    int       `json:"bondCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type WishlistDetailResponse struct {
	WishListID   string    `json:"id"`
	WishListName string    `json:"name"`
	BondCount    int       `json:"bondCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Bonds        []Bond    `json:"bonds"`
}
