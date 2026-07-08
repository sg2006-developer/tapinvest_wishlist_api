package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"tapinvest_api/models"
	"tapinvest_api/repository"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockWishlistRepo implements WishlistRepository for testing
type mockWishlistRepo struct {
	createErr error
	addErr    error
}

func (m *mockWishlistRepo) Create(ctx context.Context, name string) (*models.Wishlist, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &models.Wishlist{WishListID: "123e4567-e89b-12d3-a456-426614174000", WishListName: name, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}
func (m *mockWishlistRepo) GetAll(ctx context.Context) ([]models.WishlistResponse, error) { return nil, nil }
func (m *mockWishlistRepo) GetByID(ctx context.Context, id string, sortBy string) (*models.WishlistDetailResponse, error) { return nil, nil }
func (m *mockWishlistRepo) Update(ctx context.Context, id string, name string) (*models.Wishlist, error) { return nil, nil }
func (m *mockWishlistRepo) Delete(ctx context.Context, id string) error { return nil }
func (m *mockWishlistRepo) AddBond(ctx context.Context, wishlistID string, isin string) error { return m.addErr }
func (m *mockWishlistRepo) RemoveBond(ctx context.Context, wishlistID string, isin string) error { return nil }
func (m *mockWishlistRepo) SetBondColor(ctx context.Context, wishlistID string, isin string, color *string) error { return nil }
func (m *mockWishlistRepo) SetBondPosition(ctx context.Context, wishlistID string, isin string, position int) error { return nil }
func (m *mockWishlistRepo) SetBondPin(ctx context.Context, wishlistID string, isin string, isPinned bool) error { return nil }
func (m *mockWishlistRepo) ReorderBonds(ctx context.Context, wishlistID string, bondIsins []string) error { return nil }

func setupRouter(repo repository.WishlistRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := NewWishlistHandler(repo)
	r.POST("/wishlists", handler.CreateWishlist)
	r.POST("/wishlists/:wishlistId/bond", handler.AddBond)
	return r
}

func TestWishlistHandler_CreateWishlist(t *testing.T) {
	tests := []struct {
		name         string
		payload      map[string]string
		repoErr      error
		expectedCode int
	}{
		{
			name:         "Success",
			payload:      map[string]string{"name": "My New Wishlist"},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Empty Name",
			payload:      map[string]string{"name": ""},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Name Too Long",
			payload:      map[string]string{"name": "This name is definitely way more than twenty five characters"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Duplicate Name (Repo Error)",
			payload:      map[string]string{"name": "Duplicate"},
			repoErr:      repository.ErrWishlistNameExists,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Max 5 Wishlists Reached (Repo Error)",
			payload:      map[string]string{"name": "Valid Name"},
			repoErr:      repository.ErrMaxWishlistsReached,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockWishlistRepo{createErr: tt.repoErr}
			router := setupRouter(repo)

			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/wishlists", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestWishlistHandler_AddBond(t *testing.T) {
	repo := &mockWishlistRepo{addErr: repository.ErrBondDuplicate}
	router := setupRouter(repo)

	body, _ := json.Marshal(map[string]string{"bondIsin": "IN123456789"})
	req, _ := http.NewRequest(http.MethodPost, "/wishlists/123e4567-e89b-12d3-a456-426614174000/bond", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Since we mocked the repo to return ErrBondDuplicate, we expect 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "bond already exists")
}
