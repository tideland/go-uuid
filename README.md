# Tideland Go UUID

[![GitHub release](https://img.shields.io/github/release/tideland/go-uuid.svg)](https://github.com/tideland/go-uuid)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-uuid/master/LICENSE)
[![Go Module](https://img.shields.io/github/go-mod/go-version/tideland/go-uuid)](https://github.com/tideland/go-uuid/blob/master/go.mod)
[![GoDoc](https://godoc.org/tideland.dev/go/uuid?status.svg)](https://pkg.go.dev/mod/tideland.dev/go/uuid?tab=packages)
![Workflow](https://github.com/tideland/go-uuid/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/go-uuid)](https://goreportcard.com/report/tideland.dev/go/uuid)

## Description

**Tideland Go UUID** provides functions for the creation and working with UUIDs in versions
1, 3, 4, 5, 6, and 7 as per RFC 9562. The package supports:

- **Version 1**: Gregorian time-based with MAC address
- **Version 3**: MD5 name-based
- **Version 4**: Random/pseudorandom
- **Version 5**: SHA-1 name-based
- **Version 6**: Reordered Gregorian time-based (sortable)
- **Version 7**: Unix Epoch time-based (sortable, recommended)

## Features

- Full RFC 9562 compliance
- Time-ordered sortable UUIDs (v6, v7)
- Improved database index locality with v7
- Concurrent-safe UUID generation
- Multiple string formats (standard, short, URN, braced)
- Comprehensive test coverage

## Installation

```bash
go get tideland.dev/go/uuid
```

## Usage

### Creating UUIDs

```go
import "tideland.dev/go/uuid"

// Version 4 (random) - default
id := uuid.New()

// Version 7 (time-based, sortable) - recommended for databases
id, err := uuid.NewV7()

// Version 6 (time-based, sortable)
id, err := uuid.NewV6()

// Version 1 (time-based with MAC)
id, err := uuid.NewV1()

// Version 5 (name-based with SHA-1)
ns := uuid.NamespaceDNS()
id, err := uuid.NewV5(ns, []byte("www.example.com"))

// Version 4 (random)
id, err := uuid.NewV4()

// Version 3 (name-based with MD5)
id, err := uuid.NewV3(ns, []byte("www.example.com"))
```

### Parsing and Formatting

```go
// Parse from string
id, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

// Format as string
str := id.String()  // "123e4567-e89b-12d3-a456-426614174000"
short := id.ShortString()  // "123e4567e89b12d3a456426614174000"

// Get version and variant
version := id.Version()
variant := id.Variant()
```

### Namespaces

```go
// Predefined namespaces for name-based UUIDs
ns := uuid.NamespaceDNS()   // DNS namespace
ns := uuid.NamespaceURL()   // URL namespace
ns := uuid.NamespaceOID()   // OID namespace
ns := uuid.NamespaceX500()  // X.500 namespace
```

## Choosing a UUID Version

- **Use v7** for database primary keys, sortable IDs, or when creation time matters
- **Use v4** for general-purpose unique identifiers when randomness is preferred
- **Use v5/v3** when you need deterministic UUIDs from names
- **Use v6** when you need v1 compatibility with sorting
- **Use v1** only for legacy compatibility (consider v6 or v7 instead)

I hope you like it. ;)

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://tideland.dev)
