package models

import "time"

type Bond struct {
	Isin            string    `json:"isin"`
	BondName        string    `json:"bond_name"`
	Yield           float64   `json:"yield"`
	PayoutFrequency string    `json:"payout_frequency"`
	MaturityDate    time.Time `json:"maturity_date"`
	MinInvestment   float64   `json:"min_investment"`
	Rating          string    `json:"rating"`
	LogoUrl         *string   `json:"logo_url"`
	DetailUrl       *string   `json:"detail_url"`
	Tenure          float64   `json:"tenure"`
}
