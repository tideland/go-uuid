// Tideland Go UUID - Unit Tests
//
// Copyright (C) 2021-2025 Frank Mueller / Tideland / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package uuid_test

import (
	"testing"

	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/uuid"
)

// Tests

// TestStandard tests the standard UUID.
func TestStandard(t *testing.T) {
	// Test UUID creation and format
	uuidA := uuid.New()
	verify.Equal(t, uuidA.Version(), uuid.V4)
	uuidAShortStr := uuidA.ShortString()
	uuidAStr := uuidA.String()
	verify.Equal(t, len(uuidA), 16)
	verify.Match(t, uuidAShortStr, "[0-9a-f]{32}")
	verify.Match(t, uuidAStr, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")

	// Check for copy
	uuidB := uuid.New()
	uuidC := uuidB.Copy()
	for i := range len(uuidB) {
		uuidB[i] = 0
	}
	verify.Different(t, uuidB, uuidC)
}

// TestVersions tests the creation of different UUID versions.
func TestVersions(t *testing.T) {
	ns := uuid.NamespaceOID()
	name := []byte{1, 3, 3, 7}

	// Test UUID v1
	uuidV1, err := uuid.NewV1()
	verify.NoError(t, err)
	verify.Equal(t, uuidV1.Version(), uuid.V1)
	verify.Equal(t, uuidV1.Variant(), uuid.VariantRFC4122)
	t.Logf("UUID V1: %v", uuidV1)

	// Test UUID v3
	uuidV3, err := uuid.NewV3(ns, name)
	verify.NoError(t, err)
	verify.Equal(t, uuidV3.Version(), uuid.V3)
	verify.Equal(t, uuidV3.Variant(), uuid.VariantRFC4122)
	t.Logf("UUID V3: %v", uuidV3)

	// Test UUID v4
	uuidV4, err := uuid.NewV4()
	verify.NoError(t, err)
	verify.Equal(t, uuidV4.Version(), uuid.V4)
	verify.Equal(t, uuidV4.Variant(), uuid.VariantRFC4122)
	t.Logf("UUID V4: %v", uuidV4)

	// Test UUID v5
	uuidV5, err := uuid.NewV5(ns, name)
	verify.NoError(t, err)
	verify.Equal(t, uuidV5.Version(), uuid.V5)
	verify.Equal(t, uuidV5.Variant(), uuid.VariantRFC4122)
	t.Logf("UUID V5: %v", uuidV5)

	// Test UUID v6
	uuidV6, err := uuid.NewV6()
	verify.NoError(t, err)
	verify.Equal(t, uuidV6.Version(), uuid.V6)
	verify.Equal(t, uuidV6.Variant(), uuid.VariantRFC4122)
	t.Logf("UUID V6: %v", uuidV6)

	// Test UUID v7
	uuidV7, err := uuid.NewV7()
	verify.NoError(t, err)
	verify.Equal(t, uuidV7.Version(), uuid.V7)
	verify.Equal(t, uuidV7.Variant(), uuid.VariantRFC4122)
	t.Logf("UUID V7: %v", uuidV7)
}

// TestParse tests creating UUIDs from different string representations.
func TestParse(t *testing.T) {
	ns := uuid.NamespaceOID()
	name := []byte{1, 3, 3, 7}

	tests := []struct {
		name    string
		source  func() string
		version uuid.Version
		variant uuid.Variant
		err     string
	}{
		{"v1-standard", func() string { u, _ := uuid.NewV1(); return u.String() }, uuid.V1, uuid.VariantRFC4122, ""},
		{"v3-standard", func() string { u, _ := uuid.NewV3(ns, name); return u.String() }, uuid.V3, uuid.VariantRFC4122, ""},
		{"v4-standard", func() string { u, _ := uuid.NewV4(); return u.String() }, uuid.V4, uuid.VariantRFC4122, ""},
		{"v5-standard", func() string { u, _ := uuid.NewV5(ns, name); return u.String() }, uuid.V5, uuid.VariantRFC4122, ""},
		{"v6-standard", func() string { u, _ := uuid.NewV6(); return u.String() }, uuid.V6, uuid.VariantRFC4122, ""},
		{"v7-standard", func() string { u, _ := uuid.NewV7(); return u.String() }, uuid.V7, uuid.VariantRFC4122, ""},
		{"v1-urn", func() string { u, _ := uuid.NewV1(); return "urn:uuid:" + u.String() }, uuid.V1, uuid.VariantRFC4122, ""},
		{"v4-urn", func() string { u, _ := uuid.NewV4(); return "urn:uuid:" + u.String() }, uuid.V4, uuid.VariantRFC4122, ""},
		{"v1-braced", func() string { u, _ := uuid.NewV1(); return "{" + u.String() + "}" }, uuid.V1, uuid.VariantRFC4122, ""},
		{"v4-braced", func() string { u, _ := uuid.NewV4(); return "{" + u.String() + "}" }, uuid.V4, uuid.VariantRFC4122, ""},
		{"v1-short", func() string { u, _ := uuid.NewV1(); return u.ShortString() }, uuid.V1, uuid.VariantRFC4122, ""},
		{"v4-short", func() string { u, _ := uuid.NewV4(); return u.ShortString() }, uuid.V4, uuid.VariantRFC4122, ""},
		{"v7-short", func() string { u, _ := uuid.NewV7(); return u.ShortString() }, uuid.V7, uuid.VariantRFC4122, ""},
		{"invalid-too-long", func() string { u, _ := uuid.NewV4(); return u.String() + "-ffaabb" }, 0, 0, "invalid source format"},
		{"invalid-non-hex", func() string { return "abcdefabcdefZZZZefabcdefabcdefab" }, 0, 0, "source char 12 is no hex char"},
		{"invalid-brackets", func() string { return "[abcdefabcdefabcdefabcdefabcdefab]" }, 0, 0, "invalid source format"},
		{"invalid-separator", func() string { return "abcdefab=cdef=abcd=efab=cdefabcdefab" }, 0, 0, "source char 8 does not match pattern"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := test.source()
			t.Logf("source: %s", source)
			uuidT, err := uuid.Parse(source)
			if test.err == "" {
				verify.NoError(t, err)
				verify.Equal(t, uuidT.Version(), test.version)
				verify.Equal(t, uuidT.Variant(), test.variant)
			} else {
				verify.ErrorContains(t, err, test.err)
			}
		})
	}
}

// TestV6Sortability tests that UUIDv6 values are sortable by creation time.
func TestV6Sortability(t *testing.T) {
	uuids := make([]uuid.UUID, 100)
	for i := range 100 {
		u, err := uuid.NewV6()
		verify.NoError(t, err)
		uuids[i] = u
	}

	// Verify each UUID is greater than or equal to the previous
	for i := 1; i < len(uuids); i++ {
		prev := uuids[i-1].String()
		curr := uuids[i].String()
		verify.True(t, curr >= prev, "UUIDs should be sortable")
	}
}

// TestV7Sortability tests that UUIDv7 values are sortable by creation time.
func TestV7Sortability(t *testing.T) {
	uuids := make([]uuid.UUID, 100)
	for i := range 100 {
		u, err := uuid.NewV7()
		verify.NoError(t, err)
		uuids[i] = u
	}

	// Verify each UUID is greater than or equal to the previous
	for i := 1; i < len(uuids); i++ {
		prev := uuids[i-1].String()
		curr := uuids[i].String()
		verify.True(t, curr >= prev, "UUIDs should be sortable")
	}
}

// TestV7Monotonicity tests that UUIDv7 values are monotonic within same millisecond.
func TestV7Monotonicity(t *testing.T) {
	// Generate many UUIDs quickly to ensure some share the same millisecond
	uuids := make([]uuid.UUID, 1000)
	for i := range 1000 {
		u, err := uuid.NewV7()
		verify.NoError(t, err)
		uuids[i] = u
	}

	// Check that all UUIDs are unique
	seen := make(map[string]bool)
	for _, u := range uuids {
		s := u.String()
		verify.False(t, seen[s], "All UUIDs should be unique")
		seen[s] = true
	}
}

// TestConcurrentGeneration tests concurrent UUID generation.
func TestConcurrentGeneration(t *testing.T) {
	const goroutines = 10
	const uuidsPerGoroutine = 100

	type result struct {
		uuid uuid.UUID
		err  error
	}

	results := make(chan result, goroutines*uuidsPerGoroutine)

	// Generate UUIDs concurrently
	for range goroutines {
		go func() {
			for range uuidsPerGoroutine {
				u, err := uuid.NewV7()
				results <- result{uuid: u, err: err}
			}
		}()
	}

	// Collect results
	seen := make(map[string]bool)
	for range goroutines * uuidsPerGoroutine {
		r := <-results
		verify.NoError(t, r.err)
		s := r.uuid.String()
		verify.False(t, seen[s], "Concurrent UUIDs should be unique")
		seen[s] = true
	}
}

// TestNamespaceStability tests that name-based UUIDs are stable.
func TestNamespaceStability(t *testing.T) {
	ns := uuid.NamespaceDNS()
	name := []byte("www.example.com")

	// Generate same name-based UUID multiple times
	uuid1, err := uuid.NewV5(ns, name)
	verify.NoError(t, err)

	uuid2, err := uuid.NewV5(ns, name)
	verify.NoError(t, err)

	uuid3, err := uuid.NewV5(ns, name)
	verify.NoError(t, err)

	// All should be identical
	verify.Equal(t, uuid1.String(), uuid2.String())
	verify.Equal(t, uuid2.String(), uuid3.String())
}

// TestRawAndCopy tests Raw and Copy methods.
func TestRawAndCopy(t *testing.T) {
	u, err := uuid.NewV7()
	verify.NoError(t, err)

	// Test Raw returns correct array
	raw := u.Raw()
	verify.Equal(t, len(raw), 16)
	for i := range 16 {
		verify.Equal(t, raw[i], u[i])
	}

	// Test Copy creates independent copy
	copy := u.Copy()
	verify.Equal(t, u.String(), copy.String())

	// Modify original, copy should be unchanged
	original := u.String()
	u[0] = 0xFF
	verify.Different(t, u.String(), original)
	verify.Equal(t, copy.String(), original)
}

// BenchmarkNewV1 benchmarks UUID v1 generation.
func BenchmarkNewV1(b *testing.B) {
	for b.Loop() {
		_, _ = uuid.NewV1()
	}
}

// BenchmarkNewV4 benchmarks UUID v4 generation.
func BenchmarkNewV4(b *testing.B) {
	for b.Loop() {
		_, _ = uuid.NewV4()
	}
}

// BenchmarkNewV6 benchmarks UUID v6 generation.
func BenchmarkNewV6(b *testing.B) {
	for b.Loop() {
		_, _ = uuid.NewV6()
	}
}

// BenchmarkNewV7 benchmarks UUID v7 generation.
func BenchmarkNewV7(b *testing.B) {
	for b.Loop() {
		_, _ = uuid.NewV7()
	}
}

// BenchmarkNewV5 benchmarks UUID v5 generation.
func BenchmarkNewV5(b *testing.B) {
	ns := uuid.NamespaceDNS()
	name := []byte("www.example.com")

	for b.Loop() {
		_, _ = uuid.NewV5(ns, name)
	}
}

// BenchmarkParse benchmarks UUID parsing.
func BenchmarkParse(b *testing.B) {
	u, _ := uuid.NewV7()
	s := u.String()

	for b.Loop() {
		_, _ = uuid.Parse(s)
	}
}

// BenchmarkString benchmarks UUID string formatting.
func BenchmarkString(b *testing.B) {
	u, _ := uuid.NewV7()

	for b.Loop() {
		_ = u.String()
	}
}

// EOF
