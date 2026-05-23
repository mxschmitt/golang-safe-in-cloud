[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/mxschmitt/golang-safe-in-cloud/issues)
[![Go Reference](https://pkg.go.dev/badge/github.com/mxschmitt/golang-safe-in-cloud.svg)](https://pkg.go.dev/github.com/mxschmitt/golang-safe-in-cloud)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT)
[![Go Report](https://img.shields.io/badge/Go_report-A+-brightgreen.svg)](http://goreportcard.com/report/mxschmitt/golang-safe-in-cloud)
[![CI](https://github.com/mxschmitt/golang-safe-in-cloud/actions/workflows/ci.yml/badge.svg)](https://github.com/mxschmitt/golang-safe-in-cloud/actions/workflows/ci.yml)

# SafeInCloud Golang

Encrypt and decrypt [SafeInCloud](https://www.safe-in-cloud.com) database files in pure Go. Zero external dependencies; Go 1.24+.

## Decrypt

```go
package main

import (
    "fmt"
    "log"
    "os"

    sic "github.com/mxschmitt/golang-safe-in-cloud"
)

func main() {
    file, err := os.Open("SafeInCloud.db")
    if err != nil {
        log.Fatalf("could not read file: %v", err)
    }
    raw, err := sic.Decrypt(file, "foobar")
    if err != nil {
        log.Fatalf("could not decrypt: %v", err)
    }
    fmt.Println(string(raw))
    db, err := sic.Unmarshal(raw)
    if err != nil {
        log.Fatalf("could not unmarshal: %v", err)
    }
    fmt.Printf("data: %+v\n", db)
}
```

## Encrypt

```go
db := &sic.Database{
    Card: []sic.Card{{
        ID:    "1",
        Title: "Example",
        Field: []sic.Field{{Name: "Login", Type: "login", Text: "user@example.com"}},
    }},
}
raw, err := sic.Marshal(db)
if err != nil {
    log.Fatal(err)
}
enc, err := sic.Encrypt(raw, "foobar")
if err != nil {
    log.Fatal(err)
}
if err := os.WriteFile("SafeInCloud.db", enc, 0o600); err != nil {
    log.Fatal(err)
}
```

For tests, `sic.GenerateTestDB(password)` returns an encrypted database with sample data.

## Note

The on-disk format uses PBKDF2-SHA1 (10,000 iterations) and AES-CBC as required by SafeInCloud's file format; this is dictated by interoperability with the official client, not by current cryptographic best practice.
