// Tideland Go UUID
//
// Copyright (C) 2021-2025 Frank Mueller / Tideland / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package uuid provides a comprehensive implementation of Universally Unique Identifiers (UUIDs)
// as defined in RFC 9562. It supports UUID versions 1, 3, 4, 5, 6, and 7 with full compliance
// to the specification.
//
// Key Features:
//   - RFC 9562 compliant implementation
//   - Time-ordered sortable UUIDs (v6, v7)
//   - Improved database index locality
//   - Concurrent-safe UUID generation
//   - Multiple string format support
//   - Name-based UUID generation
//
// UUID Versions:
//
// Version 1: Gregorian time-based UUID with MAC address. Legacy format, consider v6 or v7 instead.
//
// Version 3: MD5 name-based UUID. Use v5 instead for better security.
//
// Version 4: Random or pseudorandom UUID. Good for general-purpose unique identifiers.
//
// Version 5: SHA-1 name-based UUID. Preferred over v3 for name-based UUIDs.
//
// Version 6: Reordered Gregorian time-based UUID. Compatible with v1 but sortable.
//
// Version 7: Unix Epoch time-based UUID. Recommended for database keys and sortable IDs.
//
// Basic Usage:
//
//	// Create a random UUID (v4)
//	id := uuid.New()
//
//	// Create a time-based sortable UUID (v7) - recommended
//	id, err := uuid.NewV7()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Format as string
//	fmt.Println(id.String())      // "123e4567-e89b-12d3-a456-426614174000"
//	fmt.Println(id.ShortString()) // "123e4567e89b12d3a456426614174000"
//
// Database Usage:
//
// UUIDv7 is optimized for database usage with time-ordered values that provide
// better index locality:
//
//	type Record struct {
//		ID        uuid.UUID `db:"id"`
//		CreatedAt time.Time `db:"created_at"`
//		Data      string    `db:"data"`
//	}
//
//	func CreateRecord(data string) (*Record, error) {
//		id, err := uuid.NewV7()
//		if err != nil {
//			return nil, err
//		}
//		return &Record{
//			ID:        id,
//			CreatedAt: time.Now(),
//			Data:      data,
//		}, nil
//	}
//
// Name-Based UUIDs:
//
// Generate deterministic UUIDs from names within a namespace:
//
//	ns := uuid.NamespaceDNS()
//	id, err := uuid.NewV5(ns, []byte("www.example.com"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Same input always produces the same UUID
//
// Parsing and Validation:
//
// Parse UUIDs from various string formats:
//
//	// Standard format
//	id, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
//
//	// URN format
//	id, err = uuid.Parse("urn:uuid:123e4567-e89b-12d3-a456-426614174000")
//
//	// Braced format
//	id, err = uuid.Parse("{123e4567-e89b-12d3-a456-426614174000}")
//
//	// Short format (no dashes)
//	id, err = uuid.Parse("123e4567e89b12d3a456426614174000")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Check version and variant
//	fmt.Println(id.Version()) // e.g., uuid.V7
//	fmt.Println(id.Variant()) // uuid.VariantRFC4122
//
// Choosing a UUID Version:
//
//   - Use v7 for database primary keys, sortable IDs, or when creation time matters
//   - Use v4 for general-purpose unique identifiers when randomness is preferred
//   - Use v5/v3 when you need deterministic UUIDs from names
//   - Use v6 when you need v1 compatibility with sorting
//   - Use v1 only for legacy compatibility (consider v6 or v7 instead)
//
// Concurrency:
//
// All UUID generation functions are safe for concurrent use:
//
//	var wg sync.WaitGroup
//	for i := 0; i < 100; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			id, _ := uuid.NewV7()
//			fmt.Println(id)
//		}()
//	}
//	wg.Wait()
//
// Performance:
//
// UUIDv7 generation is highly performant and suitable for high-throughput applications.
// Benchmark results show generation rates of millions of UUIDs per second on modern hardware.
package uuid

// EOF
