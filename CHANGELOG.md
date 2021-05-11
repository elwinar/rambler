# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## Unreleased

## [5.4.0](https://github.com/elwinar/rambler/releases/tag/5.4.0) - 2021-05-11
### Changed
- Open files when scanning a migration to avoid unnecessary disk operations

## [5.3.0](https://github.com/elwinar/rambler/releases/tag/5.3.0) - 2020-11-15
### Added
- Better support for PostgreSQL's sslmode

## [5.2.0](https://github.com/elwinar/rambler/releases/tag/5.2.0) - 2019-08-28
### Added
- `no-save` flag
- `migration` flag

## [5.1.0](https://github.com/elwinar/rambler/releases/tag/5.1.0) - 2019-08-28
### Added
- `dry-run` flag

## [5.0.0](https://github.com/elwinar/rambler/releases/tag/5.0.0) - 2019-05-26
### Added
- PostgreSQL support for schemas
### Changed
- It is no longer an error if no configuration file is provided as option and
  the default one doesn't exists
### Internal
- Split the general configuration and the driver configuration to allow a
  better structure

## [4.2.1](https://github.com/elwinar/rambler/releases/tag/4.2.1) - 2018-08-09
### Fixed
- Return status of the commands in case of error
- Static linking of the binary built for docker

## [4.2.0](https://github.com/elwinar/rambler/releases/tag/4.2.0) - 2018-06-05
### Added
- Multistage docker build
- DockerHub automated build

## [4.1.4](https://github.com/elwinar/rambler/releases/tag/4.1.4) - 2018-06-05
### Fixed
- PostgreSQL DSN order

## [4.1.3](https://github.com/elwinar/rambler/releases/tag/4.1.3) - 2017-02-26
### Fixed
- Respect the debug logging flag

## [4.1.2](https://github.com/elwinar/rambler/releases/tag/4.1.2) - 2017-02-26
### Fixed
- Missing formatter on environment log

## [4.1.1](https://github.com/elwinar/rambler/releases/tag/4.1.1) - 2017-02-26
### Fixed
- Logging formatting
- Missing default configuration values

## [4.1.0](https://github.com/elwinar/rambler/releases/tag/4.1.0) - 2017-01-02
### Added
- Docker image
- Logs various actions

### Internal
- Moved from github.com/codegangsta/cli to [github.com/urfave/cli](https://github.com/urfave/cli)

## [4.0.0](https://github.com/elwinar/rambler/releases/tag/4.0.0) - 2016-10-30
### Added
- MySQL driver now use the migration name as primary key (necessary for Percona
  XtraDB clusters)

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
- Possibility of using alternatives configuration file types (namely YAML and
  TOML)

### Fixed
- Behavior of the transactions
