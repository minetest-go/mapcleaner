mapcleaner
-----------------

![](https://github.com/BuckarooBanzay/mapcleaner/workflows/docker/badge.svg)
![](https://github.com/BuckarooBanzay/mapcleaner/workflows/test/badge.svg)
![](https://github.com/BuckarooBanzay/mapcleaner/workflows/jshint/badge.svg)

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
