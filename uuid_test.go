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
	for i := 0; i < len(uuidB); i++ {
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

// EOF
