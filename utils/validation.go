package utils

import (
	"strings"
)

// ValidateWishlistName checks if the name is valid (<= 25 chars and not empty after trim)
func ValidateWishlistName(name string) (string, bool) {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return "", false
	}
	if len(trimmed) > 25 {
		return "", false
	}
	return trimmed, true
}
