# Changelog

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
