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
func TestGeographicHierarchyJSONIsPreserved(t *testing.T) {
	q := Question{ID: "city-dublin", Publication: StatusPublished, Scope: Scope{ID: "DUB", Type: "city", Name: "Dublin", ParentID: "IE-D"}}
	c := Category{ID: "city", ParentCategoryID: "county", Publication: StatusPublished}
	b, e := json.Marshal(struct {
		Question Question `json:"question"`
		Category Category `json:"category"`
	}{q, c})
	if e != nil {
		t.Fatal(e)
	}
	s := string(b)
	for _, want := range []string{`"scope":{"id":"DUB","type":"city","name":"Dublin","parentId":"IE-D"}`, `"parentCategoryId":"county"`, `"publication":"published"`} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %s in %s", want, s)
		}
	}
}
func TestChoiceSourceValidation(t *testing.T) {
	if e := validateChoiceSource(ChoiceSource{Kind: "custom", Options: []ChoiceOption{{Title: "Housing"}, {Title: " housing "}}}); e == nil {
		t.Fatal("expected normalized duplicate rejection")
	}
	if e := validateChoiceSource(ChoiceSource{Kind: "predefined", EntityType: "city"}); e != nil {
		t.Fatal(e)
	}
	if e := validateChoiceSource(ChoiceSource{Kind: "predefined", EntityType: "ethnicity"}); e == nil {
		t.Fatal("unexpected unsupported identity dimension")
	}
}
func TestSlugifyStable(t *testing.T) {
	if got := Slugify("Housing Affordability!"); got != "housing-affordability" {
		t.Fatalf("got %q", got)
	}
	if got := Slugify("!@#$%^&*()"); got != "issue" {
		t.Fatalf("expected 'issue' for non-alphanumeric, got %q", got)
	}
	longTitle := strings.Repeat("long-title-with-many-words-", 10)
	if got := Slugify(longTitle); len(got) > 72 || !strings.HasPrefix(got, "long-title") {
		t.Fatalf("expected truncated slug <= 72 chars, got %q (len %d)", got, len(got))
	}
}

func TestNormalizeTitle_Branches(t *testing.T) {
	// Too short (< 3)
	if _, err := NormalizeTitle("ab"); err == nil {
		t.Fatal("expected error for title < 3 chars")
	}
	// Too long (> 120)
	if _, err := NormalizeTitle(strings.Repeat("a", 121)); err == nil {
		t.Fatal("expected error for title > 120 chars")
	}
	// Control character
	if _, err := NormalizeTitle("valid \x00 title"); err == nil {
		t.Fatal("expected error for control character")
	}
}

func TestNormalizedKey(t *testing.T) {
	k1 := NormalizedKey("hello world")
	k2 := NormalizedKey("hello world")
	if k1 == "" || k1 != k2 {
		t.Fatalf("expected stable non-empty key, got %q vs %q", k1, k2)
	}
}
