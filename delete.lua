function mapcleaner.delete_chunk(chunk_pos)
	minetest.log("warning", "[mapcleaner] removing chunk: " .. minetest.pos_to_string(chunk_pos))
	local min_mapblock_pos, max_mapblock_pos = mapcleaner.get_mapblocks_from_chunk(chunk_pos)

	for x=min_mapblock_pos.x,max_mapblock_pos.x do
		for y=min_mapblock_pos.y,max_mapblock_pos.y do
			for z=min_mapblock_pos.z,max_mapblock_pos.z do
				local mapblock_pos = {x=x, y=y, z=z}
				local pos = mapcleaner.get_blocks_from_mapblock(mapblock_pos)
				minetest.log("warning", "[mapcleaner] removing mapblock at pos: " .. minetest.pos_to_string(mapblock_pos))
				minetest.delete_area(pos, pos)
			end
		end
	end
end
