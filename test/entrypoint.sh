#!/bin/sh

minetestserver --config /minetest.conf

/mapcleaner

sqlite3 map.sqlite "select count(*) from blocks"

cat mapcleaner.json | jq