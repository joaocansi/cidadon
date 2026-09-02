package entity

import "testing"

func TestOfficialPartyReturnsCanonicalAcronym(t *testing.T) {
	party, ok := OfficialParty(" psol ")
	if !ok || party.Acronym != "PSOL" {
		t.Fatalf("expected PSOL to be in the official catalogue, got %#v (ok=%v)", party, ok)
	}
}

func TestOfficialPartyRejectsFreeFormNames(t *testing.T) {
	if _, ok := OfficialParty("Partido inventado"); ok {
		t.Fatal("free-form party names must not be accepted")
	}
}
