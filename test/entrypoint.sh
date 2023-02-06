#!/bin/sh

minetestserver --config /minetest.conf

/mapcleaner

sqlite3 map.sqlite "count(*) from blocks"