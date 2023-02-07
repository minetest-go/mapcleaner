mapcleaner
-----------------

![](https://github.com/BuckarooBanzay/mapcleaner/workflows/build/badge.svg)
![](https://github.com/BuckarooBanzay/mapcleaner/workflows/go-test/badge.svg)
![](https://github.com/BuckarooBanzay/mapcleaner/workflows/test/badge.svg)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/minetest-go/mapcleaner)
[![Go Report Card](https://goreportcard.com/badge/github.com/minetest-go/mapcleaner)](https://goreportcard.com/report/github.com/minetest-go/mapcleaner)
[![Coverage Status](https://coveralls.io/repos/github/minetest-go/v/badge.svg)](https://coveralls.io/github/minetest-go/mapcleaner)

# Overview

Removes unused and unprotected map chunks on an online server

# Features

* Removes map chunks based on a nodename whitelist
* Nodename whitelist defines which nodes should be kept
* One layer safety-range around protected chunks
* `areas` mod inetgration
* Can be paused (state and config is maintained in `mapcleaner.json`)

# Installing / Running

* [Download](https://github.com/minetest-go/mapcleaner/releases) the latest binary
* Prepare a `mapcleaner_protect.txt` file with protected nodenames and place it in the world folder (see section below)
* Optionally: configure the `mapcleaner.json` file with start/end positions
* Start the `mapcleaner` binary in your world directory

# Configuration files

## `mapcleaner_protect.txt`

The mapcleaner needs a `mapcleaner_protect.txt` file with nodenames that need to be kept.

An example:
```
bones:bones
travelnet:travelnet
protector:protect
```

With the above file all the chunks without those nodes (or not in the neighborhood of those) will be removed

**NOTE**: the mapcleaner refuses to work without this file (for good reasons)

## `mapcleaner.json`

This file is used for state tracking and cleaning range config.

An example:
```json
{
  "chunk_x": -40,
  "chunk_y": -10,
  "chunk_z": -40,
  "removed_chunks": 0,
  "retained_chunks": 0,
  "processed_chunks": 0,
  "from_x": -40,
  "from_z": -40,
  "to_x": 40,
  "to_y": 10,
  "to_z": 40
}
```

The above config will process all chunks between `(-40,-10,-40)` and `(40,10,40)` (chunk-positions).
This translates roughly to `(-3200,-800,-3200)` and `(3200,800,3200)` in node-positions.

**NOTE**: a chunk has 5x5x5 mapblocks and a mapblock 16x16x16 nodes, chunks have a mapblock-offset of `(-2,-2,-2)`

# Warnings

This program removes data from your world, try it out on a backup/throw-away world first!

# License

MIT
