package models

type WishIsin struct {
	WishListID int    `json:"wish_list_id"`
	Isin       string `json:"isin"`
}
