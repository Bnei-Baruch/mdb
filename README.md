# MDB - Metadata DB

## Overview

BB archive Metadata Database.

This system aims to be a single source of truth for all content produced by Bnei Baruch. 

## Commands
The mdb is meant to be executed as command line. 
Type `mdb <command> -h` to see how to use each command.
 
`mdb server` 

Run the server

`mdb config <path>`
 
Generate default configuration in the given path. If path is omitted STDOUT is used instead.
  *Note* that default value to config file is `config.toml` in project root directory.

## Implementation Notes

### Dates and Times
All timestamps fields are expecting values in UTC only.


### Languages
Languages are represented in the system as a two letters code adhering to the [ISO_639-1](https://en.wikipedia.org/wiki/ISO_639-1) standard.

Special values:

* Unknown - `xx` 
* Multiple languages - `zz` 

## Installation details

### Postgresql installation

https://wiki.postgresql.org/wiki/Apt

### Go related installations

```shell
sudo apt-get update
sudo curl -O https://storage.googleapis.com/golang/go1.7.3linux-amd64.tar.gz
```

Detailes can be found here: https://www.digitalocean.com/community/tutorials/how-to-install-go-1-6-on-ubuntu-14-04)

```shell
sudo tar -xvf go1.7.3.linux-amd64.tar.gz
sudo mv go /usr/local
export GOROOT=/user/local/go
```

### While at /home/kolmanv/go

```shell
export GOPATH=/home/kolmanv/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
git clone https://github.com/Bnei-Baruch/MDB.git
```

### Install Packages - for now not using any package manager.
```shell
go get gopkg.in/gin-gonic/gin.v1
```
