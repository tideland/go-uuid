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
	// Asserts.
	uuidV1, err := uuid.NewV1()
	assert.Nil(err)
	assert.Equal(uuidV1.Version(), uuid.V1)
	assert.Equal(uuidV1.Variant(), uuid.VariantRFC4122)
	assert.Logf("UUID V1: %v", uuidV1)
	uuidV3, err := uuid.NewV3(ns, []byte{4, 7, 1, 1})
	assert.Nil(err)
	assert.Equal(uuidV3.Version(), uuid.V3)
	assert.Equal(uuidV3.Variant(), uuid.VariantRFC4122)
	assert.Logf("UUID V3: %v", uuidV3)
	uuidV4, err := uuid.NewV4()
	assert.Nil(err)
	assert.Equal(uuidV4.Version(), uuid.V4)
	assert.Equal(uuidV4.Variant(), uuid.VariantRFC4122)
	assert.Logf("UUID V4: %v", uuidV4)
	uuidV5, err := uuid.NewV5(ns, []byte{4, 7, 1, 1})
	assert.Nil(err)
	assert.Equal(uuidV5.Version(), uuid.V5)
	assert.Equal(uuidV5.Variant(), uuid.VariantRFC4122)
	assert.Logf("UUID V5: %v", uuidV5)
}

// TestFromHex tests creating UUIDs from hex strings.
func TestFromHex(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	// Asserts.
	_, err := uuid.FromHex("ffff")
	assert.ErrorMatch(err, `source length is not 32`)
	_, err = uuid.FromHex("012345678901234567890123456789zz")
	assert.ErrorMatch(err, `source is no hex value: .*`)
	_, err = uuid.FromHex("012345678901234567890123456789ab")
	assert.Nil(err)
}

// EOF
