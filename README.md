# MDB - Metadata DB

[![Maintainability](https://api.codeclimate.com/v1/badges/3f7da320a15bb28e8af5/maintainability)](https://codeclimate.com/github/Bnei-Baruch/mdb/maintainability)

[![Test Coverage](https://api.codeclimate.com/v1/badges/3f7da320a15bb28e8af5/test_coverage)](https://codeclimate.com/github/Bnei-Baruch/mdb/test_coverage)

## Overview

BB archive Metadata Database.

This system aims to be a single source of truth for all content produced by Bnei Baruch. 


## Developer Environment

We assume docker and golang are already installed on your system.

### Tools

```Shell
go get -u github.com/jteeuwen/go-bindata/...
```

See **_rambler_** under Schema Migrations and **_sqlboiler_** ORM sections below.

### Useful Commands

```Shell
mdb config <path>
```

Generate default configuration in the given path. If path is omitted STDOUT is used instead.
**Note** that default value to config file is `config.toml` in project root directory.


```Shell
rambler -c migrations/rambler.json apply -a
```
Apply all DB migrations (See Schema migrations section for more information)

```Shell
 mdb migration my-migration-name
```
Create new migration



### Schema Migrations
We keep track of all changes to the MDB schema under `migrations`.
These are pure postgres sql scripts.
To create a new migration file with name <my-migration-name> run in project root directory:
```Shell
mdb migration my-migration-name
```
This will create a migration file in migrations directory with name like: `2017-01-07_14:21:02_my-migration-name.sql`

They play along well with [rambler](https://github.com/elwinar/rambler) A simple and language-independent SQL schema migration tool.
Download the rambler executable for your system from the [release page](https://github.com/elwinar/rambler/releases).
(on linux `chmod +x`)

Under `migrations` folder add a `rambler.json` config file. 
Simply copying `rambler.sample.json` should work fine but feel free to change that.
Check out the docs for configuration options.

**Important** make sure never to commit such files to SCM.

On the command line:

```Shell
rambler -c migrations/rambler.json apply -a
```


### Generating Documentation

Documentation is based on tests and will be generated automatically with each `make build`. To generate static html documentation (`docs.html`) install:

```Shell
npm install -g aglio
```

Then run:

```Shell
make api
```


## Implementation Notes

### Dates and Times
All timestamps fields are expecting values in UTC only.


### Languages
Languages are represented in the system as a two letters code adhering to the [ISO_639-1](https://en.wikipedia.org/wiki/ISO_639-1) standard.

Special values:

* Unknown - `xx` 
* Multiple languages - `zz` 


### Logging
We use [logrus](https://github.com/Sirupsen/logrus) for logging.

All logs are written to STDOUT or STDERR. It's up to the running environment
to pipe these into physical files and rotate those using `logrotate(8)`.


### Rollbar
If the `rollbar-token` config is found we'll use our own recovery middleware to send errors to [Rollbar](https://rollbar.com).
If not, we'll use gin.Recovery() to print stacktrace to console. Using rollbar is meant for production environment.

 In addition, you could log whatever error you want to rollbar directly, for example:

 ```Go
    if _, err := SomeErrorProneFunc(); err != nil {
        rollbar.Error("level", err,...)
    }
 ```

 Check out the [docs](https://godoc.org/github.com/stvp/rollbar) for more info on how to use the Rollbar client.


### ORM
We use [sqlboiler](https://github.com/volatiletech/sqlboiler) as an ORM.


In root folder add a `sqlboiler.toml` config file.
Simply copying `sqlboiler.sample.toml` should work fine but feel free to change that.
Check out the docs for configuration options.

To regenerate models run 
```Shell
make models
```

## License

MIT