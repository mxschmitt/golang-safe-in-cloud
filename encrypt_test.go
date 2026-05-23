package sic

import (
	"bytes"
	"os"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	xml, err := os.ReadFile("decrypt_test.xml")
	if err != nil {
		t.Fatalf("could not read xml: %v", err)
	}
	enc, err := Encrypt(xml, "foobar")
	if err != nil {
		t.Fatalf("could not encrypt: %v", err)
	}
	dec, err := Decrypt(bytes.NewReader(enc), "foobar")
	if err != nil {
		t.Fatalf("could not decrypt: %v", err)
	}
	if !bytes.Equal(xml, dec) {
		t.Errorf("round-trip mismatch: %d in, %d out", len(xml), len(dec))
	}
}

func TestEncryptWrongPassword(t *testing.T) {
	enc, err := Encrypt([]byte("<database/>"), "right")
	if err != nil {
		t.Fatalf("could not encrypt: %v", err)
	}
	if _, err := Decrypt(bytes.NewReader(enc), "wrong"); err == nil {
		t.Errorf("expected decryption to fail with wrong password")
	}
}

func TestGenerateTestDB(t *testing.T) {
	enc, err := GenerateTestDB("foobar")
	if err != nil {
		t.Fatalf("could not generate: %v", err)
	}
	raw, err := Decrypt(bytes.NewReader(enc), "foobar")
	if err != nil {
		t.Fatalf("could not decrypt: %v", err)
	}
	db, err := Unmarshal(raw)
	if err != nil {
		t.Fatalf("could not unmarshal: %v", err)
	}
	if len(db.Card) != 1 || db.Card[0].ID != "100" || db.Card[0].Title != "Sample" {
		t.Errorf("unexpected card: %+v", db.Card)
	}
	if len(db.Card[0].Field) != 2 || db.Card[0].Field[0].Text != "user@example.com" {
		t.Errorf("unexpected field: %+v", db.Card[0].Field)
	}
	if len(db.Label) != 1 || db.Label[0].Name != "Test" {
		t.Errorf("unexpected label: %+v", db.Label)
	}
}
