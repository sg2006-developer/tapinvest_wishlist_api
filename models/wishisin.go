package models

import "time"

type WishIsin struct {
	WishListID string    `json:"wishListId"`
	Isin       string    `json:"isin"`
	Color      *string   `json:"color"`
	Position   int       `json:"position"`
	IsPinned   bool      `json:"isPinned"`
	CreatedAt  time.Time `json:"createdAt"`
}
