
function mapcleaner.get_mapblock_from_pos(pos)
	return {
		x = math.floor(pos.x / 16),
		y = math.floor(pos.y / 16),
		z = math.floor(pos.z / 16)
	}
end


function mapcleaner.get_chunkpos_from_pos(pos)
	local mapblock = mapcleaner.get_mapblock_from_pos(pos)
	local aligned_mapblock = vector.add(mapblock, {x=2,y=2,z=2})
	return {
		x = math.floor(aligned_mapblock.x / 5),
		y = math.floor(aligned_mapblock.y / 5),
		z = math.floor(aligned_mapblock.z / 5)
	}
end

function mapcleaner.get_mapblocks_from_chunk(chunkpos)
	local minmapblock = {
		x = (chunkpos.x * 5) - 2,
		y = (chunkpos.y * 5) - 2,
		z = (chunkpos.z * 5) - 2,
	}
	local maxmapblock = vector.add(minmapblock, 4)

	return minmapblock, maxmapblock
end

function mapcleaner.get_blocks_from_mapblock(mapblock)
	local min = {
		x = (mapblock.x) * 16,
		y = (mapblock.y) * 16,
		z = (mapblock.z) * 16
	}
	local max = vector.add(min, 15)

	return min, max
end
