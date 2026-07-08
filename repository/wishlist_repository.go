package repository

import (
	"context"
	"errors"
	"fmt"
	"tapinvest_api/models"

	"github.com/google/uuid"
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
	GetByID(ctx context.Context, id string, sortBy string) (*models.WishlistDetailResponse, error)
	Update(ctx context.Context, id string, name string) (*models.Wishlist, error)
	Delete(ctx context.Context, id string) error
	AddBond(ctx context.Context, wishlistID string, isin string) error
	RemoveBond(ctx context.Context, wishlistID string, isin string) error
	SetBondColor(ctx context.Context, wishlistID string, isin string, color *string) error
	SetBondPosition(ctx context.Context, wishlistID string, isin string, position int) error
	SetBondPin(ctx context.Context, wishlistID string, isin string, isPinned bool) error
	ReorderBonds(ctx context.Context, wishlistID string, bondIsins []string) error
}

type wishlistRepository struct {
	db *pgxpool.Pool
}

func NewWishlistRepository(db *pgxpool.Pool) WishlistRepository {
	return &wishlistRepository{db: db}
}

func (r *wishlistRepository) Create(ctx context.Context, name string) (*models.Wishlist, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_lists").Scan(&count)
	if err != nil {
		return nil, err
	}
	if count >= 5 {
		return nil, ErrMaxWishlistsReached
	}

	var existing int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_lists WHERE wish_list_name = $1", name).Scan(&existing)
	if err != nil {
		return nil, err
	}
	if existing > 0 {
		return nil, ErrWishlistNameExists
	}

	newId := uuid.New().String()

	var w models.Wishlist
	err = r.db.QueryRow(ctx, 
		"INSERT INTO wish_lists (wish_list_id, wish_list_name) VALUES ($1, $2) RETURNING wish_list_id, wish_list_name, created_at, updated_at", 
		newId, name,
	).Scan(&w.WishListID, &w.WishListName, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *wishlistRepository) GetAll(ctx context.Context) ([]models.WishlistResponse, error) {
	query := `
		SELECT w.wish_list_id, w.wish_list_name, w.created_at, w.updated_at, COUNT(wi.isin) as bond_count
		FROM wish_lists w
		LEFT JOIN wish_isin wi ON w.wish_list_id = wi.wish_list_id
		GROUP BY w.wish_list_id, w.wish_list_name, w.created_at, w.updated_at
		ORDER BY w.created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []models.WishlistResponse
	for rows.Next() {
		var w models.WishlistResponse
		if err := rows.Scan(&w.WishListID, &w.WishListName, &w.CreatedAt, &w.UpdatedAt, &w.BondCount); err != nil {
			return nil, err
		}
		wishlists = append(wishlists, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	if wishlists == nil {
		wishlists = make([]models.WishlistResponse, 0)
	}
	return wishlists, nil
}

func (r *wishlistRepository) GetByID(ctx context.Context, id string, sortBy string) (*models.WishlistDetailResponse, error) {
	var w models.WishlistDetailResponse
	err := r.db.QueryRow(ctx, "SELECT wish_list_id, wish_list_name, created_at, updated_at FROM wish_lists WHERE wish_list_id = $1", id).Scan(&w.WishListID, &w.WishListName, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}

	var count int
	r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_isin WHERE wish_list_id = $1", id).Scan(&count)
	w.BondCount = count

	orderClause := "wi.is_pinned DESC, wi.position ASC" // default manual

	switch sortBy {
	case "addedRecently":
		orderClause = "wi.is_pinned DESC, wi.created_at DESC, wi.position ASC"
	case "color":
		orderClause = "wi.is_pinned DESC, wi.color ASC NULLS LAST, wi.position ASC"
	case "yield":
		orderClause = "wi.is_pinned DESC, md.yield DESC NULLS LAST, wi.position ASC"
	case "minInvestment":
		orderClause = "wi.is_pinned DESC, md.min_investment ASC NULLS LAST, wi.position ASC"
	case "tenure":
		orderClause = "wi.is_pinned DESC, md.tenure ASC, wi.position ASC"
	case "rating":
		orderClause = "wi.is_pinned DESC, md.rating ASC NULLS LAST, wi.position ASC"
	}

	query := fmt.Sprintf(`
		SELECT md.isin, md.bond_name, md.yield, md.payout_frequency, md.maturity_date, md.min_investment, md.rating, md.logo_url, md.detail_url, md.tenure, wi.color, wi.position, wi.is_pinned
		FROM master_data md
		JOIN wish_isin wi ON md.isin = wi.isin
		WHERE wi.wish_list_id = $1
		ORDER BY %s
	`, orderClause)
	
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	w.Bonds = make([]models.Bond, 0)
	for rows.Next() {
		var b models.Bond
		if err := rows.Scan(&b.Isin, &b.BondName, &b.Yield, &b.PayoutFrequency, &b.MaturityDate, &b.MinInvestment, &b.Rating, &b.LogoUrl, &b.DetailUrl, &b.Tenure, &b.Color, &b.Position, &b.IsPinned); err != nil {
			return nil, err
		}
		w.Bonds = append(w.Bonds, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *wishlistRepository) Update(ctx context.Context, id string, name string) (*models.Wishlist, error) {
	var existingId string
	err := r.db.QueryRow(ctx, "SELECT wish_list_id FROM wish_lists WHERE wish_list_id = $1", id).Scan(&existingId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}

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
		"UPDATE wish_lists SET wish_list_name = $1, updated_at = CURRENT_TIMESTAMP WHERE wish_list_id = $2 RETURNING wish_list_id, wish_list_name, created_at, updated_at", 
		name, id,
	).Scan(&w.WishListID, &w.WishListName, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *wishlistRepository) Delete(ctx context.Context, id string) error {
	cmd, err := r.db.Exec(ctx, "DELETE FROM wish_lists WHERE wish_list_id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

func (r *wishlistRepository) AddBond(ctx context.Context, wishlistID string, isin string) error {
	var wId string
	err := r.db.QueryRow(ctx, "SELECT wish_list_id FROM wish_lists WHERE wish_list_id = $1", wishlistID).Scan(&wId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrWishlistNotFound
		}
		return err
	}

	var bIsin string
	err = r.db.QueryRow(ctx, "SELECT isin FROM master_data WHERE isin = $1", isin).Scan(&bIsin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrBondNotFound
		}
		return err
	}

	var count int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_isin WHERE wish_list_id = $1", wishlistID).Scan(&count)
	if err != nil {
		return err
	}
	if count >= 10 {
		return ErrMaxBondsReached
	}

	var existing int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM wish_isin WHERE wish_list_id = $1 AND isin = $2", wishlistID, isin).Scan(&existing)
	if err != nil {
		return err
	}
	if existing > 0 {
		return ErrBondDuplicate
	}

	// Calculate max position to append to end
	var maxPos int
	err = r.db.QueryRow(ctx, "SELECT COALESCE(MAX(position), -1) FROM wish_isin WHERE wish_list_id = $1", wishlistID).Scan(&maxPos)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, "INSERT INTO wish_isin (wish_list_id, isin, position) VALUES ($1, $2, $3)", wishlistID, isin, maxPos+1)
	return err
}

func (r *wishlistRepository) RemoveBond(ctx context.Context, wishlistID string, isin string) error {
	cmd, err := r.db.Exec(ctx, "DELETE FROM wish_isin WHERE wish_list_id = $1 AND isin = $2", wishlistID, isin)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

func (r *wishlistRepository) SetBondColor(ctx context.Context, wishlistID string, isin string, color *string) error {
	cmd, err := r.db.Exec(ctx, "UPDATE wish_isin SET color = $1 WHERE wish_list_id = $2 AND isin = $3", color, wishlistID, isin)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

func (r *wishlistRepository) SetBondPosition(ctx context.Context, wishlistID string, isin string, position int) error {
	cmd, err := r.db.Exec(ctx, "UPDATE wish_isin SET position = $1 WHERE wish_list_id = $2 AND isin = $3", position, wishlistID, isin)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

func (r *wishlistRepository) SetBondPin(ctx context.Context, wishlistID string, isin string, isPinned bool) error {
	cmd, err := r.db.Exec(ctx, "UPDATE wish_isin SET is_pinned = $1 WHERE wish_list_id = $2 AND isin = $3", isPinned, wishlistID, isin)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

func (r *wishlistRepository) ReorderBonds(ctx context.Context, wishlistID string, bondIsins []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Verify all bonds exist in the wishlist and count matches
	var count int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM wish_isin WHERE wish_list_id = $1", wishlistID).Scan(&count)
	if err != nil {
		return err
	}

	if count != len(bondIsins) {
		return errors.New("bondIsins must contain all N bonds in the wishlist")
	}

	for i, isin := range bondIsins {
		cmd, err := tx.Exec(ctx, "UPDATE wish_isin SET position = $1 WHERE wish_list_id = $2 AND isin = $3", i, wishlistID, isin)
		if err != nil {
			return err
		}
		if cmd.RowsAffected() == 0 {
			return errors.New("one or more bondIsins do not exist in this wishlist")
		}
	}

	return tx.Commit(ctx)
}
