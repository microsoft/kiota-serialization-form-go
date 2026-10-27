# Changelog

All notable changes to this project will be documented in this file.

## [1.1.6](https://github.com/microsoft/kiota-serialization-form-go/compare/v1.1.5...v1.1.6) (2026-10-27)


### Bug Fixes

* **form:** add nil receiver safety guards to FormSerializationWriter ([647f7ba](https://github.com/microsoft/kiota-serialization-form-go/commit/647f7bacf062af29b006974fd12e09257cf9485f))
* **form:** support content-type parameters and case-insensitive matching in parse and writer factories ([681c90c](https://github.com/microsoft/kiota-serialization-form-go/commit/681c90c54daf3e542f10542d942b62ddb3d34644))
* **form:** support content-type parameters, trim whitespace, and add nil receiver guards ([7126fad](https://github.com/microsoft/kiota-serialization-form-go/commit/7126fada978903e069d93b059a605455be50ce84))
* **form:** support RFC3339Nano time parsing and add nil guards to FormParseNode ([d646004](https://github.com/microsoft/kiota-serialization-form-go/commit/d646004555d3cb3eddce809555f3760a0e90d104))
* **form:** support RFC3339Nano time parsing and add nil guards to FormParseNode ([3b0b1f8](https://github.com/microsoft/kiota-serialization-form-go/commit/3b0b1f85d53f819590cf8ade437224f0b47b6d95))

## [1.1.5](https://github.com/microsoft/kiota-serialization-form-go/compare/v1.1.4...v1.1.5) (2026-10-21)


### Bug Fixes

* **form:** allow scientific notation fallback on ParseInt range errors ([424c789](https://github.com/microsoft/kiota-serialization-form-go/commit/424c789ad544946fd8fa5d6ee401a08dabc4fd07))
* **form:** improve int32 parsing and support int8 primitive collections ([60c8680](https://github.com/microsoft/kiota-serialization-form-go/commit/60c8680fe3bff66e05ba4c75b1b9f3f8e19d2e55))
* **form:** improve int32 parsing and support int8 primitive collections ([464ea91](https://github.com/microsoft/kiota-serialization-form-go/commit/464ea91ea0e9461eaeb7679c3fa48297f2315aec))
* **form:** validate float bounds before integer conversion and sequence float check ([75b6b0b](https://github.com/microsoft/kiota-serialization-form-go/commit/75b6b0b8b3cc99bad6196ecc718b454d05146d60))
* **form:** validate full precision exponent tokens using big.Rat to prevent float64 rounding overflow ([66eff78](https://github.com/microsoft/kiota-serialization-form-go/commit/66eff78ac8132d50ab2ccf4f28e099cc4c14147c))

## [1.1.4](https://github.com/microsoft/kiota-serialization-form-go/compare/v1.1.3...v1.1.4) (2026-10-14)


### Bug Fixes

* preserve decoded form field values and embedded equals ([3b51599](https://github.com/microsoft/kiota-serialization-form-go/commit/3b51599fe1715386779c73eab79b1b0b88baca44))
* preserve form field values during deserialization ([14e3f79](https://github.com/microsoft/kiota-serialization-form-go/commit/14e3f79f77c54b23ffe015c84eb9b9d6d8451503))
* preserve int64 precision and unsigned byte values in form parsing ([51c2cf7](https://github.com/microsoft/kiota-serialization-form-go/commit/51c2cf794776f338335c9a6641c88626d490fd39))
* preserve int64 precision and unsigned byte values in form parsing ([fca913c](https://github.com/microsoft/kiota-serialization-form-go/commit/fca913c8b2b4218e5ca7eab9b3aac9642953fe0b))
* preserve range errors for form int64 values ([b353fb9](https://github.com/microsoft/kiota-serialization-form-go/commit/b353fb93332ac59f2821446e5815d94a5f3edafe))

## [1.1.3](https://github.com/microsoft/kiota-serialization-form-go/compare/v1.1.2...v1.1.3) (2026-03-13)


### Bug Fixes

* reduce parse node allocations when deserializing primitive types collections ([d871ef8](https://github.com/microsoft/kiota-serialization-form-go/commit/d871ef8be0272296fb53389725ca01516d188599))


### Performance Improvements

* extract private helpers and avoid parse node allocations in primitive collection deserialization ([ba61e0b](https://github.com/microsoft/kiota-serialization-form-go/commit/ba61e0bb7331cf8a563317b6d9d79a95582c2dfc))
* reduce parse node allocations when deserializing primitive types ([61219f6](https://github.com/microsoft/kiota-serialization-form-go/commit/61219f6d0054f12523acc89e87bfff7b07a5e36e))

## [1.1.2](https://github.com/microsoft/kiota-serialization-form-go/compare/v1.1.1...v1.1.2) (2025-04-02)


### Bug Fixes

* removes common go dependency ([ce7bc35](https://github.com/microsoft/kiota-serialization-form-go/commit/ce7bc35ef5af4073244de8bdab87332dd4897864))
* removes common go dependency ([1265450](https://github.com/microsoft/kiota-serialization-form-go/commit/1265450a46777d7fd22500f782e587f9eae85217))

## [1.1.1](https://github.com/microsoft/kiota-serialization-form-go/compare/v1.1.0...v1.1.1) (2025-03-24)


### Bug Fixes

* upgrades common go dependency to solve triming issues ([d3e8da1](https://github.com/microsoft/kiota-serialization-form-go/commit/d3e8da1ce7cae10758a8874373cf656a8bee7645))
* upgrades common go dependency to solve triming issues ([6dda99a](https://github.com/microsoft/kiota-serialization-form-go/commit/6dda99ae94e501484e486ea907f682fb2ee7100f))

## [1.1.0](https://github.com/microsoft/kiota-serialization-form-go/compare/v1.0.0...v1.1.0) (2025-03-13)


### Features

* upgrades required go version from go1.18 to go 1.22 ([d92bfde](https://github.com/microsoft/kiota-serialization-form-go/commit/d92bfdea60153a0cc1f7e38b23db33a313a6fd46))

## [1.0.0] - 2023-05-04

### Changed

- GA Release.

## [0.9.0] - 2023-03-01

### Changed

- Update version number to 0.9.0. Fixes error in release number

## [0.4.0] - 2023-03-01

### Added

- Adds support for serialization of collections by serialization writer and parsenode.

## [0.3.0] - 2023-01-26

### Changed

- Added support for backing store.

## [0.2.0] - 2022-12-16

### Changed

- Fixed key encoding

## [0.1.0] - 2022-12-15

### Added

- Initial release of the package.
