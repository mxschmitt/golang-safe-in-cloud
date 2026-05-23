// SPDX-License-Identifier: MIT

package sic

import (
	"encoding/xml"
	"fmt"
)

// Marshal converts a Database into the XML representation expected by Encrypt.
func Marshal(db *Database) ([]byte, error) {
	body, err := xml.MarshalIndent(db, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("could not Marshal xml: %w", err)
	}
	out := make([]byte, 0, len(xml.Header)+len(body))
	out = append(out, []byte(xml.Header)...)
	out = append(out, body...)
	return out, nil
}
