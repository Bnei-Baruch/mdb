# MDB - Metadata DB

## Overview

Metadata database for BB content.

The main purpose of this repository is to hold the sql migrations and a basic, slim, go ORM model layer.


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

### Install Packages - Using godep
```shell
go get github.com/tools/godep
godep save
# go get gopkg.in/gin-gonic/gin.v1
# go get github.com/lib/pq
```
