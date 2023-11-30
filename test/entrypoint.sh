#!/bin/sh
set -e
set -x

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

# export mode including protected nodes
rm -rf area-export

/mapcleaner -mode export_protected -export-all

test -d area-export
test -f area-export/map.sqlite

# 2 chunks should be exported: one for areas and one with bones
# Each chunk has 5x5x5 blocks so there should be 250 blocks
exported_blocks="$(sqlite3 area-export/map.sqlite "select count(*) from blocks")"
test "${exported_blocks}" == "250" || exit 1