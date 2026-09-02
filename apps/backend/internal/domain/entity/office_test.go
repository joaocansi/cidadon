package entity

import "testing"

func TestOfficeSlugIsReadableNormalizedAndUniquePerCouncillor(t *testing.T) {
	first := OfficeSlug("Partido Verde", "Ana de Ávila", 41)
	second := OfficeSlug("Partido Verde", "Ana de Ávila", 42)

	if first != "partido-verde-ana-de-avila-41" {
		t.Fatalf("unexpected slug: %q", first)
	}
	if first == second {
		t.Fatal("slugs must differ for distinct councillors with identical party and name")
	}
}
