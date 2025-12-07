# Changelog

## v0.3.2 (2025-12-07)

### Fixed
* **CRITICAL**: Fixed UUIDv7 sortability by implementing monotonic counter
* UUIDv7 now uses 12-bit sequence counter for same-millisecond UUIDs
* Ensures strict lexicographic ordering even when generating thousands of UUIDs per millisecond
* Handles clock rollback gracefully
* Thread-safe implementation with mutex protection

### Changed
* Enhanced UUIDv7 test coverage with 10,000 UUID generation test
* Improved TestV7Sortability to verify strict monotonic ordering
* Updated TestV7Monotonicity to test 10,000 UUIDs instead of 1,000
* Added detailed logging to sortability tests

### Technical Details
* Implements RFC 9562 Section 6.2 recommendations for monotonic counters
* 48-bit millisecond timestamp + 12-bit sequence + 62-bit random data
* Sequence counter randomly initialized per millisecond
* Automatic sequence increment for same-millisecond UUIDs
* Handles sequence overflow by waiting for next millisecond

## v0.3.1 (2025-12-07)

### Added
* Comprehensive Makefile for automated build process
* Makefile targets: all, help, tidy, lint, build, test, bench, fuzz, coverage, clean, install-tools, ci
* Color-coded output for better readability
* Automatic dependency handling between targets
* Coverage report generation (HTML format)
* Development tools installation script
* Detailed Makefile usage documentation in README.md

### Changed
* Enhanced README.md with complete Makefile documentation
* Added development workflow section
* Improved contributor onboarding with make targets

## v0.3.0 (2025-12-06)

### Added
* UUID version 2 (DCE Security) implementation with POSIX UID/GID support
* `NewV2()` function for custom domain and ID
* `NewV2Person()` convenience function using current process UID
* `NewV2Group()` convenience function using current process GID
* `Domain` type with Person, Group, and Org constants
* `Domain()` method to extract domain from UUID v2
* `ID()` method to extract local identifier from UUID v2
* `Domain.String()` method for human-readable domain names
* Comprehensive UUID v2 test suite with 95%+ coverage
* UUID v2 benchmarks (NewV2, NewV2Person, NewV2Group)
* Fuzz tests for UUID v2 domain values and parsing
* Enhanced documentation with UUID v2 usage examples

### Changed
* Updated to Go 1.25
* Enhanced golangci-lint configuration
* Improved test coverage across all UUID versions
* Updated documentation to include DCE Security UUID examples

## v0.2.0 (2025-12-04)

### Added
* UUID version 6 (reordered Gregorian time-based) implementation per RFC 9562
* UUID version 7 (Unix Epoch time-based) implementation per RFC 9562
* Comprehensive test suite with sortability, monotonicity, and concurrency tests
* Benchmark tests for all UUID versions
* Enhanced documentation with usage examples

### Changed
* Updated to Go 1.24
* Migrated from `tideland.dev/go/audit/asserts` to `tideland.dev/go/asserts/verify`
* Modernized test patterns following current Tideland Go standards
* Updated copyright to 2021-2025
* Improved documentation in README.md with comprehensive examples
* Enhanced doc.go to reflect RFC 9562 compliance

### Fixed
* Package import path updated to modern standards
* Test coverage improvements
* Code style alignment with other Tideland Go packages

## v0.1.1 (2021-08-31)

* Renamed FromHex() to Parse() and improved it
* Fix GitHub workflow

## v0.1.0 (2021-08-28)

* Migrated UUID code from former DSA Identifier package
