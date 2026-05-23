[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/mxschmitt/golang-safe-in-cloud/issues)
[![GoDoc](https://godoc.org/github.com/mxschmitt/golang-safe-in-cloud?status.svg)](http://godoc.org/github.com/mxschmitt/golang-safe-in-cloud)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT)
[![Go Report](https://img.shields.io/badge/Go_report-A+-brightgreen.svg)](http://goreportcard.com/report/mxschmitt/golang-safe-in-cloud)
[![CI](https://github.com/mxschmitt/golang-safe-in-cloud/actions/workflows/ci.yml/badge.svg)](https://github.com/mxschmitt/golang-safe-in-cloud/actions/workflows/ci.yml)

# SafeInCloud Golang Decryption

Provides decryption of a SafeInCloud database in Golang.

# Example

```golang
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/mxschmitt/golang-safe-in-cloud"
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
    x, err := sic.Unmarshal(raw)
    if err != nil {
        log.Fatalf("could not unmarshal: %v", err)
    }
    fmt.Printf("data: %+v\n", x)
}
```
