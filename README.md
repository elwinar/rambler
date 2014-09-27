# rambler

A simple and language-independent SQL schema migration tool

## Installation

Go users can simply compile it from source and install it as a go executable using the following command :

```
go install github.com/elwinar/rambler
```

Others will have to wait for me to have time to cross-compile the executables… nothing really fancy (thank to golang), but not done yet.

## Usage

### Migrations

In rambler, migrations are kept in the simplest form possible: a migration is a list of sections (`up` and `down`, each section being *chunks* of sql. Example:

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

Sections are delimited by SQL comments prefixed by the rambler marker (currently sensitive to white-spaces). While applying a migration, rambler will execute each `up` section in order, and while reversing it it will execute each `down` section in reverse order.

Migrations filename must be of the form `version_comment.sql`, version being an integer value, and comment an underscored string. Examples:

* `201409272258_Added_table_foo.sql`
* `01_First_migration.sql`

### Configuration

Rambler configuration is lightweight: just dump the credentials of your database and the path to your migrations' directory into a JSON (or YAML or TOML) file, and you're done. Here is an example or JSON configuration file with the default values of rambler:

```json
{
	"driver": "mysql",
	"protocol": "tcp",
	"host": "localhost",
	"port": 3306,
	"user": "root",
	"password": "",
	"database": "",
	"migrations": "."
}
```

When running, rambler will try to find a configuration file in the working directory and use its values to connect to the managed database. Every option can be overriden at runtime by the matching command-line option (`rambler help` to get the shorthands).

### Applying a migration

To apply a migration, use the `apply` command.

```
rambler apply
```

Rambler will compare the migrations already applied and the available migrations in increasing version order to find the next migration to apply, then execute all its `up` sections' statements in order. 

### Reversing a migration

To reverse a migration, use the `reverse` command.

```
rambler reverse
```

Rambler will compare the migrations already applied and the available migrations in decreasing version order to find the last applied migrations, then execute all its `down` sections' statements in reverse order.

### Options

You can tell rambler to repeat the process while there is a migration to apply with the `all` flag (or its shorthand, `a`).

### Errors

To ensure database schema consistency, rambler will complain and stop when encountering a new migration in the middle of the already existing ones or if it can't find a migration already applied.

## TODO

* Add a `refresh` command that will reverse the last migration then apply it again
* Add a `number` option to choose the number of migrations to apply, reverse, or refresh
* Add a `ignore-missing` option to ignore missing migrations and continue
* Add a `ignore-out-of-order` option to ignore out-of-order migrations and apply, reverse or refresh them anyway

## Feedback and contributions

Feel free to give feedback, make pull requests or simply open issues if you find a bug or have an idea.
