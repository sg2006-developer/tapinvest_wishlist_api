package models

type Wishlist struct {
	WishListID   int    `json:"wish_list_id"`
	WishListName string `json:"wish_list_name"`
}

type WishlistResponse struct {
	WishListID   int    `json:"wish_list_id"`
	WishListName string `json:"wish_list_name"`
	BondCount    int    `json:"bond_count"`
}

type WishlistDetailResponse struct {
	WishListID   int    `json:"wish_list_id"`
	WishListName string `json:"wish_list_name"`
	Bonds        []Bond `json:"bonds"`
}
