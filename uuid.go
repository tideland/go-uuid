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
	"os"
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
	V2 Version = 2
	V3 Version = 3
	V4 Version = 4
	V5 Version = 5
	V6 Version = 6
	V7 Version = 7
)

// Variant represents a UUID's variant.
type Variant byte

const (
	VariantNCS       Variant = 0 // Reserved, NCS backward compatibility.
	VariantRFC4122   Variant = 4 // The variant specified in RFC4122.
	VariantMicrosoft Variant = 6 // Reserved, Microsoft Corporation backward compatibility.
	VariantFuture    Variant = 7 // Reserved for future definition.
)

// Domain represents a DCE Security (Version 2) UUID domain.
type Domain byte

// Domain constants for DCE Security (Version 2) UUIDs.
const (
	Person Domain = 0 // POSIX UID domain
	Group  Domain = 1 // POSIX GID domain
	Org    Domain = 2 // Organization domain
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

// NewV2 generates a new UUID based on version 2 (DCE Security).
// The domain should be one of Person, Group, or Org.
// On a POSIX system, the id should be the user's UID for the Person
// domain and the user's GID for the Group domain. The meaning of id for
// the Org domain or on non-POSIX systems is site-defined.
//
// For a given domain/id pair, the same token may be returned for up to
// 7 minutes and 10 seconds.
func NewV2(domain Domain, id uint32) (UUID, error) {
	// Start with a v1 UUID
	uuid, err := NewV1()
	if err != nil {
		return uuid, err
	}

	// Replace time_low with the local ID
	binary.BigEndian.PutUint32(uuid[0:4], id)

	// Set the domain in the clock_seq_low position
	uuid[9] = byte(domain)

	// Set version to 2
	uuid.setVersion(V2)

	return uuid, nil
}

// NewV2Person returns a DCE Security (Version 2) UUID in the Person
// domain with the id returned by os.Getuid.
func NewV2Person() (UUID, error) {
	return NewV2(Person, uint32(os.Getuid()))
}

// NewV2Group returns a DCE Security (Version 2) UUID in the Group
// domain with the id returned by os.Getgid.
func NewV2Group() (UUID, error) {
	return NewV2(Group, uint32(os.Getgid()))
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

// NewV6 generates a new UUID based on version 6 (reordered Gregorian timestamp).
// UUIDv6 is a field-compatible version of UUIDv1, reordered for improved DB locality.
// The timestamp bytes are stored from most to least significant for better sortability.
func NewV6() (UUID, error) {
	uuid := UUID{}
	epoch := int64(0x01b21dd213814000)
	now := uint64(time.Now().UnixNano()/100 + epoch)

	clockSeqRand := [2]byte{}
	if _, err := rand.Read(clockSeqRand[:]); err != nil {
		return uuid, err
	}
	clockSeq := binary.LittleEndian.Uint16(clockSeqRand[:])

	// Extract timestamp components
	timeHigh := uint32((now >> 28) & 0xffffffff) // Most significant 32 bits
	timeMid := uint16((now >> 12) & 0xffff)      // Middle 16 bits
	timeLow := uint16(now & 0x0fff)              // Least significant 12 bits
	clockSeq &= 0x3fff

	// Store in big-endian order for v6
	binary.BigEndian.PutUint32(uuid[0:4], timeHigh)
	binary.BigEndian.PutUint16(uuid[4:6], timeMid)
	binary.BigEndian.PutUint16(uuid[6:8], timeLow)
	binary.BigEndian.PutUint16(uuid[8:10], clockSeq)
	copy(uuid[10:16], cachedMACAddress)

	uuid.setVersion(V6)
	uuid.setVariant(VariantRFC4122)
	return uuid, nil
}

// NewV7 generates a new UUID based on version 7 (Unix Epoch timestamp).
// UUIDv7 features a time-ordered value field derived from Unix Epoch timestamp
// in milliseconds with improved entropy characteristics.
func NewV7() (UUID, error) {
	uuid := UUID{}

	// Get Unix timestamp in milliseconds
	now := time.Now()
	unixMs := uint64(now.UnixMilli())

	// Fill first 48 bits with timestamp
	uuid[0] = byte(unixMs >> 40)
	uuid[1] = byte(unixMs >> 32)
	uuid[2] = byte(unixMs >> 24)
	uuid[3] = byte(unixMs >> 16)
	uuid[4] = byte(unixMs >> 8)
	uuid[5] = byte(unixMs)

	// Fill remaining bits with random data
	randData := make([]byte, 10)
	if _, err := rand.Read(randData); err != nil {
		return uuid, err
	}
	copy(uuid[6:], randData)

	uuid.setVersion(V7)
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

// dump creates a copy as a byte slice.
func (uuid UUID) dump() []byte {
	dump := make([]byte, len(uuid))

	copy(dump, uuid[:])

	return dump
}

// setVersion sets the version part of the UUID.
func (uuid *UUID) setVersion(v Version) {
	uuid[6] = (uuid[6] & 0x0f) | (byte(v) << 4)
}

// setVariant sets the variant part of the UUID.
func (uuid *UUID) setVariant(v Variant) {
	uuid[8] = (uuid[8] & 0x1f) | (byte(v) << 5)
}

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
