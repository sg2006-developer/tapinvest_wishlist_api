package repository

import (
	"context"
	"errors"
	"fmt"
	"tapinvest_api/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMaxWishlistsReached = errors.New("maximum of 5 wishlists allowed")
	ErrWishlistNotFound    = errors.New("wishlist not found")
	ErrWishlistNameExists  = errors.New("wishlist name already exists")
	ErrMaxBondsReached     = errors.New("maximum of 10 bonds allowed per wishlist")
	ErrBondDuplicate       = errors.New("bond already exists in the wishlist")
	ErrBondNotFound        = errors.New("bond does not exist in master data")
)

type WishlistRepository interface {
	Create(ctx context.Context, name string) (*models.Wishlist, error)
	GetAll(ctx context.Context) ([]models.WishlistResponse, error)
	GetByID(ctx context.Context, id int) (*models.WishlistDetailResponse, error)
	Update(ctx context.Context, id int, name string) (*models.Wishlist, error)
	Delete(ctx context.Context, id int) error
	AddBond(ctx context.Context, wishlistID int, isin string) error
	RemoveBond(ctx context.Context, wishlistID int, isin string) error
}

type wishlistRepository struct {
	db *pgxpool.Pool
}

func NewWishlistRepository(db *pgxpool.Pool) WishlistRepository {
	return &wishlistRepository{db: db}
}

func (r *wishlistRepository) Create(ctx context.Context, name string) (*models.Wishlist, error) {
	// Check max wishlists
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_lists").Scan(&count)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, ErrMaxWishlistsReached
	}

	// Check if name already exists
	var existing int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_lists WHERE wish_list_name = $1", name).Scan(&existing)
	if err != nil {
		return nil, err
	}
	if existing > 0 {
		return nil, ErrWishlistNameExists
	}

	// Insert
	var w models.Wishlist
	err = r.db.QueryRow(ctx, 
		"INSERT INTO wish_lists (wish_list_name) VALUES ($1) RETURNING wish_list_id, wish_list_name", 
		name,
	).Scan(&w.WishListID, &w.WishListName)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *wishlistRepository) GetAll(ctx context.Context) ([]models.WishlistResponse, error) {
	query := `
		SELECT w.wish_list_id, w.wish_list_name, COUNT(wi.isin) as bond_count
		FROM wish_lists w
		LEFT JOIN wish_isin wi ON w.wish_list_id = wi.wish_list_id
		GROUP BY w.wish_list_id, w.wish_list_name
		ORDER BY w.wish_list_id ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []models.WishlistResponse
	for rows.Next() {
		var w models.WishlistResponse
		if err := rows.Scan(&w.WishListID, &w.WishListName, &w.BondCount); err != nil {
			return nil, err
		}
		wishlists = append(wishlists, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice instead of nil if no wishlists exist
	if wishlists == nil {
		wishlists = make([]models.WishlistResponse, 0)
	}
	return wishlists, nil
}

func (r *wishlistRepository) GetByID(ctx context.Context, id int) (*models.WishlistDetailResponse, error) {
	// First check if wishlist exists and get its details
	var w models.WishlistDetailResponse
	err := r.db.QueryRow(ctx, "SELECT wish_list_id, wish_list_name FROM wish_lists WHERE wish_list_id = $1", id).Scan(&w.WishListID, &w.WishListName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}

	// Get bonds
	query := `
		SELECT md.isin, md.bond_name, md.yield, md.payout_frequency, md.maturity_date, md.min_investment, md.rating, md.logo_url, md.detail_url, md.tenure
		FROM master_data md
		JOIN wish_isin wi ON md.isin = wi.isin
		WHERE wi.wish_list_id = $1
	`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	w.Bonds = make([]models.Bond, 0)
	for rows.Next() {
		var b models.Bond
		if err := rows.Scan(&b.Isin, &b.BondName, &b.Yield, &b.PayoutFrequency, &b.MaturityDate, &b.MinInvestment, &b.Rating, &b.LogoUrl, &b.DetailUrl, &b.Tenure); err != nil {
			return nil, err
		}
		w.Bonds = append(w.Bonds, b)
	}

	return &w, nil
}

func (r *wishlistRepository) Update(ctx context.Context, id int, name string) (*models.Wishlist, error) {
	// Check if wishlist exists
	var existingId int
	err := r.db.QueryRow(ctx, "SELECT wish_list_id FROM wish_lists WHERE wish_list_id = $1", id).Scan(&existingId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}

	// Check for name collision
	var count int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_lists WHERE wish_list_name = $1 AND wish_list_id != $2", name, id).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrWishlistNameExists
	}

	var w models.Wishlist
	err = r.db.QueryRow(ctx, 
		"UPDATE wish_lists SET wish_list_name = $1 WHERE wish_list_id = $2 RETURNING wish_list_id, wish_list_name", 
		name, id,
	).Scan(&w.WishListID, &w.WishListName)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *wishlistRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, "DELETE FROM wish_lists WHERE wish_list_id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

func (r *wishlistRepository) AddBond(ctx context.Context, wishlistID int, isin string) error {
	// 1. Wishlist exists?
	var wId int
	err := r.db.QueryRow(ctx, "SELECT wish_list_id FROM wish_lists WHERE wish_list_id = $1", wishlistID).Scan(&wId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrWishlistNotFound
		}
		return err
	}

	// 2. Bond exists in master data?
	var bIsin string
	err = r.db.QueryRow(ctx, "SELECT isin FROM master_data WHERE isin = $1", isin).Scan(&bIsin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrBondNotFound
		}
		return err
	}

	// 3. Wishlist has 10 bonds?
	var count int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_isin WHERE wish_list_id = $1", wishlistID).Scan(&count)
	if err != nil {
		return err
	}
	if count >= 10 {
		return ErrMaxBondsReached
	}

	// 4. Bond already exists? (Duplicate check)
	var existing int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_isin WHERE wish_list_id = $1 AND isin = $2", wishlistID, isin).Scan(&existing)
	if err != nil {
		return err
	}
	if existing > 0 {
		return ErrBondDuplicate
	}

	// 5. Insert
	_, err = r.db.Exec(ctx, "INSERT INTO wish_isin (wish_list_id, isin) VALUES ($1, $2)", wishlistID, isin)
	return err
}

func (r *wishlistRepository) RemoveBond(ctx context.Context, wishlistID int, isin string) error {
	cmd, err := r.db.Exec(ctx, "DELETE FROM wish_isin WHERE wish_list_id = $1 AND isin = $2", wishlistID, isin)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("bond not found in wishlist")
	}
	return nil
}
