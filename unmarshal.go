package sic

import (
	"encoding/xml"
	"fmt"
)

type Database struct {
	XMLName xml.Name `xml:"database"`
	Notes   []string `xml:"notes"`
	LabelID []string `xml:"label_id"`
	File    [][]File `xml:"file"`
	Ghost   []Ghost  `xml:"ghost"`
	Label   []Label  `xml:"label"`
	Card    []Card   `xml:"card"`
}

type Ghost struct {
	ID        string `xml:"id,attr"`
	TimeStamp string `xml:"time_stamp,attr"`
}

type Label struct {
	Type      string `xml:"type,attr"`
	TimeStamp string `xml:"time_stamp,attr"`
	ID        string `xml:"id,attr"`
	Name      string `xml:"name,attr"`
}

type Card struct {
	ID          string  `xml:"id,attr"`
	Symbol      string  `xml:"symbol,attr"`
	Template    string  `xml:"template,attr"`
	Type        string  `xml:"type,attr"`
	WebsiteIcon string  `xml:"website_icon,attr"`
	TimeStamp   string  `xml:"time_stamp,attr"`
	Deleted     string  `xml:"deleted,attr"`
	Title       string  `xml:"title,attr"`
	Color       string  `xml:"color,attr"`
	Star        string  `xml:"star,attr"`
	Field       []Field `xml:"field"`
}

type Field struct {
	Hash    string `xml:"hash,attr"`
	History string `xml:"history,attr"`
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Text    string `xml:",chardata"`
	Score   string `xml:"score,attr"`
}

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
