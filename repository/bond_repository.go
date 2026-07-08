package repository

import (
	"context"
	"tapinvest_api/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BondRepository interface {
	GetAll(ctx context.Context) ([]models.Bond, error)
}

type bondRepository struct {
	db *pgxpool.Pool
}

func NewBondRepository(db *pgxpool.Pool) BondRepository {
	return &bondRepository{db: db}
}

func (r *bondRepository) GetAll(ctx context.Context) ([]models.Bond, error) {
	query := `SELECT isin, bond_name, yield, payout_frequency, maturity_date, min_investment, rating, logo_url, detail_url, tenure FROM master_data ORDER BY isin ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bonds []models.Bond
	for rows.Next() {
		var b models.Bond
		if err := rows.Scan(&b.Isin, &b.BondName, &b.Yield, &b.PayoutFrequency, &b.MaturityDate, &b.MinInvestment, &b.Rating, &b.LogoUrl, &b.DetailUrl, &b.Tenure); err != nil {
			return nil, err
		}
		bonds = append(bonds, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bonds, nil
}
