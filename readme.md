mapcleaner
-----------------

![](https://github.com/BuckarooBanzay/mapcleaner/workflows/build/badge.svg)
![](https://github.com/BuckarooBanzay/mapcleaner/workflows/go-test/badge.svg)
![](https://github.com/BuckarooBanzay/mapcleaner/workflows/test/badge.svg)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/minetest-go/mapcleaner)
[![Go Report Card](https://goreportcard.com/badge/github.com/minetest-go/mapcleaner)](https://goreportcard.com/report/github.com/minetest-go/mapcleaner)
[![Coverage Status](https://coveralls.io/repos/github/minetest-go/v/badge.svg)](https://coveralls.io/github/minetest-go/mapcleaner)

# Overview

Removes unprotected map chunks on an online server

# Rationale

Cleaning up unused and stale mapgen only blocks

# Process

Runs in the background, checks and removes chunks from bottom to top.

# Protected nodes

TODO

# Safety range

* For a chunk to be removed all surrounding chunks have to be unprotected

# License

MIT
