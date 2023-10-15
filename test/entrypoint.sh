#!/bin/sh
set -e

minetestserver --config /minetest.conf

sqlite3 map.sqlite "select count(*) from blocks"

# prune mode
/mapcleaner

sqlite3 map.sqlite "select count(*) from blocks"

retained=$(cat mapcleaner.json | jq -r ".retained_chunks")

# 4x3x3 chunks remain
test "${retained}" == "36" || exit 1

# export mode
/mapcleaner -mode export_protected

test -d area-export
test -f area-export/map.sqlite

sqlite3 area-export/map.sqlite "select count(*) from blocks"
