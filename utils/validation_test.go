package utils

import (
	"testing"
)

func TestValidateWishlistName(t *testing.T) {
	// Covering scenarios from QA Section 1: Wishlist Name Limiting
	tests := []struct {
		name     string
		input    string
		expected string
		isValid  bool
	}{
		{"Valid Name", "My Wishlist", "My Wishlist", true},
		{"Emojis (currently allowed)", "My Wishlist 🚀", "My Wishlist 🚀", true},
		{"Starts with underscore (allowed)", "_MyWishlist", "_MyWishlist", true},
		{"Underscores in middle/end", "My_Wishlist_", "My_Wishlist_", true},
		{"Special characters (allowed)", "My @#$% Wishlist", "My @#$% Wishlist", true},
		{"Numbers (allowed)", "Wishlist 123", "Wishlist 123", true},
		{"Spaces between words (allowed)", "My   Wishlist", "My   Wishlist", true},
		{"Starts with a space (trimmed)", "  My Wishlist", "My Wishlist", true},
		{"Only spaces", "     ", "", false},
		{"Blank/Empty", "", "", false},
		{"Exceeds max limit (25 chars)", "This wishlist name is way too long", "", false},
		{"Exactly 25 chars", "Exactly twenty five chars", "Exactly twenty five chars", true},
		{"One character long", "A", "A", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, valid := ValidateWishlistName(tt.input)
			if valid != tt.isValid {
				t.Errorf("expected isValid %v, got %v for input %q", tt.isValid, valid, tt.input)
			}
			if valid && result != tt.expected {
				t.Errorf("expected result %q, got %q", tt.expected, result)
			}
		})
	}
}
