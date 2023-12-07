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

# export mode including user specified protected nodes
rm -rf area-export
/mapcleaner -mode export_protected -export-all
test -d area-export
test -f area-export/map.sqlite
# Two chunks are preserved (one with bones and one with areas protection)
# 3x3x3 chunks around each are exported as well, but since they are next
# to each other, some mapblocks are shared and should not be counted twice:
# 3x3x4 chunks should be exported. Each has 5x5x5 mapblocks = 4500 total exported mapblocks.
exported_blocks="$(sqlite3 area-export/map.sqlite "select count(*) from blocks")"
test "${exported_blocks}" == "4500" || exit 1
