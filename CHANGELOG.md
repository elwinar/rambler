# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## Unreleased

## [3.4.0](https://github.com/elwinar/rambler/releases/tag/v3.4.0) - 2016-10-27

### Added

- HJSON support for the configuration file

## [3.3.0](https://github.com/elwinar/rambler/releases/tag/v3.3.0) - 2016-07-13

### Added

- SQLite driver
- Exit with error code 1 in case of error
- Dynamic version into the build
- A simple makefile for convenience

### Removed

- GoXC configuration file

## [3.2.0](https://github.com/elwinar/rambler/releases/tag/v3.2.0) - 2015-09-19

### Added

- Add GoXC configuration file to automate the release process

## [3.1.0](https://github.com/elwinar/rambler/releases/tag/v3.1.0) - 2015-09-18

### Added

- Continuous integration tooling with Wercker
- [Unlicense](http://unlicense.org/) LICENSE file
- Ego-growing badges on the top of the README file

### Fixed

- Various linting warnings

## [3.0.0](https://github.com/elwinar/rambler/releases/tag/v3.0.0) - 2015-08-21

### Changed

- Complete rewrite of the software for simplicity

### Fixed

- Configuration bug due to vendor breaking change

## [2.1.0](https://github.com/elwinar/rambler/releases/tag/v2.1.0) - 2014-12-24

### Added 

- Postgresql driver from the good work of [cjhubert](https://github.com/cjhubert)

## [2.0.0](https://github.com/elwinar/rambler/releases/tag/v2.0.0) - 2014-12-24

### Changed

- Complete rewrite to add unit-testing

### Removed

- The command-line options to override the configuration
- Posibility of using alternatives configuration file types (namely YAML and TOML)

### Fixed

- Behavior of the transactions

