package publicqa

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNormalizeTitle(t *testing.T) {
	v, e := NormalizeTitle("  Cost   OF living ")
	if e != nil || v != "cost of living" {
		t.Fatalf("got %q, %v", v, e)
	}
	if _, e = NormalizeTitle("https://spam.example"); e == nil {
		t.Fatal("expected URL rejection")
	}
}
func TestPublicIssueJSONDoesNotLeakCreatorUID(t *testing.T) {
	b, err := json.Marshal(Issue{ID: "i", CreatorUID: "private-user", Attribution: "anonymous"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "private-user") || strings.Contains(string(b), "creatorUID") {
		t.Fatalf("private identity leaked: %s", b)
	}
}
func TestSlugifyStable(t *testing.T) {
	if got := Slugify("Housing Affordability!"); got != "housing-affordability" {
		t.Fatalf("got %q", got)
	}
}
