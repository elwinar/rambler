# rambler 
[![build](https://app.wercker.com/status/b645428b6f548288d71d3ba83cc1a783/s/master "wercker status")](https://app.wercker.com/project/bykey/b645428b6f548288d71d3ba83cc1a783)
[![Coverage Status](https://coveralls.io/repos/elwinar/rambler/badge.svg?branch=master&service=github)](https://coveralls.io/github/elwinar/rambler?branch=master)

A simple and language-independent SQL schema migration tool

## Installation

You can download the latest release on the [release
page](https://github.com/elwinar/rambler/releases) of the project.

Go users can also simply compile it from source and install it as a go
executable using the following command :

```
go install github.com/elwinar/rambler
```

Releases are compiled using the wonderful
[XGo](https://github.com/karalabe/xgo). Don't hesitate to check it out, it
really kicks some serious ass.

## Usage

### Migrations

In rambler, migrations are kept in the simplest form possible: a migration is a
list of sections (`up` and `down`), each section being an SQL statement.
Example:

```sql
-- rambler up

CREATE TABLE foo (
	id INTEGER UNSIGNED AUTO_INCREMENT,
	bar VARCHAR(60),
	PRIMARY KEY (id)
);

-- rambler down

DROP TABLE foo;
```

Sections are delimited by SQL comments suffixed by the rambler marker
(white-spaces sensitive). While applying a migration, rambler will execute each
`up` section in order, and while reversing it it will execute each `down`
section in reverse order.

Migrations are executed in alphabetical order, thus a versioning scheme of the
form `version_description.sql` is highly recommended, version being an integer
value, and description an underscored string. Examples:

* `201409272258_Added_table_foo.sql`
* `01_First_migration.sql`

The migrations applied to the database are stored in a table named `migration`
(can be changed with the `table` configuration option).

### Configuration

Rambler configuration is lightweight: just dump the credentials of your
database and the path to your migrations' directory into a JSON file, and
you're done. Here is an example or JSON configuration file with the default
values of rambler:

```json
{
	"driver": "mysql",
	"protocol": "tcp",
	"host": "localhost",
	"port": 3306,
	"user": "root",
	"password": "",
	"database": "",
	"schema": "",
	"directory": ".",
	"table": "migrations"
}
```

When running, rambler will try to find a configuration file in the working
directory and use its values to connect to the managed database.

#### HJSON

Rambler now supports [HJSON](http://hjson.org/) configuration files, which is
by the way retrocompatible with JSON.

#### Environment Variables

Alternatively, Rambler can read configuration from environment variables. The
environment variables can override any of the confifuration file values and
are prefixed with `RAMBLER_`.

| Env Var           | Config    |
|-------------------|-----------|
| RAMBLER_DRIVER    | driver    |
| RAMBLER_PROTOCOL  | protocol  |
| RAMBLER_HOST      | host      |
| RAMBLER_PORT      | port      |
| RAMBLER_USER      | user      |
| RAMBLER_PASSWORD  | password  |
| RAMBLER_DATABASE  | database  |
| RAMBLER_SCHEMA    | schema    |
| RAMBLER_DIRECTORY | directory |
| RAMBLER_TABLE     | table     |

#### Drivers

Rambler supports actually 3 drivers:

- `mysql`
- `postgresql`
- `sqlite`

Don't hesitate to get in touch if you want to see another one supported,
provided a golang `database/sql` driver exist for your database vendor.

*Note that some configuration options only apply to some drivers.*

### Applying a migration

To apply a migration, use the `apply` command.

```
rambler apply
```

Rambler will compare the migrations already applied and the available
migrations in increasing order to find the next migration to apply, then
execute all its `up` sections' statements in order. 

### Reversing a migration

To reverse a migration, use the `reverse` command.

```
rambler reverse
```

Rambler will compare the migrations already applied and the available
migrations in decreasing order to find the last applied migrations, then
execute all its `down` sections' statements in reverse order.

### Options

- `--all, -a` repeat the command until there is no more migration to
  apply/reverse. This flag is exclusive with `--migration`.
- `--no-save` doesn't save the applied/reversed migration. This option is
  mainly destined to allow faster iteration when writing the migration.
- `--migration` bypass the migration automatic discovery and apply the migration
  whose path is given. This is exclusive with `--all`.

### Errors

To ensure database schema consistency, rambler will complain and stop when
encountering a new migration in the middle of the already existing ones or if
it can't find a migration already applied.

### Environments

An environment is an additional configuration that is given a name, and can be
used to create multiple configurations for a single application (for example,
to differentiate production, testing, etc).

Environments are defined in the configuration file, under the `environments`
item.  Each environment is defined as an attribute of this item, the key being
the name and the value being the configuration options.

Environments configuration are derived from the default configuration of
rambler (at the configuration file's root), so you only need to override the
needed options:

```json
{
    "driver": "mysql",
    "protocol": "tcp",
    "port": 3306,
    "user": "root",
    "password": "",
    "database": "rambler_default",
    "directory": "migrations",
    "table": "migrations",
    "environments": {
        "development": {
            "database": "rambler_development"
        },
        "testing": {
            "database": "rambler_testing"
        }
    }
}
```

Here we have three environments defined:
- `default`, will use the `rambler_default` database,
- `development`, will use the `rambler_development` database,
- `testing`, will use the `rambler_testing` database;

### Logging

Rambler will log a few important informations for monitoring what is happening
on stdout. If you suspect something of being wrong, you can also use the debug
mode by adding `--debug` to your command line.

### Dry-run

The `--dry-run` flag will print statements instead of executing them.

## CONTRIBUTORS

- [cjhubert](https://github.com/cjhubert)
- [shawndellysse](https://github.com/shawndellysse)

## Feedback and contributions

Feel free to give feedback, make pull requests or simply open issues if you
find a bug or have an idea.
