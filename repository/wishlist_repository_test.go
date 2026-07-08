package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	_ = godotenv.Load("../.env")
	
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err, "Failed to connect to database. Ensure PostgreSQL is running and .env is configured.")

	// Clean up any old test data before starting
	_, _ = pool.Exec(context.Background(), "DELETE FROM wish_lists WHERE wish_list_name LIKE 'TEST_QA_%'")

	cleanup := func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM wish_lists WHERE wish_list_name LIKE 'TEST_QA_%'")
		pool.Close()
	}

	return pool, cleanup
}

func TestWishlistRepository_LimitsAndDuplicates(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()
	repo := NewWishlistRepository(pool)

	// Ensure we are below the 5 wishlist limit before testing the max limit rule
	var count int
	_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM wish_lists").Scan(&count)
	if count > 0 {
		t.Skip("Skipping max wishlist test because existing wishlists interfere. Clear DB to run fully.")
	}

	// 1. Create 5 wishlists (Limit testing)
	var createdIds []int
	for i := 1; i <= 5; i++ {
		w, err := repo.Create(ctx, fmt.Sprintf("TEST_QA_Wishlist_%d", i))
		require.NoError(t, err)
		createdIds = append(createdIds, w.WishListID)
	}

	// 2. Try to create 6th wishlist - should fail
	_, err := repo.Create(ctx, "TEST_QA_Wishlist_6")
	assert.ErrorIs(t, err, ErrMaxWishlistsReached, "Should block 6th wishlist")

	// 3. Duplicate Group Names
	// Need to delete one to free up space
	err = repo.Delete(ctx, createdIds[4])
	require.NoError(t, err)

	_, err = repo.Create(ctx, "TEST_QA_Wishlist_1")
	assert.ErrorIs(t, err, ErrWishlistNameExists, "Should block duplicate wishlist name")

	// Create a valid one for bond testing
	targetWishlist, err := repo.Create(ctx, "TEST_QA_Bond_Testing")
	require.NoError(t, err)

	// Setup a dummy bond in master_data for testing bonds
	testIsin := "TEST_ISIN_123"
	_, _ = pool.Exec(ctx, "INSERT INTO master_data (bond_name, yield, payout_frequency, maturity_date, min_investment, rating, isin, tenure) VALUES ('Test Bond', 10.5, 'Monthly', '2030-01-01', 10000, 'AAA', $1, 12) ON CONFLICT DO NOTHING", testIsin)
	
	// 4. Duplicate Bonds Not Allowed in the Same Group
	err = repo.AddBond(ctx, targetWishlist.WishListID, testIsin)
	require.NoError(t, err, "First addition should succeed")

	err = repo.AddBond(ctx, targetWishlist.WishListID, testIsin)
	assert.ErrorIs(t, err, ErrBondDuplicate, "Duplicate bond addition should fail")

	// 5. Max 10 Bonds Allowed
	// We need 10 distinct bonds.
	for i := 1; i <= 10; i++ {
		isin := fmt.Sprintf("TEST_ISIN_%d", i)
		_, _ = pool.Exec(ctx, "INSERT INTO master_data (bond_name, yield, payout_frequency, maturity_date, min_investment, rating, isin, tenure) VALUES ('Test', 10, 'M', '2030-01-01', 10000, 'A', $1, 12) ON CONFLICT DO NOTHING", isin)
		
		if i == 1 {
			continue // Already added TEST_ISIN_123, we'll just add 9 more to reach 10
		}
		err = repo.AddBond(ctx, targetWishlist.WishListID, isin)
		require.NoError(t, err)
	}

	// Try adding 11th bond
	err = repo.AddBond(ctx, targetWishlist.WishListID, "TEST_ISIN_10") 
	assert.ErrorIs(t, err, ErrMaxBondsReached, "Should block 11th bond")

	// 6. Group Deletion Behavior (Cascading)
	err = repo.Delete(ctx, targetWishlist.WishListID)
	require.NoError(t, err, "Should delete group successfully")

	// Verify bonds mapped to this group are deleted
	var mappingCount int
	_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM wish_isin WHERE wish_list_id = $1", targetWishlist.WishListID).Scan(&mappingCount)
	assert.Equal(t, 0, mappingCount, "All wish_isin mappings should be cascade deleted")
}
