package sic

import "fmt"

// GenerateTestDB builds a small SafeInCloud database with sample data and
// returns it encrypted under the given password. Useful for tests that need
// a realistic .db payload without depending on a checked-in fixture.
func GenerateTestDB(password string) ([]byte, error) {
	db := &Database{
		Label: []Label{
			{ID: "1", Name: "Test"},
		},
		Card: []Card{
			{
				ID:    "100",
				Title: "Sample",
				Field: []Field{
					{Name: "Login", Type: "login", Text: "user@example.com"},
					{Name: "Password", Type: "password", Text: "hunter2"},
				},
			},
		},
	}
	raw, err := Marshal(db)
	if err != nil {
		return nil, fmt.Errorf("could not marshal sample db: %w", err)
	}
	return Encrypt(raw, password)
}
