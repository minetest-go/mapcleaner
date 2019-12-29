-- returns true if the chunk is generated
-- TODO: returns false on edge of the world blocks
function mapcleaner.is_generated(chunk_pos)

	local min_mapblock_pos = mapcleaner.get_mapblocks_from_chunk(chunk_pos)
	local min_pos = mapcleaner.get_blocks_from_mapblock(min_mapblock_pos)

	-- load area
	minetest.get_voxel_manip(min_pos, min_pos)

	local node = minetest.get_node(min_pos)

	return node.name ~= "ignore"
end
