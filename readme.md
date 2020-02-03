mapcleaner
-----------------

A mod for [minetest](http://www.minetest.net)

![](https://github.com/thomasrudin-mt/mapcleaner/workflows/luacheck/badge.svg)

# Overview

Removes unprotected map chunks on an online server

# Rationale

Cleaning up unused an stale mapgen only blocks

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

# Chatcommands

* **/mapcleaner_status** shows the current status (current position, statistics)
* **/mapcleaner_max_time [microseconds]** shows or sets the max cpu-time per step
* **/mapcleaner_step_interval [seconds]** shows or sets the seconds between intervals
* **/mapcleaner_run [true|false]** show, start or stop the process
* **mapcleaner_max_lag [seconds]** sets the max lag value where the mapcleaner stops running

All settings are persisted across restarts

**NOTE**: by default the process is stopped and has to be started initially with **/mapcleaner_run true**

# License

MIT
