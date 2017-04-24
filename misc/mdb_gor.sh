#!/usr/bin/env bash
# HTTP traffic log using gor.
# Checkout https://goreplay.org/

/sites/mdb/gor --input-raw :8080 \
--input-raw-track-response \
--http-allow-method POST \
--http-allow-method DELETE \
--http-allow-method PUT \
--http-allow-method PATCH \
--output-stdout
#--debug --stats --verbose
