package utils

import (
	"regexp"
	"strings"
)

// GenerateSlug merubah "Halo Dunia 2026!" menjadi "halo-dunia-2026"
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}