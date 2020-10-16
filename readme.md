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

List of protected nodes:

```
protector:protect
protector:protect2
travelnet:travelnet
xp_redo:protector
priv_protector:protector
default:chest_protected
bones:bones
advtrains:dtrack_st
advtrains:dtrack_st_45
advtrains:dtrack_cr_60
group:save_in_at_nodedb
```

# Safety range

* Player distance: 5*80 blocks
* For a chunk to be removed all surrounding chunks have to be unprotected

# License

MIT
