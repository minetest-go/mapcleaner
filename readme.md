mapcleaner
-----------------

A mod for [minetest](http://www.minetest.net)

![](https://github.com/thomasrudin-mt/mapcleaner/workflows/luacheck/badge.svg)

# Overview

Removes unprotected map chunks on an online server

# Rationale

TODO

# Process

TODO

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
```

# Safety range

* Player distance: 5*80 blocks
* For a chunk to be removed all surrounding chunks have to be unprotected

# Settings

TODO

# License

MIT
