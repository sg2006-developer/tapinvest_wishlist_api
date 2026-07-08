package models

import "time"

type Bond struct {
	Isin            string    `json:"isin"`
	BondName        string    `json:"bondName"`
	Yield           float64   `json:"bondYield"`
	PayoutFrequency string    `json:"payoutFrequency"`
	MaturityDate    time.Time `json:"maturityDate"`
	MinInvestment   float64   `json:"minInvestment"`
	Rating          string    `json:"rating"`
	LogoUrl         *string   `json:"logoUrl"`
	DetailUrl       *string   `json:"detailUrl"`
	Tenure          float64   `json:"tenure"`
	Color           *string   `json:"color,omitempty"`
	Position        int       `json:"position,omitempty"`
	IsPinned        bool      `json:"isPinned,omitempty"`
}
