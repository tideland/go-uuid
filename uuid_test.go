// Tideland Go UUID - Unit Tests
//
// Copyright (C) 2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package uuid_test // import "tideland.dev/go/uuid"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/uuid"
)

//--------------------
// TESTS
//--------------------

// TestStandard tests the standard UUID.
func TestStandard(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	// Asserts.
	uuidA := uuid.New()
	assert.Equal(uuidA.Version(), uuid.V4)
	uuidAShortStr := uuidA.ShortString()
	uuidAStr := uuidA.String()
	assert.Equal(len(uuidA), 16)
	assert.Match(uuidAShortStr, "[0-9a-f]{32}")
	assert.Match(uuidAStr, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
	// Check for copy.
	uuidB := uuid.New()
	uuidC := uuidB.Copy()
	for i := 0; i < len(uuidB); i++ {
		uuidB[i] = 0
	}
	assert.Different(uuidB, uuidC)
}

// TestVersions tests the creation of different UUID versions.
func TestVersions(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ns := uuid.NamespaceOID()
	name := []byte{1, 3, 3, 7}
	// Asserts.
	uuidV1, err := uuid.NewV1()
	assert.Nil(err)
	assert.Equal(uuidV1.Version(), uuid.V1)
	assert.Equal(uuidV1.Variant(), uuid.VariantRFC4122)
	assert.Logf("UUID V1: %v", uuidV1)
	uuidV3, err := uuid.NewV3(ns, name)
	assert.Nil(err)
	assert.Equal(uuidV3.Version(), uuid.V3)
	assert.Equal(uuidV3.Variant(), uuid.VariantRFC4122)
	assert.Logf("UUID V3: %v", uuidV3)
	uuidV4, err := uuid.NewV4()
	assert.Nil(err)
	assert.Equal(uuidV4.Version(), uuid.V4)
	assert.Equal(uuidV4.Variant(), uuid.VariantRFC4122)
	assert.Logf("UUID V4: %v", uuidV4)
	uuidV5, err := uuid.NewV5(ns, name)
	assert.Nil(err)
	assert.Equal(uuidV5.Version(), uuid.V5)
	assert.Equal(uuidV5.Variant(), uuid.VariantRFC4122)
	assert.Logf("UUID V5: %v", uuidV5)
}

// TestParse tests creating UUIDs from different string representations.
func TestParse(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ns := uuid.NamespaceOID()
	name := []byte{1, 3, 3, 7}
	// Asserts.
	tests := []struct {
		source  func() string
		version uuid.Version
		variant uuid.Variant
		err     string
	}{
		{func() string { u, _ := uuid.NewV1(); return u.String() }, uuid.V1, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV3(ns, name); return u.String() }, uuid.V3, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return u.String() }, uuid.V4, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV5(ns, name); return u.String() }, uuid.V5, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV1(); return "urn:uuid:" + u.String() }, uuid.V1, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV3(ns, name); return "urn:uuid:" + u.String() }, uuid.V3, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return "urn:uuid:" + u.String() }, uuid.V4, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV5(ns, name); return "urn:uuid:" + u.String() }, uuid.V5, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV1(); return "{" + u.String() + "}" }, uuid.V1, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV3(ns, name); return "{" + u.String() + "}" }, uuid.V3, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return "{" + u.String() + "}" }, uuid.V4, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV5(ns, name); return "{" + u.String() + "}" }, uuid.V5, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV1(); return u.ShortString() }, uuid.V1, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV3(ns, name); return u.ShortString() }, uuid.V3, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return u.ShortString() }, uuid.V4, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV5(ns, name); return u.ShortString() }, uuid.V5, uuid.VariantRFC4122, ""},
		{func() string { u, _ := uuid.NewV4(); return u.String() + "-ffaabb" }, 0, 0, "invalid source format"},
		{func() string { u, _ := uuid.NewV4(); return u.String() + "-ffxxyy" }, 0, 0, "invalid source format"},
		{func() string { u, _ := uuid.NewV4(); return "uuid:" + u.String() }, 0, 0, "invalid source format"},
		{func() string { u, _ := uuid.NewV4(); return "{" + u.ShortString() + "}" }, 0, 0, "invalid source format"},
		{func() string { return "ababababababababab" }, 0, 0, "invalid source format"},
		{func() string { return "abcdefabcdefZZZZefabcdefabcdefab" }, 0, 0, "source char 12 is no hex char"},
		{func() string { return "[abcdefabcdefabcdefabcdefabcdefab]" }, 0, 0, "invalid source format"},
		{func() string { return "abcdefab=cdef=abcd=efab=cdefabcdefab" }, 0, 0, "source char 8 does not match pattern"},
	}
	for i, test := range tests {
		source := test.source()
		assert.Logf("test #%d source %s", i, source)
		uuidT, err := uuid.Parse(source)
		if test.err == "" {
			assert.NoError(err)
			assert.Equal(uuidT.Version(), test.version)
			assert.Equal(uuidT.Variant(), test.variant)
		} else {
			assert.ErrorContains(err, test.err)
		}
	}
}

// EOF
