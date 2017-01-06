# MDB - Metadata DB

## Overview

BB archive Metadata Database.

This system aims to be a single source of truth for all content produced by Bnei Baruch. 


## Commands
The mdb is meant to be executed as command line. 
Type `mdb <command> -h` to see how to use each command.
 
```Shell
mdb server
```

Run the server

```Shell
mdb config <path>
```

Generate default configuration in the given path. If path is omitted STDOUT is used instead.
**Note** that default value to config file is `config.toml` in project root directory.

```Shell
 mdb migration my-migration-name
```
Create new migration. (See Schema migrations section for more information).

```Shell
mdb version
```

Print the version of MDB

## Implementation Notes

### Dates and Times
All timestamps fields are expecting values in UTC only.


### Languages
Languages are represented in the system as a two letters code adhering to the [ISO_639-1](https://en.wikipedia.org/wiki/ISO_639-1) standard.

Special values:

* Unknown - `xx` 
* Multiple languages - `zz` 


## Release and Deployment

Once development is done, all tests are green, we want to go live.
All we have to do is simply execute `misc/release.sh`.

To add a pre-release tag, add the relevant environment variable. For example,

```Shell
PRE_RELEASE=rc.1 misc/release.sh
```



## Schema Migrations
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

Under `migrations` folder add a `rambler.json` config file. An example:

```JSON
{
  "driver": "postgresql",
  "protocol": "tcp",
  "host": "localhost",
  "port": 5432,
  "user": "",
  "password": "",
  "database": "mdb",
  "directory": ".",
  "table": "migrations"
}
```

**Important** make sure never to commit such files to SCM.

On the command line:

```Shell
rambler apply -a
```


## Logging
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


## Installation details

### Postgresql installation

https://wiki.postgresql.org/wiki/Apt

### Go related installations

```Shell
sudo apt-get update
sudo curl -O https://storage.googleapis.com/golang/go1.7.3linux-amd64.tar.gz
```

Detailes can be found here: https://www.digitalocean.com/community/tutorials/how-to-install-go-1-6-on-ubuntu-14-04)

```Shell
sudo tar -xvf go1.7.3.linux-amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/user/local/go
```

### While at /home/kolmanv/go

```Shell
export GOPATH=/home/kolmanv/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
git clone https://github.com/Bnei-Baruch/mdb.git
```

### Install Packages - Using godep
```Shell
go get gopkg.in/gin-gonic/gin.v1
go get github.com/lib/pq
go get github.com/tools/godep
# https://github.com/tools/godep
godep save
```
