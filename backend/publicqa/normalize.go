package publicqa

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var errInvalidTitle = errors.New("issue title must be 3 to 120 characters and contain no URL or control characters")
var urlPattern = regexp.MustCompile(`(?i)(https?://|www\.)`)
var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func NormalizeTitle(title string) (string, error) {
	title = strings.Join(strings.Fields(strings.TrimSpace(title)), " ")
	if len([]rune(title)) < 3 || len([]rune(title)) > 120 || urlPattern.MatchString(title) {
		return "", errInvalidTitle
	}
	for _, r := range title {
		if unicode.IsControl(r) {
			return "", errInvalidTitle
		}
	}
	return strings.ToLower(title), nil
}
func normalizedHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:16])
}
func NormalizedKey(s string) string { return normalizedHash(s) }
func Slugify(title string) string {
	s := strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(strings.TrimSpace(title)), "-"), "-")
	if len(s) > 72 {
		s = strings.Trim(s[:72], "-")
	}
	if s == "" {
		return "issue"
	}
	return s
}
