package repository

import (
	"context"
	"fmt"
	"strings"
	"tapinvest_api/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BondRepository interface {
	GetAll(ctx context.Context, sortBy string, sortOrder string) ([]models.Bond, error)
	Search(ctx context.Context, query string) ([]models.Bond, error)
}

type bondRepository struct {
	db *pgxpool.Pool
}

func NewBondRepository(db *pgxpool.Pool) BondRepository {
	return &bondRepository{db: db}
}

func (r *bondRepository) GetAll(ctx context.Context, sortBy string, sortOrder string) ([]models.Bond, error) {
	orderClause := "isin ASC"

	direction := "ASC"
	if strings.ToLower(sortOrder) == "desc" {
		direction = "DESC"
	}

	nullsLast := "NULLS LAST"

	switch sortBy {
	case "bondYield":
		orderClause = fmt.Sprintf("yield %s %s, isin %s", direction, nullsLast, direction)
	case "minInvestment":
		orderClause = fmt.Sprintf("min_investment %s %s, isin %s", direction, nullsLast, direction)
	case "tenure":
		orderClause = fmt.Sprintf("tenure %s, isin %s", direction, direction)
	case "rating":
		// Example custom rating sort (simple alphabetical if acceptable, or requires custom case)
		// For now using standard sort which works well if ratings are A, AA, AAA
		orderClause = fmt.Sprintf("rating %s %s, isin %s", direction, nullsLast, direction)
	case "isin":
		orderClause = fmt.Sprintf("isin %s", direction)
	}

	queryStr := fmt.Sprintf(`SELECT isin, bond_name, yield, payout_frequency, maturity_date, min_investment, rating, logo_url, detail_url, tenure FROM master_data ORDER BY %s`, orderClause)
	rows, err := r.db.Query(ctx, queryStr)
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

	if bonds == nil {
		bonds = make([]models.Bond, 0)
	}
	return bonds, nil
}

func (r *bondRepository) Search(ctx context.Context, q string) ([]models.Bond, error) {
	queryStr := `SELECT isin, bond_name, yield, payout_frequency, maturity_date, min_investment, rating, logo_url, detail_url, tenure 
				 FROM master_data 
				 WHERE isin ILIKE $1 OR bond_name ILIKE $1 OR similarity(bond_name, $2) > 0.2
				 ORDER BY similarity(bond_name, $2) DESC, isin ASC`
	rows, err := r.db.Query(ctx, queryStr, "%"+q+"%", q)
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
	if bonds == nil {
		bonds = make([]models.Bond, 0)
	}
	return bonds, nil
}
