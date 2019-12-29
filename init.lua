

local function get_mapblock_from_pos(pos)
	return {
		x = math.floor(pos.x / 16),
		y = math.floor(pos.y / 16),
		z = math.floor(pos.z / 16)
	}
end


local function get_chunkpos_from_pos(pos)
	local mapblock = get_mapblock_from_pos(pos)
	local aligned_mapblock = vector.add(mapblock, {x=2,y=2,z=2})
	return {
		x = math.floor(aligned_mapblock.x / 5),
		y = math.floor(aligned_mapblock.y / 5),
		z = math.floor(aligned_mapblock.z / 5)
	}
end

local function get_mapblocks_from_chunk(chunkpos)
	local minmapblock = {
		x = (chunkpos.x * 5) - 2,
		y = (chunkpos.y * 5) - 2,
		z = (chunkpos.z * 5) - 2,
	}
	local maxmapblock = vector.add(minmapblock, 4)

	return minmapblock, maxmapblock
end

local function get_blocks_from_mapblock(mapblock)
	local min = {
		x = (mapblock.x) * 16,
		y = (mapblock.y) * 16,
		z = (mapblock.z) * 16
	}
	local max = vector.add(min, 15)

	return min, max
end

local function is_mapblock_protected(mapblock_pos)
	local min, max = get_blocks_from_mapblock(mapblock_pos)

	local nodes = minetest.find_nodes_in_area(min, max, {
		"protector:protect",
		"protector:protect2"
	})

	return nodes and #nodes > 0
end

minetest.register_chatcommand("mapcleaner_check", {
	description = "checks the current chunk",
	privs = { server = true },
	func = function(name)
		local player = minetest.get_player_by_name(name)
		if not player then
			return true, "no such player: " .. name
		end

		local pos = player:get_pos()
		local mapblock = get_mapblock_from_pos(pos)
		local chunk = get_chunkpos_from_pos(pos)
		local protected = is_mapblock_protected(mapblock)
		local protected_str = "false"
		if protected then
			protected_str = "true"
		end

		return true, "Mapblock: " .. minetest.pos_to_string(mapblock) ..
			" Chunk: " .. minetest.pos_to_string(chunk) ..
			" Protected: " .. protected_str

	end
})

--[[
minetest.register_on_generated(function(minp, maxp)
	minetest.chat_send_all("minp: " .. minetest.pos_to_string(minp) ..
		" Chunk: " .. minetest.pos_to_string(get_chunkpos_from_pos(minp)))
	minetest.chat_send_all("maxp: " .. minetest.pos_to_string(maxp) ..
		" Chunk: " .. minetest.pos_to_string(get_chunkpos_from_pos(maxp)))
end)
--]]
