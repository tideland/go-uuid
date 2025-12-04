// Tideland Go UUID
//
// Copyright (C) 2021-2025 Frank Mueller / Tideland / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package uuid

//--------------------
// IMPORTS
//--------------------

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

//--------------------
// UUID
//--------------------

// Version represents a UUID's version.
type Version byte

// UUID versions and variants.
const (
	V1 Version = 1
	V3 Version = 3
	V4 Version = 4
	V5 Version = 5
)

// Variant represents a UUID's variant.
type Variant byte

const (
	VariantNCS       Variant = 0 // Reserved, NCS backward compatibility.
	VariantRFC4122   Variant = 4 // The variant specified in RFC4122.
	VariantMicrosoft Variant = 6 // Reserved, Microsoft Corporation backward compatibility.
	VariantFuture    Variant = 7 // Reserved for future definition.
)

// UUID represents a universal identifier with 16 bytes.
// See http://en.wikipedia.org/wiki/Universally_unique_identifier.
type UUID [16]byte

// New returns a new UUID with based on the default version 4.
func New() UUID {
	uuid, err := NewV4()
	if err != nil {
		// Panic due to compatibility reasons.
		panic(err)
	}
	return uuid
}

// NewV1 generates a new UUID based on version 1 (MAC address and
// date-time).
func NewV1() (UUID, error) {
	uuid := UUID{}
	epoch := int64(0x01b21dd213814000)
	now := uint64(time.Now().UnixNano()/100 + epoch)

	clockSeqRand := [2]byte{}
	if _, err := rand.Read(clockSeqRand[:]); err != nil {
		return uuid, err
	}
	clockSeq := binary.LittleEndian.Uint16(clockSeqRand[:])

	timeLow := uint32(now & (0x100000000 - 1))
	timeMid := uint16((now >> 32) & 0xffff)
	timeHighVer := uint16((now >> 48) & 0x0fff)
	clockSeq &= 0x3fff

	binary.LittleEndian.PutUint32(uuid[0:4], timeLow)
	binary.LittleEndian.PutUint16(uuid[4:6], timeMid)
	binary.LittleEndian.PutUint16(uuid[6:8], timeHighVer)
	binary.LittleEndian.PutUint16(uuid[8:10], clockSeq)
	copy(uuid[10:16], cachedMACAddress)

	uuid.setVersion(V1)
	uuid.setVariant(VariantRFC4122)
	return uuid, nil
}

// NewV3 generates a new UUID based on version 3 (MD5 hash of a namespace
// and a name).
func NewV3(ns UUID, name []byte) (UUID, error) {
	uuid := UUID{}
	hash := md5.New()
	if _, err := hash.Write(ns.dump()); err != nil {
		return uuid, err
	}
	if _, err := hash.Write(name); err != nil {
		return uuid, err
	}
	copy(uuid[:], hash.Sum([]byte{})[:16])

	uuid.setVersion(V3)
	uuid.setVariant(VariantRFC4122)
	return uuid, nil
}

// NewV4 generates a new UUID based on version 4 (strong random number).
func NewV4() (UUID, error) {
	uuid := UUID{}
	_, err := rand.Read([]byte(uuid[:]))
	if err != nil {
		return uuid, err
	}

	uuid.setVersion(V4)
	uuid.setVariant(VariantRFC4122)
	return uuid, nil
}

// NewV5 generates a new UUID based on version 5 (SHA1 hash of a namespace
// and a name).
func NewV5(ns UUID, name []byte) (UUID, error) {
	uuid := UUID{}
	hash := sha1.New()
	if _, err := hash.Write(ns.dump()); err != nil {
		return uuid, err
	}
	if _, err := hash.Write(name); err != nil {
		return uuid, err
	}
	copy(uuid[:], hash.Sum([]byte{})[:16])

	uuid.setVersion(V5)
	uuid.setVariant(VariantRFC4122)
	return uuid, nil
}

// Parse creates a UUID based on the given hex string which has to have
// one of the following formats:
//
// - xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// - urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// - {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
// - xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//
// The net data always has to have the length of 32 bytes.
func Parse(source string) (UUID, error) {
	var uuid UUID
	var hexSource string
	var err error
	switch len(source) {
	case 36:
		hexSource, err = parseSource(source, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	case 36 + 9:
		hexSource, err = parseSource(source, "urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	case 36 + 2:
		hexSource, err = parseSource(source, "{xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}")
	case 32:
		hexSource, err = parseSource(source, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	default:
		return uuid, fmt.Errorf("invalid source format: %q", source)
	}
	if err != nil {
		return uuid, err
	}
	hexData, err := hex.DecodeString(hexSource)
	if err != nil {
		return uuid, fmt.Errorf("source is no hex value: %w", err)
	}
	copy(uuid[:], hexData)
	// TODO: Validate UUID (version, variant).
	return uuid, nil
}

// Version returns the version number of the UUID algorithm.
func (uuid UUID) Version() Version {
	return Version(uuid[6] & 0xf0 >> 4)
}

// Variant returns the variant of the UUID.
func (uuid UUID) Variant() Variant {
	return Variant(uuid[8] & 0xe0 >> 5)
}

// Copy returns a copy of the UUID.
func (uuid UUID) Copy() UUID {
	uuidCopy := uuid
	return uuidCopy
}

// Raw returns a copy of the UUID bytes.
func (uuid UUID) Raw() [16]byte {
	uuidCopy := uuid.Copy()
	return [16]byte(uuidCopy)
}

// dump creates a copy as a byte slice.
func (uuid UUID) dump() []byte {
	dump := make([]byte, len(uuid))

	copy(dump, uuid[:])

	return dump
}

// ShortString returns a hexadecimal string representation
// without separators.
func (uuid UUID) ShortString() string {
	return fmt.Sprintf("%x", uuid[0:16])
}

// String returns a hexadecimal string representation with
// standardized separators.
func (uuid UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// setVersion sets the version part of the UUID.
func (uuid *UUID) setVersion(v Version) {
	uuid[6] = (uuid[6] & 0x0f) | (byte(v) << 4)
}

// setVariant sets the variant part of the UUID.
func (uuid *UUID) setVariant(v Variant) {
	uuid[8] = (uuid[8] & 0x1f) | (byte(v) << 5)
}

// NamespaceDNS returns the DNS namespace UUID for a v3 or a v5.
func NamespaceDNS() UUID {
	uuid, _ := Parse("6ba7b8109dad11d180b400c04fd430c8")
	return uuid
}

// NamespaceURL returns the URL namespace UUID for a v3 or a v5.
func NamespaceURL() UUID {
	uuid, _ := Parse("6ba7b8119dad11d180b400c04fd430c8")
	return uuid
}

// NamespaceOID returns the OID namespace UUID for a v3 or a v5.
func NamespaceOID() UUID {
	uuid, _ := Parse("6ba7b8129dad11d180b400c04fd430c8")
	return uuid
}

// NamespaceX500 returns the X.500 namespace UUID for a v3 or a v5.
func NamespaceX500() UUID {
	uuid, _ := Parse("6ba7b8149dad11d180b400c04fd430c8")
	return uuid
}

//--------------------
// PRIVATE HELPERS
//--------------------

// parseSource parses a source based on the given pattern. Only the
// char x of the pattern is interpreted as hex char. If the result is
// longer than 32 bytes it's an error.
func parseSource(source, pattern string) (string, error) {
	lower := []byte(strings.ToLower(source))
	raw := make([]byte, 32)
	rawPos := 0
	patternPos := 0
	patternLen := len(pattern)
	for i, b := range lower {
		if patternPos == patternLen {
			return "", fmt.Errorf("source %q too long for pattern %q", source, pattern)
		}
		switch pattern[patternPos] {
		case 'x':
			if (b < '0' || b > '9') && (b < 'a' || b > 'f') {
				return "", fmt.Errorf("source char %d is no hex char: %c", i, b)
			}
			raw[rawPos] = b
			rawPos++
			patternPos++
		default:
			if b != pattern[patternPos] {
				return "", fmt.Errorf("source char %d does not match pattern: %x is not %c", i, b, pattern[patternPos])
			}
			patternPos++
		}
	}
	return string(raw), nil
}

// macAddress retrieves the MAC address of the computer.
func macAddress() []byte {
	address := [6]byte{}
	ifaces, err := net.Interfaces()
	// Try to get address from interfaces.
	if err == nil {
		set := false
		for _, iface := range ifaces {
			if len(iface.HardwareAddr.String()) != 0 {
				copy(address[:], []byte(iface.HardwareAddr))
				set = true
				break
			}
		}
		if set {
			// Had success.
			return address[:]
		}
	}
	// Need a random address.
	if _, err := rand.Read(address[:]); err != nil {
		panic(err)
	}
	address[0] |= 0x01
	return address[:]
}

var cachedMACAddress []byte

func init() {
	cachedMACAddress = macAddress()
}

// EOF
