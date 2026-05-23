// SPDX-License-Identifier: MIT

package sic

import (
	"encoding/xml"
	"fmt"
)

// Database is the root SafeInCloud database, holding labels, templates, cards,
// notes, attachments, and tombstones for deleted entries.
type Database struct {
	XMLName xml.Name `xml:"database"`
	Notes   []string `xml:"notes"`
	LabelID []string `xml:"label_id"`
	File    [][]File `xml:"file"`
	Ghost   []Ghost  `xml:"ghost"`
	Label   []Label  `xml:"label"`
	Card    []Card   `xml:"card"`
}

// Ghost is a tombstone for a deleted card, retained so sync can propagate the
// deletion to other clients.
type Ghost struct {
	ID        string `xml:"id,attr"`
	TimeStamp string `xml:"time_stamp,attr"`
}

// Label is a user- or schema-defined tag that can be attached to cards.
type Label struct {
	Type      string `xml:"type,attr"`
	TimeStamp string `xml:"time_stamp,attr"`
	ID        string `xml:"id,attr"`
	Name      string `xml:"name,attr"`
}

// Card is a single SafeInCloud entry such as a login, credit card, or note.
// Cards with template="true" define the schema for entries created from them.
type Card struct {
	ID          string  `xml:"id,attr"`
	Symbol      string  `xml:"symbol,attr"`
	Template    string  `xml:"template,attr"`
	Type        string  `xml:"type,attr"`
	WebsiteIcon string  `xml:"website_icon,attr"`
	TimeStamp   string  `xml:"time_stamp,attr"`
	FirstStamp  string  `xml:"first_stamp,attr,omitempty"`
	Deleted     string  `xml:"deleted,attr"`
	Title       string  `xml:"title,attr"`
	Color       string  `xml:"color,attr"`
	Star        string  `xml:"star,attr"`
	Autofill    string  `xml:"autofill,attr,omitempty"`
	Field       []Field `xml:"field"`
}

// Field is a single labelled value within a Card (e.g. a username, password,
// or PIN).
type Field struct {
	Hash     string `xml:"hash,attr"`
	History  string `xml:"history,attr"`
	Name     string `xml:"name,attr"`
	Type     string `xml:"type,attr"`
	Text     string `xml:",chardata"`
	Score    string `xml:"score,attr"`
	Autofill string `xml:"autofill,attr,omitempty"`
}

// File is a binary attachment encoded inline within a Card.
type File struct {
	Name string `xml:"name,attr"`
	Text string `xml:",chardata"`
}

// Unmarshal converts the xml in []byte into a Go struct
func Unmarshal(raw []byte) (*Database, error) {
	var db *Database
	if err := xml.Unmarshal(raw, &db); err != nil {
		return nil, fmt.Errorf("could not Unmarshal xml: %w", err)
	}
	return db, nil
}
