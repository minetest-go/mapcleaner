-- returns true if the chunk is generated
function mapcleaner.is_generated(chunk_pos)

	local min_mapblock_pos, max_mapblock_pos = mapcleaner.get_mapblocks_from_chunk(chunk_pos)
	local min_pos = mapcleaner.get_blocks_from_mapblock(min_mapblock_pos)
	local _, max_pos = mapcleaner.get_blocks_from_mapblock(max_mapblock_pos)

	-- load area
	minetest.get_voxel_manip(min_pos, min_pos)
	minetest.get_voxel_manip(max_pos, max_pos)

	local min_node = minetest.get_node(min_pos)
	local max_node = minetest.get_node(max_pos)

	local is_generated = min_node.name ~= "ignore" or max_node.name ~= "ignore"

	if not is_generated then
		-- clean up afterwards
		-- looks like the above calls create "ignore" only mapblocks on the database
		minetest.delete_area(min_pos, min_pos)
		minetest.delete_area(max_pos, max_pos)
	end

	return is_generated
end
