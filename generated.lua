-- returns true if the chunk is generated
function mapcleaner.is_generated(chunk_pos)

	local min_mapblock_pos, max_mapblock_pos = mapcleaner.get_mapblocks_from_chunk(chunk_pos)
	local min_pos = mapcleaner.get_blocks_from_mapblock(min_mapblock_pos)
	local _, max_pos = mapcleaner.get_blocks_from_mapblock(max_mapblock_pos)

	local check_pos = { x=min_pos.x, y=min_pos.y, z=min_pos.z }

	if chunk_pos.x < 0 then
		check_pos.x = max_pos.x
	end

	if chunk_pos.y < 0 then
		check_pos.y = max_pos.y
	end

	if chunk_pos.z < 0 then
		check_pos.z = max_pos.z
	end

	-- load area
	minetest.get_voxel_manip(check_pos, check_pos)

	local node = minetest.get_node(check_pos)

	local is_generated = node.name ~= "ignore"

	if not is_generated then
		-- clean up afterwards
		-- looks like the above calls create "ignore" only mapblocks on the database
		minetest.delete_area(check_pos, check_pos)
	end

	return is_generated
end
