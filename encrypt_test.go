// SPDX-License-Identifier: MIT

package sic

import (
	"bytes"
	"errors"
	"io"
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

func TestWriteByteArrayTooLarge(t *testing.T) {
	if err := writeByteArray(io.Discard, make([]byte, 256)); err == nil {
		t.Error("expected error for >255 bytes, got nil")
	}
}

func TestDecryptCorruptedZlib(t *testing.T) {
	enc, err := Encrypt([]byte("<database/>"), "foobar")
	if err != nil {
		t.Fatalf("could not encrypt: %v", err)
	}
	enc[len(enc)-1] ^= 0xFF
	_, err = Decrypt(bytes.NewReader(enc), "foobar")
	if err == nil {
		t.Fatal("expected error from corrupted body")
	}
	if errors.Is(err, ErrIncorrectPassword) {
		t.Errorf("corrupted body should not surface as ErrIncorrectPassword: %v", err)
	}
}

func TestMarshalRoundTripWithNewAttributes(t *testing.T) {
	in := &Database{
		Card: []Card{{
			ID:         "1",
			Title:      "T",
			Autofill:   "on",
			FirstStamp: "12345",
			Field: []Field{{
				Name:     "Login",
				Type:     "login",
				Text:     "u",
				Autofill: "username",
			}},
		}},
	}
	raw, err := Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out, err := Unmarshal(raw)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Card[0].Autofill != "on" || out.Card[0].FirstStamp != "12345" {
		t.Errorf("card attrs lost: %+v", out.Card[0])
	}
	if out.Card[0].Field[0].Autofill != "username" {
		t.Errorf("field autofill lost: %+v", out.Card[0].Field[0])
	}

	bare, err := Marshal(&Database{Card: []Card{{ID: "1"}}})
	if err != nil {
		t.Fatalf("marshal bare: %v", err)
	}
	if bytes.Contains(bare, []byte("autofill=")) || bytes.Contains(bare, []byte("first_stamp=")) {
		t.Errorf("omitempty failed; bare marshal contained new attributes:\n%s", bare)
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
