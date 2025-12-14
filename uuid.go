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
	"sync"
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
	uuid.setVariant()
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
	uuid.setVariant()
	return uuid, nil
}

// NewV4 generates a new UUID based on version 4 (strong random number).
func NewV4() (UUID, error) {
	uuidRand := make([]byte, 16)
	_, err := rand.Read(uuidRand)
	if err != nil {
		return UUID(uuidRand), err
	}

	uuid := UUID(uuidRand)
	uuid.setVersion(V4)
	uuid.setVariant()
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
	uuid.setVariant()
	return uuid, nil
}

// NewV6 generates a new UUID based on version 6 (reordered Gregorian timestamp).
// UUIDv6 is a field-compatible version of UUIDv1, reordered for improved DB locality.
// The timestamp bytes are stored from most to least significant for better sortability.
//
// This implementation ensures monotonicity by using a counter for the clock sequence
// when UUIDs are generated within the same timestamp period.
func NewV6() (UUID, error) {
	uuid := UUID{}
	epoch := int64(0x01b21dd213814000)
	now := uint64(time.Now().UnixNano()/100 + epoch)

	// Get monotonic clock sequence and adjusted timestamp
	adjustedNow, clockSeq, err := getV6ClockSeq(now)
	if err != nil {
		return uuid, err
	}

	// Extract timestamp components
	timeHigh := uint32((adjustedNow >> 28) & 0xffffffff) // Most significant 32 bits
	timeMid := uint16((adjustedNow >> 12) & 0xffff)      // Middle 16 bits
	timeLow := uint16(adjustedNow & 0x0fff)              // Least significant 12 bits

	// Store in big-endian order for v6
	binary.BigEndian.PutUint32(uuid[0:4], timeHigh)
	binary.BigEndian.PutUint16(uuid[4:6], timeMid)
	binary.BigEndian.PutUint16(uuid[6:8], timeLow)
	binary.BigEndian.PutUint16(uuid[8:10], clockSeq)
	copy(uuid[10:16], cachedMACAddress)

	uuid.setVersion(V6)
	uuid.setVariant()
	return uuid, nil
}

// NewV7 generates a new UUID based on version 7 (Unix Epoch timestamp).
// UUIDv7 features a time-ordered value field derived from Unix Epoch timestamp
// in milliseconds with improved entropy characteristics.
//
// This implementation ensures monotonicity by using a counter for UUIDs
// generated within the same millisecond, as recommended by RFC 9562 Section 6.2.
func NewV7() (UUID, error) {
	uuid := UUID{}

	// Get time and sequence with monotonicity guarantee
	ms, seq, err := getV7Time()
	if err != nil {
		return uuid, err
	}

	// Fill first 48 bits with timestamp (milliseconds)
	uuid[0] = byte(ms >> 40)
	uuid[1] = byte(ms >> 32)
	uuid[2] = byte(ms >> 24)
	uuid[3] = byte(ms >> 16)
	uuid[4] = byte(ms >> 8)
	uuid[5] = byte(ms)

	// Fill next 12 bits (after version) with sequence counter
	// This ensures monotonicity within the same millisecond
	uuid[6] = byte(seq >> 8)
	uuid[7] = byte(seq)

	// Fill remaining 62 bits (after variant) with random data
	randData := make([]byte, 8)
	if _, err := rand.Read(randData); err != nil {
		return uuid, err
	}
	uuid[8] = randData[0]
	uuid[9] = randData[1]
	uuid[10] = randData[2]
	uuid[11] = randData[3]
	uuid[12] = randData[4]
	uuid[13] = randData[5]
	uuid[14] = randData[6]
	uuid[15] = randData[7]

	uuid.setVersion(V7)
	uuid.setVariant()
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

// Domain returns the domain for a Version 2 UUID.
// Domains are only defined for Version 2 UUIDs.
func (uuid UUID) Domain() Domain {
	return Domain(uuid[9])
}

// ID returns the local identifier for a Version 2 UUID.
// IDs are only defined for Version 2 UUIDs.
func (uuid UUID) ID() uint32 {
	return binary.BigEndian.Uint32(uuid[0:4])
}

// String returns the string representation of a Domain.
func (d Domain) String() string {
	switch d {
	case Person:
		return "Person"
	case Group:
		return "Group"
	case Org:
		return "Org"
	}
	return fmt.Sprintf("Domain%d", int(d))
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

// setVariant sets the variant part of the UUID, always RFC4122.
// Used to keep source more consistent with version and variant.
func (uuid *UUID) setVariant() {
	uuid[8] = (uuid[8] & 0x1f) | (byte(VariantRFC4122) << 5)
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

// v6State holds the state for monotonic UUID v6 generation.
type v6State struct {
	mu           sync.Mutex
	lastTime     uint64 // Last timestamp (in 100ns intervals)
	lastClockSeq uint16 // Last clock sequence
}

// v7State holds the state for monotonic UUID v7 generation.
type v7State struct {
	mu      sync.Mutex
	lastMs  int64  // Last millisecond timestamp
	lastSeq uint16 // Last sequence number
}

var v6Generator = &v6State{}
var v7Generator = &v7State{}

// getV7Time returns the current time in milliseconds and a monotonic sequence number.
// The returned values ensure that each UUID v7 is greater than the previous one,
// even when multiple UUIDs are generated within the same millisecond.
func getV7Time() (ms int64, seq uint16, err error) {
	v7Generator.mu.Lock()
	defer v7Generator.mu.Unlock()

	// Get current time in milliseconds - capture once to avoid inconsistencies
	now := time.Now().UnixMilli()

	switch {
	case now == v7Generator.lastMs:
		// Same millisecond: increment sequence
		if v7Generator.lastSeq == 0x0FFF {
			// Sequence overflow - this is extremely rare but we should handle it
			// Wait for the next millisecond
			for {
				time.Sleep(time.Microsecond * 100)
				newNow := time.Now().UnixMilli()
				if newNow > now {
					now = newNow
					break
				}
			}
			// Start with sequence 0 for new millisecond after overflow
			v7Generator.lastSeq = 0
			v7Generator.lastMs = now
		} else {
			v7Generator.lastSeq++
		}
	case now > v7Generator.lastMs:
		// New millisecond: initialize with random sequence
		randBytes := make([]byte, 2)
		if _, err := rand.Read(randBytes); err != nil {
			return 0, 0, err
		}
		v7Generator.lastSeq = binary.BigEndian.Uint16(randBytes) & 0x0FFF
		v7Generator.lastMs = now
	default:
		// Clock went backwards - this is problematic
		// Use the last known time and increment sequence
		if v7Generator.lastSeq == 0x0FFF {
			v7Generator.lastSeq = 0
			v7Generator.lastMs++
		} else {
			v7Generator.lastSeq++
		}
		now = v7Generator.lastMs
	}

	return now, v7Generator.lastSeq, nil
}

// getV6ClockSeq returns a monotonic clock sequence for UUID v6 generation.
// The returned clock sequence ensures that each UUID v6 is greater than the previous one,
// even when multiple UUIDs are generated within the same timestamp period.
func getV6ClockSeq(timestamp uint64) (adjustedTimestamp uint64, clockSeq uint16, err error) {
	v6Generator.mu.Lock()
	defer v6Generator.mu.Unlock()

	switch {
	case timestamp == v6Generator.lastTime:
		// Same timestamp: increment clock sequence
		if v6Generator.lastClockSeq == 0x3FFF {
			// Clock sequence overflow - this is extremely rare
			// Wait for the next timestamp unit (100ns)
			for timestamp == v6Generator.lastTime {
				time.Sleep(time.Nanosecond * 100)
				epoch := int64(0x01b21dd213814000)
				timestamp = uint64(time.Now().UnixNano()/100 + epoch)
			}
			// Initialize new clock sequence with random value
			randBytes := make([]byte, 2)
			if _, err := rand.Read(randBytes); err != nil {
				return 0, 0, err
			}
			v6Generator.lastClockSeq = binary.BigEndian.Uint16(randBytes) & 0x3FFF
		} else {
			v6Generator.lastClockSeq++
		}
	case timestamp > v6Generator.lastTime:
		// New timestamp: initialize with random clock sequence
		randBytes := make([]byte, 2)
		if _, err := rand.Read(randBytes); err != nil {
			return 0, 0, err
		}
		v6Generator.lastClockSeq = binary.BigEndian.Uint16(randBytes) & 0x3FFF
		v6Generator.lastTime = timestamp
	default:
		// Clock went backwards - this is problematic
		// Use the last known time and increment clock sequence
		if v6Generator.lastClockSeq == 0x3FFF {
			v6Generator.lastClockSeq = 0
			v6Generator.lastTime++
		} else {
			v6Generator.lastClockSeq++
		}
		timestamp = v6Generator.lastTime
	}

	return timestamp, v6Generator.lastClockSeq, nil
}

// EOF
