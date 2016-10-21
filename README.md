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

