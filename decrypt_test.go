// SPDX-License-Identifier: MIT

package sic

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestDecryptionSuccess(t *testing.T) {
	file, err := os.Open("decrypt_test.db")
	if err != nil {
		t.Errorf("could not read file: %v", err)
	}
	raw, err := Decrypt(file, "foobar")
	if err != nil {
		t.Errorf("could not decrypt: %v", err)
	}
	db, err := Unmarshal(raw)
	if err != nil {
		t.Errorf("could not unmarshal xml: %v", err)
	}
	tt := []struct {
		expected string
		got      string
		key      string
	}{
		{
			expected: "101",
			got:      db.Card[0].ID,
			key:      "database>card[0].id",
		},
		{
			expected: "Nummer",
			got:      db.Card[0].Field[0].Name,
			key:      "database>card[0]>field[0].name",
		},
		{
			expected: "Geschäftlich",
			got:      db.Label[0].Name,
			key:      "database>label[0].name",
		},
		{
			expected: "3",
			got:      fmt.Sprintf("%d", len(db.Card[1].Field)),
			key:      "len(database>card[1].field)",
		},
	}
	for _, tc := range tt {
		if tc.expected != tc.got {
			t.Errorf("expected: %s; got: %s; diffs: %s", tc.expected, tc.got, tc.key)
		}
	}
}

func TestDecryptionInvalidPassword(t *testing.T) {
	file, err := os.Open("decrypt_test.db")
	if err != nil {
		t.Errorf("could not read file: %v", err)
	}
	if _, err = Decrypt(file, "definitely not correct"); !errors.Is(err, ErrIncorrectPassword) {
		t.Errorf("expected ErrIncorrectPassword, got: %v", err)
	}
}
