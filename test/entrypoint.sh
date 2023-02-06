#!/bin/sh
set -e

minetestserver --config /minetest.conf

sqlite3 map.sqlite "select count(*) from blocks"

/mapcleaner

sqlite3 map.sqlite "select count(*) from blocks"

retained=$(cat mapcleaner.json | jq -r ".retained_chunks")

# 4x3x3 chunks remain
test "${retained}" == "36" || exit 1