#!/bin/sh

docker run --rm -it \
	-u root:root \
  -v $(pwd)/minetest.conf:/data/minetest.conf \
  -v $(pwd)/world.mt:/data/world/world.mt \
  -v $(pwd)/worldmods:/data/world/worldmods \
	--network host \
	buckaroobanzay/minetest:5.2.0-r1
